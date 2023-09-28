package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/libsv/go-bt/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/ordishs/go-bitcoin"
	"github.com/ordishs/go-utils"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	"github.com/bitcoin-sv/arc/api"
	"github.com/bitcoin-sv/arc/api/transactionHandler"
	"github.com/bitcoin-sv/arc/metamorph/metamorph_api"
	"github.com/bitcoin-sv/arc/validator"
	defaultValidator "github.com/bitcoin-sv/arc/validator/default"
)

type ArcDefaultHandler struct {
	TransactionHandler transactionHandler.TransactionHandler
	NodePolicy         *bitcoin.Settings
	logger             utils.Logger
}

func NewDefault(logger utils.Logger, transactionHandler transactionHandler.TransactionHandler, policy *bitcoin.Settings) (api.HandlerInterface, error) {

	handler := &ArcDefaultHandler{
		TransactionHandler: transactionHandler,
		NodePolicy:         policy,
		logger:             logger,
	}

	return handler, nil
}

func (m ArcDefaultHandler) GETPolicy(ctx echo.Context) error {
	span, _ := opentracing.StartSpanFromContext(ctx.Request().Context(), "ArcDefaultHandler:GETPolicy")
	defer span.Finish()

	satoshis, bytes := calcFeesFromBSVPerKB(m.NodePolicy.MinMiningTxFee)

	return ctx.JSON(http.StatusOK, api.PolicyResponse{
		Policy: api.Policy{
			Maxscriptsizepolicy:     uint64(m.NodePolicy.MaxScriptSizePolicy),
			Maxtxsigopscountspolicy: uint64(m.NodePolicy.MaxTxSigopsCountsPolicy),
			Maxtxsizepolicy:         uint64(m.NodePolicy.MaxTxSizePolicy),
			MiningFee: api.FeeAmount{
				Bytes:    bytes,
				Satoshis: satoshis,
			},
		},
		Timestamp: time.Now().UTC(),
	})
}

func calcFeesFromBSVPerKB(feePerKB float64) (uint64, uint64) {
	bytes := uint64(1000)
	fSatoshis := feePerKB * 1e8
	satoshis := uint64(fSatoshis)

	// increment bytes and satoshis by a factor of 10 until satoshis is not a fraction
	for fSatoshis != float64(satoshis) {
		fSatoshis *= 10
		satoshis = uint64(fSatoshis)
		bytes *= 10
	}

	return satoshis, bytes
}

// POSTTransaction ...
func (m ArcDefaultHandler) POSTTransaction(ctx echo.Context, params api.POSTTransactionParams) error {
	span, tracingCtx := opentracing.StartSpanFromContext(ctx.Request().Context(), "ArcDefaultHandler:POSTTransactions")
	defer span.Finish()

	transactionOptions, err := getTransactionOptions(params)
	if err != nil {
		e := api.ErrBadRequest
		e.ExtraInfo = ptr.To(err.Error())
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return ctx.JSON(int(api.ErrStatusBadRequest), e)
	}

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		e := api.ErrBadRequest
		e.ExtraInfo = ptr.To(err.Error())
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return ctx.JSON(http.StatusBadRequest, e)
	}

	var transaction *bt.Tx
	contentType := ctx.Request().Header.Get("Content-Type")
	switch contentType {
	case "text/plain":
		if transaction, err = bt.NewTxFromString(string(body)); err != nil {
			e := api.ErrBadRequest
			e.ExtraInfo = ptr.To(err.Error())
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(err))
			return ctx.JSON(int(api.ErrStatusBadRequest), e)
		}
	case "application/json":
		var txHex string
		var txBody api.POSTTransactionJSONRequestBody
		if err = json.Unmarshal(body, &txBody); err != nil {
			e := api.ErrBadRequest
			e.ExtraInfo = ptr.To(err.Error())
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(err))
			return ctx.JSON(int(api.ErrStatusBadRequest), e)
		}
		txHex = txBody.RawTx

		if transaction, err = bt.NewTxFromString(txHex); err != nil {
			e := api.ErrBadRequest
			e.ExtraInfo = ptr.To(err.Error())
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(err))
			return ctx.JSON(int(api.ErrStatusBadRequest), e)
		}
	case "application/octet-stream":
		if transaction, err = bt.NewTxFromBytes(body); err != nil {
			e := api.ErrBadRequest
			e.ExtraInfo = ptr.To(err.Error())
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(err))
			return ctx.JSON(e.Status, e)
		}
	default:
		e := api.ErrBadRequest
		e.ExtraInfo = ptr.To(fmt.Sprintf("given content-type %s does not match any of the allowed content-types", contentType))
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return ctx.JSON(e.Status, e)
	}

	span.SetTag("txid", transaction.TxID())

	status, response, responseErr := m.processTransaction(tracingCtx, transaction, transactionOptions)
	if responseErr != nil {
		// if an error is returned, the processing failed, and we should return a 500 error
		span.LogFields(log.Error(responseErr))
		return responseErr
	}

	sizingInfo := make([][]uint64, 1)

	normalBytes, dataBytes, feeAmount := getSizings(transaction)
	sizingInfo[0] = []uint64{normalBytes, dataBytes, feeAmount}
	sizingCtx := context.WithValue(ctx.Request().Context(), api.ContextSizings, sizingInfo)
	ctx.SetRequest(ctx.Request().WithContext(sizingCtx))

	return ctx.JSON(int(status), response)
}

// GETTransactionStatus ...
func (m ArcDefaultHandler) GETTransactionStatus(ctx echo.Context, id string) error {
	span, tracingCtx := opentracing.StartSpanFromContext(ctx.Request().Context(), "ArcDefaultHandler:GETTransactionStatus")
	defer span.Finish()

	span.SetTag("txid", id)

	tx, err := m.getTransactionStatus(tracingCtx, id)
	if err != nil {
		if errors.Is(err, transactionHandler.ErrTransactionNotFound) {
			e := api.ErrNotFound
			e.Detail = err.Error()
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(err))
			return echo.NewHTTPError(http.StatusNotFound, e)
		}

		e := api.ErrGeneric
		e.ExtraInfo = ptr.To(err.Error())
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return ctx.JSON(int(api.ErrStatusGeneric), e)
	}

	if tx == nil {
		return echo.NewHTTPError(http.StatusNotFound, api.ErrNotFound)
	}

	return ctx.JSON(http.StatusOK, api.TransactionStatus{
		BlockHash:   &tx.BlockHash,
		BlockHeight: &tx.BlockHeight,
		TxStatus:    &tx.Status,
		Timestamp:   time.Now(),
		Txid:        tx.TxID,
		MerklePath:  &tx.MerklePath,
	})
}

// POSTTransactions ...
func (m ArcDefaultHandler) POSTTransactions(ctx echo.Context, params api.POSTTransactionsParams) error {
	span, tracingCtx := opentracing.StartSpanFromContext(ctx.Request().Context(), "ArcDefaultHandler:POSTTransactions")
	defer span.Finish()

	// set the globals for all transactions in this request
	transactionOptions, err := getTransactionsOptions(params)
	if err != nil {
		e := api.ErrBadRequest
		e.ExtraInfo = ptr.To(err.Error())
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return ctx.JSON(e.Status, e)
	}

	// Set the transaction reader function to read a text/plain by default.
	// If the mimetype is application/octet-stream, then we will replace this
	// function with one that reads the raw bytes in the switch statement below.
	transactionReaderFn := func(r io.Reader) (*bt.Tx, error) {
		reader := bufio.NewReader(r)
		b, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}

		return bt.NewTxFromString(string(b))
	}

	var transactions []interface{}
	sizingInfo := make([][]uint64, 0)

	var status api.StatusCode

	contentType := ctx.Request().Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		body, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			e := api.ErrBadRequest
			e.ExtraInfo = ptr.To(err.Error())
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(err))
			return ctx.JSON(e.Status, e)
		}

		var txBody api.POSTTransactionsJSONBody
		if err = json.Unmarshal(body, &txBody); err != nil {
			e := api.ErrBadRequest
			e.ExtraInfo = ptr.To(err.Error())
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(err))
			return ctx.JSON(e.Status, e)
		}

		sizingMap := make(map[string][]uint64)
		transactionInputs := make([]*bt.Tx, 0, len(txBody))
		for index, tx := range txBody {
			transaction, err := bt.NewTxFromString(tx.RawTx)
			if err != nil {
				e := api.ErrBadRequest
				e.ExtraInfo = ptr.To(err.Error())
				transactions[index] = e
				return ctx.JSON(e.Status, e)
			}
			transactionInputs = append(transactionInputs, transaction)
			normalBytes, dataBytes, feeAmount := getSizings(transaction)
			sizingMap[transaction.TxID()] = []uint64{normalBytes, dataBytes, feeAmount}
		}
		// submit for processing
		status, transactions, err = m.processTransactions(tracingCtx, transactionInputs, transactionOptions)
		if err != nil {
			e := api.ErrGeneric
			e.ExtraInfo = ptr.To(fmt.Sprintf("failed to process transactions: %v", err))
			return ctx.JSON(e.Status, e)
		}

		for _, btTx := range transactions {
			if tx, ok := btTx.(api.TransactionResponse); ok {
				sizing, found := sizingMap[tx.Txid]
				if !found {
					m.logger.Warnf("tx id %s not found in sizing map", tx.Txid)
				}

				sizingInfo = append(sizingInfo, sizing)
			}
		}
	case "application/octet-stream":
		transactionReaderFn = func(r io.Reader) (*bt.Tx, error) {
			btTx := new(bt.Tx)
			if _, err := btTx.ReadFrom(r); err != nil {
				return nil, err
			}
			return btTx, nil
		}
		fallthrough
	case "text/plain":
		reader := ctx.Request().Body
		transactionInputs := make([]*bt.Tx, 0)
		sizingsMap := make(map[string][]uint64)

		isFirstTransaction := true
		var err error

		// parse each transaction from request body and prepare
		// slice of transactions to process before submitting to metamorph
		for {
			btTx, err := transactionReaderFn(reader)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					e := api.ErrBadRequest
					e.ExtraInfo = ptr.To(err.Error())
					span.SetTag(string(ext.Error), true)
					span.LogFields(log.Error(err))
					return ctx.JSON(e.Status, e)
				}
			}

			// we reached the end of request body
			if btTx == nil {
				if isFirstTransaction {
					// no transactions found in the request body
					e := api.ErrBadRequest
					e.ExtraInfo = ptr.To("no transactions found in the request body")
					span.SetTag(string(ext.Error), true)
					span.LogFields(log.Error(err))
					return ctx.JSON(e.Status, e)
				}
				// no more transaction data found, stop the loop
				break
			}

			isFirstTransaction = false
			transactionInputs = append(transactionInputs, btTx)
			normalBytes, dataBytes, feeAmount := getSizings(btTx)
			sizingsMap[btTx.TxID()] = []uint64{normalBytes, dataBytes, feeAmount}
		}
		// process all transactions before submitting to metamorph
		status, transactions, err = m.processTransactions(tracingCtx, transactionInputs, transactionOptions)
		if err != nil {
			e := api.ErrGeneric
			e.ExtraInfo = ptr.To(fmt.Sprintf("failed to process transactions: %v", err))
			span.SetTag(string(ext.Error), true)
			span.LogFields(log.Error(err))
			return ctx.JSON(e.Status, e)
		}

		for _, btTx := range transactions {
			if tx, ok := btTx.(api.TransactionResponse); ok {
				sizing, found := sizingsMap[tx.Txid]
				if !found {
					m.logger.Warnf("tx id %s not found in sizing map", tx.Txid)
				}

				sizingInfo = append(sizingInfo, sizing)
			}
		}

	default:
		e := api.ErrBadRequest
		e.ExtraInfo = ptr.To(fmt.Sprintf("given content-type %s does not match any of the allowed content-types", contentType))
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return ctx.JSON(e.Status, e)
	}

	sizingCtx := context.WithValue(ctx.Request().Context(), api.ContextSizings, sizingInfo)
	ctx.SetRequest(ctx.Request().WithContext(sizingCtx))
	// we cannot really return any other status here
	// each transaction in the slice will have the result of the transaction submission
	return ctx.JSON(int(status), transactions)
}

func getTransactionOptions(params api.POSTTransactionParams) (*api.TransactionOptions, error) {
	return getTransactionsOptions(api.POSTTransactionsParams(params))
}

func getTransactionsOptions(params api.POSTTransactionsParams) (*api.TransactionOptions, error) {
	transactionOptions := &api.TransactionOptions{}
	if params.XCallbackUrl != nil {
		_, err := url.ParseRequestURI(*params.XCallbackUrl)

		if err != nil {
			return nil, fmt.Errorf("invalid callback URL [%w]", err)
		}

		transactionOptions.CallbackURL = *params.XCallbackUrl
		if params.XCallbackToken != nil {
			transactionOptions.CallbackToken = *params.XCallbackToken
		}
	}
	if params.XMerkleProof != nil {
		if *params.XMerkleProof == "true" || *params.XMerkleProof == "1" {
			transactionOptions.MerkleProof = true
		}
	}
	if params.XWaitForStatus != nil {
		if *params.XWaitForStatus >= 2 && *params.XWaitForStatus <= 6 {
			transactionOptions.WaitForStatus = metamorph_api.Status(*params.XWaitForStatus)
		}
	}

	return transactionOptions, nil
}

func (m ArcDefaultHandler) processTransaction(ctx context.Context, transaction *bt.Tx, transactionOptions *api.TransactionOptions) (api.StatusCode, interface{}, error) {
	span, tracingCtx := opentracing.StartSpanFromContext(ctx, "ArcDefaultHandler:processTransaction")
	defer span.Finish()

	txValidator := defaultValidator.New(m.NodePolicy)

	// the validator expects an extended transaction
	// we must enrich the transaction with the missing data
	if !txValidator.IsExtended(transaction) {
		err := m.extendTransaction(tracingCtx, transaction)
		if err != nil {
			statusCode, arcError, errHandled := m.handleError(tracingCtx, transaction, err)
			m.logger.Errorf("failed to extend transaction with ID %s, status Code: %d: %v", transaction.TxID(), statusCode, err)
			return statusCode, arcError, errHandled
		}
	}

	validateSpan, validateCtx := opentracing.StartSpanFromContext(tracingCtx, "ArcDefaultHandler:ValidateTransaction")
	if err := txValidator.ValidateTransaction(transaction); err != nil {
		validateSpan.Finish()
		statusCode, arcError, errHandled := m.handleError(validateCtx, transaction, err)
		m.logger.Errorf("failed to validate transaction with ID %s, status Code: %d: %v", transaction.TxID(), statusCode, errHandled)
		return statusCode, arcError, errHandled
	}
	validateSpan.Finish()

	tx, err := m.TransactionHandler.SubmitTransaction(tracingCtx, transaction.Bytes(), transactionOptions)
	if err != nil {
		statusCode, arcError, errHandled := m.handleError(tracingCtx, transaction, err)
		m.logger.Errorf("failed to submit transaction with ID %s, status Code: %d: %v", transaction.TxID(), statusCode, errHandled)
		return statusCode, arcError, errHandled
	}

	txID := tx.TxID
	if txID == "" {
		txID = transaction.TxID()
	}

	var extraInfo string
	if tx.ExtraInfo != "" {
		extraInfo = tx.ExtraInfo
	}

	return api.StatusOK, api.TransactionResponse{
		Status:      int(api.StatusOK),
		Title:       "OK",
		BlockHash:   &tx.BlockHash,
		BlockHeight: &tx.BlockHeight,
		TxStatus:    tx.Status,
		ExtraInfo:   &extraInfo,
		Timestamp:   time.Now(),
		Txid:        txID,
		MerklePath:  &tx.MerklePath,
	}, nil
}

// processTransactions validates all the transactions in the array and submits to metamorph for processing
func (m ArcDefaultHandler) processTransactions(ctx context.Context, transactions []*bt.Tx, transactionOptions *api.TransactionOptions) (api.StatusCode, []interface{}, error) {
	span, tracingCtx := opentracing.StartSpanFromContext(ctx, "ArcDefaultHandler:processTransactions")
	defer span.Finish()

	// validate before submitting array of transactions to metamorph
	transactionsInput := make([][]byte, 0, len(transactions))
	txErrors := make([]interface{}, 0, len(transactions))

	for _, transaction := range transactions {
		txValidator := defaultValidator.New(m.NodePolicy)

		// the validator expects an extended transaction
		// we must enrich the transaction with the missing data
		if !txValidator.IsExtended(transaction) {
			err := m.extendTransaction(tracingCtx, transaction)
			if err != nil {
				_, arcError, _ := m.handleError(tracingCtx, transaction, err)
				txErrors = append(txErrors, arcError)
				continue
			}
		}

		// validate transaction
		validateSpan, validateCtx := opentracing.StartSpanFromContext(tracingCtx, "ArcDefaultHandler:ValidateTransactions")
		if err := txValidator.ValidateTransaction(transaction); err != nil {
			validateSpan.Finish()
			_, arcError, _ := m.handleError(validateCtx, transaction, err)
			txErrors = append(txErrors, arcError)
			continue
		}
		validateSpan.Finish()
		transactionsInput = append(transactionsInput, transaction.Bytes())
	}

	// submit all the validated array of transactiosn to metamorph endpoint
	txStatuses, err := m.TransactionHandler.SubmitTransactions(tracingCtx, transactionsInput, transactionOptions)
	if err != nil {
		statusCode, arcError, errHandled := m.handleError(tracingCtx, nil, err)
		m.logger.Errorf("failed to submit %d transactions, status Code: %d: %v", len(transactions), statusCode, errHandled)
		return statusCode, []interface{}{arcError}, errHandled
	}

	// process returned transaction statuses and return to user
	var transactionOutput []interface{} = make([]interface{}, 0, len(transactions))

	for ind, tx := range txStatuses {
		transactionOutput = append(transactionOutput, api.TransactionResponse{
			Status:      int(api.StatusOK),
			Title:       "OK",
			BlockHash:   &txStatuses[ind].BlockHash,
			BlockHeight: &txStatuses[ind].BlockHeight,
			TxStatus:    tx.Status,
			ExtraInfo:   &txStatuses[ind].ExtraInfo,
			Timestamp:   time.Now(),
			Txid:        transactions[ind].TxID(),
			MerklePath:  &txStatuses[ind].MerklePath,
		})
	}

	transactionOutput = append(transactionOutput, txErrors...)

	return api.StatusOK, transactionOutput, nil
}

func (m ArcDefaultHandler) extendTransaction(ctx context.Context, transaction *bt.Tx) (err error) {
	parentTxBytes := make(map[string][]byte)
	var btParentTx *bt.Tx

	// get the missing input data for the transaction
	for _, input := range transaction.Inputs {
		parentTxIDStr := input.PreviousTxIDStr()
		b, ok := parentTxBytes[parentTxIDStr]
		if !ok {
			b, err = m.getTransaction(ctx, parentTxIDStr)
			if err != nil {
				return err
			}
			parentTxBytes[parentTxIDStr] = b
		}

		btParentTx, err = bt.NewTxFromBytes(b)
		if err != nil {
			return err
		}

		if len(btParentTx.Outputs) < int(input.PreviousTxOutIndex) {
			return fmt.Errorf("output %d not found in transaction %s", input.PreviousTxOutIndex, parentTxIDStr)
		}
		output := btParentTx.Outputs[input.PreviousTxOutIndex]

		input.PreviousTxScript = output.LockingScript
		input.PreviousTxSatoshis = output.Satoshis
	}

	return nil
}

func (m ArcDefaultHandler) getTransactionStatus(ctx context.Context, id string) (*transactionHandler.TransactionStatus, error) {
	tx, err := m.TransactionHandler.GetTransactionStatus(ctx, id)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (m ArcDefaultHandler) handleError(_ context.Context, transaction *bt.Tx, submitErr error) (api.StatusCode, interface{}, error) {
	status := api.ErrStatusGeneric
	isArcError, ok := submitErr.(*validator.Error)
	if ok {
		status = isArcError.ArcErrorStatus
	}

	if errors.Is(submitErr, transactionHandler.ErrParentTransactionNotFound) {
		status = api.ErrStatusTxFormat
	}

	// enrich the response with the error details
	arcError := api.ErrByStatus[status]
	if arcError == nil {
		arcError = &api.ErrGeneric
	}

	if transaction != nil {
		txID := transaction.TxID()
		arcError.Txid = &txID
	}
	if submitErr != nil {
		extraInfo := submitErr.Error()
		arcError.ExtraInfo = &extraInfo
	}

	return status, arcError, nil
}

// getTransaction returns the transaction with the given id from a store
func (m ArcDefaultHandler) getTransaction(ctx context.Context, inputTxID string) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ArcDefaultHandler:getTransaction")
	defer span.Finish()

	// get from our transaction handler
	txBytes, _ := m.TransactionHandler.GetTransaction(ctx, inputTxID)
	// ignore error, we try other options if we don't find it
	if txBytes != nil {
		return txBytes, nil
	}

	// get from node
	txBytes, err := getTransactionFromNode(ctx, inputTxID)
	if err != nil {
		m.logger.Warnf("failed to get transaction %s from node: %v", inputTxID, err)
	}
	// we can ignore any error here, we just check whether we have the transaction
	if txBytes != nil {
		return txBytes, nil
	}

	// get from woc
	txBytes, err = getTransactionFromWhatsOnChain(ctx, inputTxID)
	if err != nil {
		m.logger.Warnf("failed to get transaction %s from WhatsOnChain: %v", inputTxID, err)
	}
	// we can ignore any error here, we just check whether we have the transaction
	if txBytes != nil {
		return txBytes, nil
	}

	return nil, transactionHandler.ErrParentTransactionNotFound
}

func getSizings(tx *bt.Tx) (uint64, uint64, uint64) {
	var feeAmount uint64

	for _, in := range tx.Inputs {
		feeAmount += in.PreviousTxSatoshis
	}

	var dataBytes uint64
	for _, out := range tx.Outputs {
		if feeAmount >= out.Satoshis {
			feeAmount -= out.Satoshis
		} else {
			feeAmount = 0
		}

		script := *out.LockingScript
		if out.Satoshis == 0 && len(script) > 0 && (script[0] == 0x6a || (script[0] == 0x00 && script[1] == 0x6a)) {
			dataBytes += uint64(len(script))
		}
	}

	normalBytes := uint64(len(tx.Bytes())) - dataBytes

	return normalBytes, dataBytes, feeAmount
}
