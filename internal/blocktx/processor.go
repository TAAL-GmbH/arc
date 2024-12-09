package blocktx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/libsv/go-bc"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/libsv/go-p2p/wire"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	blockchain "github.com/bitcoin-sv/arc/internal/blocktx/blockchain_communication"
	blocktx_p2p "github.com/bitcoin-sv/arc/internal/blocktx/blockchain_communication/p2p"
	"github.com/bitcoin-sv/arc/internal/blocktx/blocktx_api"
	"github.com/bitcoin-sv/arc/internal/blocktx/store"
	"github.com/bitcoin-sv/arc/internal/tracing"
)

var (
	ErrFailedToSubscribeToTopic        = errors.New("failed to subscribe to register topic")
	ErrFailedToCreateBUMP              = errors.New("failed to create new bump for tx hash from merkle tree and index")
	ErrFailedToGetStringFromBUMPHex    = errors.New("failed to get string from bump for tx hash")
	ErrFailedToParseBlockHash          = errors.New("failed to parse block hash")
	ErrFailedToInsertBlockTransactions = errors.New("failed to insert block transactions")
	ErrBlockAlreadyExists              = errors.New("block already exists in the database")
	ErrUnexpectedBlockStatus           = errors.New("unexpected block status")
	ErrFailedToProcessBlock            = errors.New("failed to process block")
	ErrFailedToStartCollectingStats    = errors.New("failed to start collecting stats")
)

const (
	transactionStoringBatchsizeDefault = 8192 // power of 2 for easier memory allocation
	maxRequestBlocks                   = 10
	maxBlocksInProgress                = 1
	registerTxsIntervalDefault         = time.Second * 10
	registerRequestTxsIntervalDefault  = time.Second * 5
	registerTxsBatchSizeDefault        = 100
	registerRequestTxBatchSizeDefault  = 100
	waitForBlockProcessing             = 5 * time.Minute
)

type Processor struct {
	hostname                    string
	blockRequestCh              chan blocktx_p2p.BlockRequest
	blockProcessCh              chan *blockchain.BlockMessage
	store                       store.BlocktxStore
	logger                      *slog.Logger
	transactionStorageBatchSize int
	dataRetentionDays           int
	mqClient                    MessageQueueClient
	registerTxsChan             chan []byte
	requestTxChannel            chan []byte
	registerTxsInterval         time.Duration
	registerRequestTxsInterval  time.Duration
	registerTxsBatchSize        int
	registerRequestTxsBatchSize int
	tracingEnabled              bool
	tracingAttributes           []attribute.KeyValue
	processGuardsMap            sync.Map
	stats                       *processorStats
	statCollectionInterval      time.Duration

	now                        func() time.Time
	maxBlockProcessingDuration time.Duration

	waitGroup *sync.WaitGroup
	cancelAll context.CancelFunc
	ctx       context.Context
}

func NewProcessor(
	logger *slog.Logger,
	storeI store.BlocktxStore,
	blockRequestCh chan blocktx_p2p.BlockRequest,
	blockProcessCh chan *blockchain.BlockMessage,
	opts ...func(*Processor),
) (*Processor, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	p := &Processor{
		store:                       storeI,
		logger:                      logger.With(slog.String("module", "processor")),
		blockRequestCh:              blockRequestCh,
		blockProcessCh:              blockProcessCh,
		transactionStorageBatchSize: transactionStoringBatchsizeDefault,
		registerTxsInterval:         registerTxsIntervalDefault,
		registerRequestTxsInterval:  registerRequestTxsIntervalDefault,
		registerTxsBatchSize:        registerTxsBatchSizeDefault,
		registerRequestTxsBatchSize: registerRequestTxBatchSizeDefault,
		maxBlockProcessingDuration:  waitForBlockProcessing,
		hostname:                    hostname,
		stats:                       newProcessorStats(),
		statCollectionInterval:      statCollectionIntervalDefault,
		now:                         time.Now,
		waitGroup:                   &sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(p)
	}

	ctx, cancelAll := context.WithCancel(context.Background())
	p.cancelAll = cancelAll
	p.ctx = ctx

	return p, nil
}

func (p *Processor) Start(statsEnabled bool) error {
	err := p.mqClient.Subscribe(RegisterTxTopic, func(msg []byte) error {
		p.registerTxsChan <- msg
		return nil
	})
	if err != nil {
		return errors.Join(ErrFailedToSubscribeToTopic, fmt.Errorf("topic: %s", RegisterTxTopic), err)
	}

	err = p.mqClient.Subscribe(RequestTxTopic, func(msg []byte) error {
		p.requestTxChannel <- msg
		return nil
	})
	if err != nil {
		return errors.Join(ErrFailedToSubscribeToTopic, fmt.Errorf("topic: %s", RequestTxTopic), err)
	}

	if statsEnabled {
		err = p.StartCollectStats()
		if err != nil {
			return errors.Join(ErrFailedToStartCollectingStats, err)
		}
	}
	p.StartBlockRequesting()
	p.StartBlockProcessing()
	p.StartProcessRegisterTxs()
	p.StartProcessRequestTxs()

	return nil
}

func (p *Processor) StartBlockRequesting() {
	p.waitGroup.Add(1)

	waitUntilFree := func(ctx context.Context) bool {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		for {
			bhs, err := p.store.GetBlockHashesProcessingInProgress(p.ctx, p.hostname)
			if err != nil {
				p.logger.Error("failed to get block hashes where processing in progress", slog.String("err", err.Error()))
			}

			if len(bhs) < maxBlocksInProgress && err == nil {
				return true
			}

			select {
			case <-ctx.Done():
				return false

			case <-t.C:
			}
		}
	}

	go func() {
		defer p.waitGroup.Done()
		for {
			select {
			case <-p.ctx.Done():
				return
			case req := <-p.blockRequestCh:
				hash := req.Hash
				peer := req.Peer

				if ok := waitUntilFree(p.ctx); !ok {
					continue
				}

				// lock block for the current instance to process
				processedBy, err := p.store.SetBlockProcessing(p.ctx, hash, p.hostname)
				if err != nil {
					// block is already being processed by another blocktx instance
					if errors.Is(err, store.ErrBlockProcessingDuplicateKey) {
						p.logger.Debug("block processing already in progress", slog.String("hash", hash.String()), slog.String("processed_by", processedBy))
						continue
					}

					p.logger.Error("failed to set block processing", slog.String("hash", hash.String()), slog.String("err", err.Error()))
					continue
				}

				p.logger.Info("Sending block request", slog.String("hash", hash.String()))
				msg := wire.NewMsgGetDataSizeHint(1)
				_ = msg.AddInvVect(wire.NewInvVect(wire.InvTypeBlock, hash)) // ignore error at this point
				peer.WriteMsg(msg)

				p.startBlockProcessGuard(p.ctx, hash)
				p.logger.Info("Block request message sent to peer", slog.String("hash", hash.String()), slog.String("peer", peer.String()))
			}
		}
	}()
}

func (p *Processor) StartBlockProcessing() {
	p.waitGroup.Add(1)

	go func() {
		defer p.waitGroup.Done()

		for {
			select {
			case <-p.ctx.Done():
				return
			case blockMsg := <-p.blockProcessCh:
				var err error
				timeStart := time.Now()

				p.logger.Info("received block", slog.String("hash", blockMsg.Hash.String()))

				err = p.processBlock(blockMsg)
				if err != nil {
					p.logger.Error("block processing failed", slog.String("hash", blockMsg.Hash.String()), slog.String("err", err.Error()))
					p.unlockBlock(p.ctx, blockMsg.Hash)
					continue
				}

				storeErr := p.store.MarkBlockAsDone(p.ctx, blockMsg.Hash, blockMsg.Size, uint64(len(blockMsg.TransactionHashes)))
				if storeErr != nil {
					p.logger.Error("unable to mark block as processed", slog.String("hash", blockMsg.Hash.String()), slog.String("err", storeErr.Error()))
					p.unlockBlock(p.ctx, blockMsg.Hash)
					continue
				}

				// add the total block processing time to the stats
				p.logger.Info("Processed block", slog.String("hash", blockMsg.Hash.String()), slog.Int("txs", len(blockMsg.TransactionHashes)), slog.String("duration", time.Since(timeStart).String()))
			}
		}
	}()
}

func (p *Processor) startBlockProcessGuard(ctx context.Context, hash *chainhash.Hash) {
	p.waitGroup.Add(1)

	execCtx, stopFn := context.WithCancel(ctx)
	p.processGuardsMap.Store(*hash, stopFn)

	go func() {
		defer p.waitGroup.Done()
		defer p.processGuardsMap.Delete(*hash)

		select {
		case <-execCtx.Done():
			// we may do nothing here:
			// 1. block processing is completed, or
			// 2. processor is shutting down – all unprocessed blocks are released in the Shutdown func
			return

		case <-time.After(p.maxBlockProcessingDuration):
			// check if block was processed successfully
			block, _ := p.store.GetBlock(execCtx, hash)

			if block != nil {
				return // success
			}

			p.logger.Warn(fmt.Sprintf("block was not processed after %v. Unlock the block to be processed later", waitForBlockProcessing), slog.String("hash", hash.String()))
			p.unlockBlock(execCtx, hash)
		}
	}()
}

func (p *Processor) stopBlockProcessGuard(hash *chainhash.Hash) {
	stopFn, found := p.processGuardsMap.Load(*hash)
	if found {
		stopFn.(context.CancelFunc)()
	}
}

// unlock block for future processing
func (p *Processor) unlockBlock(ctx context.Context, hash *chainhash.Hash) {
	// use closures for retries
	unlockFn := func() error {
		_, err := p.store.DelBlockProcessing(ctx, hash, p.hostname)
		if errors.Is(err, store.ErrBlockNotFound) {
			return nil // block is already unlocked
		}

		return err
	}

	var bo backoff.BackOff
	bo = backoff.NewConstantBackOff(100 * time.Millisecond)
	bo = backoff.WithContext(bo, ctx)
	bo = backoff.WithMaxRetries(bo, 5)

	if unlockErr := backoff.Retry(unlockFn, bo); unlockErr != nil {
		p.logger.ErrorContext(ctx, "failed to delete block processing", slog.String("hash", hash.String()), slog.String("err", unlockErr.Error()))
	}
}

func (p *Processor) StartProcessRegisterTxs() {
	p.waitGroup.Add(1)
	txHashes := make([][]byte, 0, p.registerTxsBatchSize)

	ticker := time.NewTicker(p.registerTxsInterval)
	go func() {
		defer p.waitGroup.Done()
		for {
			select {
			case <-p.ctx.Done():
				return
			case txHash := <-p.registerTxsChan:
				txHashes = append(txHashes, txHash)

				if len(txHashes) < p.registerTxsBatchSize {
					continue
				}

				p.registerTransactions(txHashes[:])
				txHashes = txHashes[:0]
				ticker.Reset(p.registerTxsInterval)

			case <-ticker.C:
				if len(txHashes) == 0 {
					continue
				}

				p.registerTransactions(txHashes[:])
				txHashes = txHashes[:0]
				ticker.Reset(p.registerTxsInterval)
			}
		}
	}()
}

func (p *Processor) StartProcessRequestTxs() {
	p.waitGroup.Add(1)

	txHashes := make([]*chainhash.Hash, 0, p.registerRequestTxsBatchSize)

	ticker := time.NewTicker(p.registerRequestTxsInterval)

	go func() {
		defer p.waitGroup.Done()

		for {
			select {
			case <-p.ctx.Done():
				return
			case txHash := <-p.requestTxChannel:
				tx, err := chainhash.NewHash(txHash)
				if err != nil {
					p.logger.Error("Failed to create hash from byte array", slog.String("err", err.Error()))
					continue
				}

				txHashes = append(txHashes, tx)

				if len(txHashes) < p.registerRequestTxsBatchSize || len(txHashes) == 0 {
					continue
				}

				err = p.publishMinedTxs(txHashes)
				if err != nil {
					p.logger.Error("failed to publish mined txs", slog.String("err", err.Error()))
					continue // retry, don't clear the txHashes slice
				}

				txHashes = make([]*chainhash.Hash, 0, p.registerRequestTxsBatchSize)
				ticker.Reset(p.registerRequestTxsInterval)

			case <-ticker.C:
				if len(txHashes) == 0 {
					continue
				}

				err := p.publishMinedTxs(txHashes)
				if err != nil {
					p.logger.Error("failed to publish mined txs", slog.String("err", err.Error()))
					ticker.Reset(p.registerRequestTxsInterval)
					continue // retry, don't clear the txHashes slice
				}

				txHashes = make([]*chainhash.Hash, 0, p.registerRequestTxsBatchSize)
				ticker.Reset(p.registerRequestTxsInterval)
			}
		}
	}()
}

func (p *Processor) publishMinedTxs(txHashes []*chainhash.Hash) error {
	hashesBytes := make([][]byte, len(txHashes))
	for i, h := range txHashes {
		hashesBytes[i] = h[:]
	}

	minedTxs, err := p.store.GetMinedTransactions(p.ctx, hashesBytes, false)
	if err != nil {
		return fmt.Errorf("failed to get mined transactions: %v", err)
	}

	for _, minedTx := range minedTxs {
		txBlock := &blocktx_api.TransactionBlock{
			TransactionHash: minedTx.TxHash,
			BlockHash:       minedTx.BlockHash,
			BlockHeight:     minedTx.BlockHeight,
			MerklePath:      minedTx.MerklePath,
			BlockStatus:     minedTx.BlockStatus,
		}
		err = p.mqClient.PublishMarshal(p.ctx, MinedTxsTopic, txBlock)
	}

	if err != nil {
		return fmt.Errorf("failed to publish mined transactions: %v", err)
	}

	return nil
}

func (p *Processor) registerTransactions(txHashes [][]byte) {
	updatedTxs, err := p.store.RegisterTransactions(p.ctx, txHashes)
	if err != nil {
		p.logger.Error("failed to register transactions", slog.String("err", err.Error()))
	}

	if len(updatedTxs) > 0 {
		err = p.publishMinedTxs(updatedTxs)
		if err != nil {
			p.logger.Error("failed to publish mined txs", slog.String("err", err.Error()))
		}
	}
}

func (p *Processor) buildMerkleTreeStoreChainHash(ctx context.Context, txids []*chainhash.Hash) []*chainhash.Hash {
	_, span := tracing.StartTracing(ctx, "buildMerkleTreeStoreChainHash", p.tracingEnabled, p.tracingAttributes...)
	defer tracing.EndTracing(span, nil)

	return bc.BuildMerkleTreeStoreChainHash(txids)
}

func (p *Processor) processBlock(blockMsg *blockchain.BlockMessage) (err error) {
	ctx := p.ctx

	var block *blocktx_api.Block

	// release guardian
	defer p.stopBlockProcessGuard(blockMsg.Hash)

	ctx, span := tracing.StartTracing(ctx, "processBlock", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.String("hash", blockMsg.Hash.String()))
			if block != nil {
				span.SetAttributes(attribute.String("status", block.Status.String()))
			}
		}

		tracing.EndTracing(span, err)
	}()

	p.logger.Info("processing incoming block", slog.String("hash", blockMsg.Hash.String()), slog.Uint64("height", blockMsg.Height))

	// check if we've already processed that block
	existingBlock, _ := p.store.GetBlock(ctx, blockMsg.Hash)

	if existingBlock != nil {
		p.logger.Warn("ignoring already existing block", slog.String("hash", blockMsg.Hash.String()), slog.Uint64("height", blockMsg.Height))
		return nil
	}

	block, err = p.verifyAndInsertBlock(ctx, blockMsg)
	if err != nil {
		return err
	}

	var longestTxs, staleTxs []store.TransactionBlock
	var ok bool

	switch block.Status {
	case blocktx_api.Status_LONGEST:
		longestTxs, ok = p.getRegisteredTransactions(ctx, []*blocktx_api.Block{block})
	case blocktx_api.Status_STALE:
		longestTxs, staleTxs, ok = p.handleStaleBlock(ctx, block)
	case blocktx_api.Status_ORPHANED:
		longestTxs, staleTxs, ok = p.handleOrphans(ctx, block)
	default:
		return ErrUnexpectedBlockStatus
	}

	if !ok {
		// error is already logged in each method above
		return ErrFailedToProcessBlock
	}

	p.publishTxsToMetamorph(ctx, longestTxs)
	p.publishTxsToMetamorph(ctx, staleTxs)

	return nil
}

func (p *Processor) verifyAndInsertBlock(ctx context.Context, blockMsg *blockchain.BlockMessage) (incomingBlock *blocktx_api.Block, err error) {
	ctx, span := tracing.StartTracing(ctx, "verifyAndInsertBlock", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	previousBlockHash := blockMsg.Header.PrevBlock
	merkleRoot := blockMsg.Header.MerkleRoot

	incomingBlock = &blocktx_api.Block{
		Hash:         blockMsg.Hash[:],
		PreviousHash: previousBlockHash[:],
		MerkleRoot:   merkleRoot[:],
		Height:       blockMsg.Height,
		Chainwork:    calculateChainwork(blockMsg.Header.Bits).String(),
	}

	err = p.assignBlockStatus(ctx, incomingBlock, previousBlockHash)
	if err != nil {
		p.logger.Error("unable to assign block status", slog.String("hash", blockMsg.Hash.String()), slog.Uint64("height", incomingBlock.Height), slog.String("err", err.Error()))
		return nil, err
	}

	p.logger.Info("Inserting block", slog.String("hash", blockMsg.Hash.String()), slog.Uint64("height", incomingBlock.Height), slog.String("status", incomingBlock.Status.String()))

	err = p.insertBlockAndStoreTransactions(ctx, incomingBlock, blockMsg.TransactionHashes, blockMsg.Header.MerkleRoot)
	if err != nil {
		p.logger.Error("unable to insert block and store its transactions", slog.String("hash", blockMsg.Hash.String()), slog.Uint64("height", incomingBlock.Height), slog.String("err", err.Error()))
		return nil, err
	}

	return incomingBlock, nil
}

func (p *Processor) assignBlockStatus(ctx context.Context, block *blocktx_api.Block, prevBlockHash chainhash.Hash) (err error) {
	ctx, span := tracing.StartTracing(ctx, "assignBlockStatus", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	prevBlock, _ := p.store.GetBlock(ctx, &prevBlockHash)

	if prevBlock == nil {
		// This check is only in case there's a fresh, empty database
		// with no blocks, to mark the first block as the LONGEST chain
		var longestTipExists bool
		longestTipExists, err = p.longestTipExists(ctx)
		if err != nil {
			p.logger.Error("unable to verify the longest tip existence in db", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("err", err.Error()))
			return err
		}

		// if there's no longest block in the
		// database - mark this block as LONGEST
		// otherwise - it's an orphan
		if !longestTipExists {
			block.Status = blocktx_api.Status_LONGEST
		} else {
			block.Status = blocktx_api.Status_ORPHANED
		}
		return nil
	}

	if prevBlock.Status == blocktx_api.Status_LONGEST {
		var competingBlock *blocktx_api.Block
		competingBlock, err = p.store.GetLongestBlockByHeight(ctx, block.Height)
		if err != nil && !errors.Is(err, store.ErrBlockNotFound) {
			p.logger.Error("unable to get the competing block from db", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("err", err.Error()))
			return err
		}

		if competingBlock == nil {
			block.Status = blocktx_api.Status_LONGEST
			return nil
		}

		if bytes.Equal(block.Hash, competingBlock.Hash) {
			// this means that another instance is already processing
			// or have processed this block that we're processing here
			// so we can throw an error and finish processing
			err = ErrBlockAlreadyExists
			return err
		}

		block.Status = blocktx_api.Status_STALE
		return nil
	}

	// ORPHANED or STALE
	block.Status = prevBlock.Status

	return nil
}

func (p *Processor) longestTipExists(ctx context.Context) (bool, error) {
	_, err := p.store.GetChainTip(ctx)
	if err != nil && !errors.Is(err, store.ErrBlockNotFound) {
		return false, err
	}

	if errors.Is(err, store.ErrBlockNotFound) {
		return false, nil
	}

	return true, nil
}

func (p *Processor) getRegisteredTransactions(ctx context.Context, blocks []*blocktx_api.Block) (txsToPublish []store.TransactionBlock, ok bool) {
	var err error
	ctx, span := tracing.StartTracing(ctx, "getRegisteredTransactions", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	blockHashes := make([][]byte, len(blocks))
	for i, b := range blocks {
		blockHashes[i] = b.Hash
	}

	txsToPublish, err = p.store.GetRegisteredTxsByBlockHashes(ctx, blockHashes)
	if err != nil {
		block := blocks[len(blocks)-1]
		p.logger.Error("unable to get registered transactions", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("err", err.Error()))
		return nil, false
	}

	return txsToPublish, true
}

func (p *Processor) insertBlockAndStoreTransactions(ctx context.Context, incomingBlock *blocktx_api.Block, txHashes []*chainhash.Hash, merkleRoot chainhash.Hash) (err error) {
	ctx, span := tracing.StartTracing(ctx, "insertBlockAndStoreTransactions", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	calculatedMerkleTree := p.buildMerkleTreeStoreChainHash(ctx, txHashes)
	if !merkleRoot.IsEqual(calculatedMerkleTree[len(calculatedMerkleTree)-1]) {
		p.logger.Error("merkle root mismatch", slog.String("hash", getHashStringNoErr(incomingBlock.Hash)))
		return err
	}

	blockID, err := p.store.UpsertBlock(ctx, incomingBlock)
	if err != nil {
		p.logger.Error("unable to insert block at given height", slog.String("hash", getHashStringNoErr(incomingBlock.Hash)), slog.Uint64("height", incomingBlock.Height), slog.String("err", err.Error()))
		return err
	}

	if err = p.storeTransactions(ctx, blockID, incomingBlock, calculatedMerkleTree); err != nil {
		p.logger.Error("unable to store transactions from block", slog.String("hash", getHashStringNoErr(incomingBlock.Hash)), slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (p *Processor) storeTransactions(ctx context.Context, blockID uint64, block *blocktx_api.Block, merkleTree []*chainhash.Hash) (err error) {
	ctx, span := tracing.StartTracing(ctx, "storeTransactions", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	txs := make([]store.TxWithMerklePath, 0, p.transactionStorageBatchSize)
	leaves := merkleTree[:(len(merkleTree)+1)/2]

	blockhash, err := chainhash.NewHash(block.Hash)
	if err != nil {
		return errors.Join(ErrFailedToParseBlockHash, fmt.Errorf("block height: %d", block.Height), err)
	}

	var totalSize int
	for totalSize = 1; totalSize < len(leaves); totalSize++ {
		if leaves[totalSize] == nil {
			// Everything to the right of the first nil will also be nil, as this is just padding upto the next PoT.
			break
		}
	}

	progress := progressIndices(totalSize, 5)
	now := time.Now()

	var iterateMerkleTree trace.Span
	ctx, iterateMerkleTree = tracing.StartTracing(ctx, "iterateMerkleTree", p.tracingEnabled, p.tracingAttributes...)

	for txIndex, hash := range leaves {
		// Everything to the right of the first nil will also be nil, as this is just padding upto the next PoT.
		if hash == nil {
			break
		}

		bump, err := bc.NewBUMPFromMerkleTreeAndIndex(block.Height, merkleTree, uint64(txIndex)) // #nosec G115
		if err != nil {
			return errors.Join(ErrFailedToCreateBUMP, fmt.Errorf("tx hash %s, block height: %d", hash.String(), block.Height), err)
		}

		bumpHex, err := bump.String()
		if err != nil {
			return errors.Join(ErrFailedToGetStringFromBUMPHex, err)
		}

		txs = append(txs, store.TxWithMerklePath{
			Hash:       hash[:],
			MerklePath: bumpHex,
		})

		if (txIndex+1)%p.transactionStorageBatchSize == 0 {
			err := p.store.UpsertBlockTransactions(ctx, blockID, txs)
			if err != nil {
				return errors.Join(ErrFailedToInsertBlockTransactions, err)
			}
			// free up memory
			txs = txs[:0]
		}

		if percentage, found := progress[txIndex+1]; found {
			if totalSize > 0 {
				p.logger.Info(fmt.Sprintf("%d txs out of %d stored", txIndex+1, totalSize), slog.Int("percentage", percentage), slog.String("hash", blockhash.String()), slog.Uint64("height", block.Height), slog.String("duration", time.Since(now).String()))
			}
		}
	}

	tracing.EndTracing(iterateMerkleTree, nil)

	// update all remaining transactions
	err = p.store.UpsertBlockTransactions(ctx, blockID, txs)
	if err != nil {
		return errors.Join(ErrFailedToInsertBlockTransactions, fmt.Errorf("block height: %d", block.Height), err)
	}

	return nil
}

func (p *Processor) handleStaleBlock(ctx context.Context, block *blocktx_api.Block) (longestTxs, staleTxs []store.TransactionBlock, ok bool) {
	var err error
	ctx, span := tracing.StartTracing(ctx, "handleStaleBlock", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	staleBlocks, err := p.store.GetStaleChainBackFromHash(ctx, block.Hash)
	if err != nil {
		p.logger.Error("unable to get STALE blocks to verify chainwork", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("err", err.Error()))
		return nil, nil, false
	}

	lowestHeight := block.Height
	if len(staleBlocks) > 0 {
		lowestHeight = staleBlocks[0].Height
	}

	longestBlocks, err := p.store.GetLongestChainFromHeight(ctx, lowestHeight)
	if err != nil {
		p.logger.Error("unable to get LONGEST blocks to verify chainwork", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("err", err.Error()))
		return nil, nil, false
	}

	staleChainwork := sumChainwork(staleBlocks)
	longestChainwork := sumChainwork(longestBlocks)

	if longestChainwork.Cmp(staleChainwork) < 0 {
		p.logger.Info("chain reorg detected", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height))

		longestTxs, staleTxs, err = p.performReorg(ctx, staleBlocks, longestBlocks)
		if err != nil {
			p.logger.Error("unable to perform reorg", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("err", err.Error()))
			return nil, nil, false
		}
		return longestTxs, staleTxs, true
	}

	return nil, nil, true
}

func (p *Processor) performReorg(ctx context.Context, staleBlocks []*blocktx_api.Block, longestBlocks []*blocktx_api.Block) (longestTxs, staleTxs []store.TransactionBlock, err error) {
	ctx, span := tracing.StartTracing(ctx, "performReorg", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	staleHashes := make([][]byte, len(staleBlocks))
	longestHashes := make([][]byte, len(longestBlocks))

	blockStatusUpdates := make([]store.BlockStatusUpdate, len(longestBlocks)+len(staleBlocks))

	for i, b := range longestBlocks {
		longestHashes[i] = b.Hash

		b.Status = blocktx_api.Status_STALE
		update := store.BlockStatusUpdate{Hash: b.Hash, Status: b.Status}
		blockStatusUpdates[i] = update
	}

	for i, b := range staleBlocks {
		staleHashes[i] = b.Hash

		b.Status = blocktx_api.Status_LONGEST
		update := store.BlockStatusUpdate{Hash: b.Hash, Status: b.Status}
		blockStatusUpdates[i+len(longestBlocks)] = update
	}

	err = p.store.UpdateBlocksStatuses(ctx, blockStatusUpdates)
	if err != nil {
		return nil, nil, err
	}

	p.logger.Info("reorg performed successfully")

	// now the previously stale chain is the longest,
	// so longestTxs are from previously stale block hashes
	longestTxs, err = p.store.GetRegisteredTxsByBlockHashes(ctx, staleHashes)
	if err != nil {
		return nil, nil, err
	}

	// now the previously longest chain is stale,
	// so staleTxs are from previously longest block hashes
	staleTxs, err = p.store.GetRegisteredTxsByBlockHashes(ctx, longestHashes)
	if err != nil {
		return nil, nil, err
	}

	staleTxs = exclusiveRightTxs(longestTxs, staleTxs)

	return longestTxs, staleTxs, nil
}

func (p *Processor) handleOrphans(ctx context.Context, block *blocktx_api.Block) (longestTxs, staleTxs []store.TransactionBlock, ok bool) {
	var err error
	ctx, span := tracing.StartTracing(ctx, "handleOrphans", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	orphans, ancestor, err := p.store.GetOrphansBackToNonOrphanAncestor(ctx, block.Hash)
	if err != nil {
		p.logger.Error("unable to get ORPHANED blocks", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("err", err.Error()))
		return nil, nil, false
	}

	if ancestor == nil || len(orphans) == 0 {
		return nil, nil, true
	}

	p.logger.Info("orphaned chain found", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("status", block.Status.String()))

	if ancestor.Status == blocktx_api.Status_STALE {
		ok = p.acceptIntoChain(ctx, orphans, ancestor.Status)
		if !ok {
			return nil, nil, false
		}

		block.Status = blocktx_api.Status_STALE
		return p.handleStaleBlock(ctx, block)
	}

	if ancestor.Status == blocktx_api.Status_LONGEST {
		// If there is competing block at the height of
		// the first orphan, then we need to mark them
		// all as stale and recheck for reorg.
		//
		// If there's no competing block at the height
		// of the first orphan, then we can assume that
		// there's no competing chain at all.

		var competingBlock *blocktx_api.Block
		competingBlock, err = p.store.GetLongestBlockByHeight(ctx, orphans[0].Height)
		if err != nil && !errors.Is(err, store.ErrBlockNotFound) {
			p.logger.Error("unable to get competing block when handling orphans", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height), slog.String("err", err.Error()))
			return nil, nil, false
		}

		if competingBlock != nil && !bytes.Equal(competingBlock.Hash, orphans[0].Hash) {
			ok = p.acceptIntoChain(ctx, orphans, blocktx_api.Status_STALE)
			if !ok {
				return nil, nil, false
			}

			block.Status = blocktx_api.Status_STALE
			return p.handleStaleBlock(ctx, block)
		}

		ok = p.acceptIntoChain(ctx, orphans, ancestor.Status) // LONGEST
		if !ok {
			return nil, nil, false
		}

		p.logger.Info("orphaned chain accepted into LONGEST chain", slog.String("hash", getHashStringNoErr(block.Hash)), slog.Uint64("height", block.Height))
		longestTxs, ok = p.getRegisteredTransactions(ctx, orphans)
		return longestTxs, nil, ok
	}

	return nil, nil, true
}

func (p *Processor) acceptIntoChain(ctx context.Context, blocks []*blocktx_api.Block, chain blocktx_api.Status) (ok bool) {
	var err error
	ctx, span := tracing.StartTracing(ctx, "acceptIntoChain", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	blockStatusUpdates := make([]store.BlockStatusUpdate, len(blocks))

	for i, b := range blocks {
		b.Status = chain
		blockStatusUpdates[i] = store.BlockStatusUpdate{
			Hash:   b.Hash,
			Status: b.Status,
		}
	}

	tip := blocks[len(blocks)-1]

	err = p.store.UpdateBlocksStatuses(ctx, blockStatusUpdates)
	if err != nil {
		p.logger.Error("unable to accept blocks into chain", slog.String("hash", getHashStringNoErr(tip.Hash)), slog.Uint64("height", tip.Height), slog.String("chain", chain.String()), slog.String("err", err.Error()))
		return false
	}

	p.logger.Info("blocks successfully accepted into chain", slog.String("hash", getHashStringNoErr(tip.Hash)), slog.Uint64("height", tip.Height), slog.String("chain", chain.String()))
	return true
}

func (p *Processor) publishTxsToMetamorph(ctx context.Context, txs []store.TransactionBlock) {
	var publishErr error
	ctx, span := tracing.StartTracing(ctx, "publish transactions", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, publishErr)
	}()

	for _, tx := range txs {
		txBlock := &blocktx_api.TransactionBlock{
			BlockHash:       tx.BlockHash,
			BlockHeight:     tx.BlockHeight,
			TransactionHash: tx.TxHash,
			MerklePath:      tx.MerklePath,
			BlockStatus:     tx.BlockStatus,
		}

		err := p.mqClient.PublishMarshal(ctx, MinedTxsTopic, txBlock)
		if err != nil {
			p.logger.Error("failed to publish mined txs", slog.String("blockHash", getHashStringNoErr(tx.BlockHash)), slog.Uint64("height", tx.BlockHeight), slog.String("txHash", getHashStringNoErr(tx.TxHash)), slog.String("err", err.Error()))
			publishErr = err
		}
	}
}

func (p *Processor) Shutdown() {
	p.cancelAll()
	p.waitGroup.Wait()

	// unlock unprocessed blocks
	bhs, err := p.store.GetBlockHashesProcessingInProgress(context.Background(), p.hostname)
	if err != nil {
		p.logger.Error("reading unprocessing blocks on shutdown failed", slog.Any("err", err))
		return
	}

	for _, bh := range bhs {
		_, err := p.store.DelBlockProcessing(context.Background(), bh, p.hostname)
		if err != nil {
			p.logger.Error("unlocking unprocessed block on shutdown failed", slog.String("hash", bh.String()), slog.Any("err", err))
		}
	}
}
