package mcast

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/libsv/go-p2p/wire"

	"github.com/bitcoin-sv/arc/internal/blocktx/bcnet"
	"github.com/bitcoin-sv/arc/internal/blocktx/store"
	"github.com/bitcoin-sv/arc/internal/multicast"
)

var _ multicast.MessageHandlerI = (*Listner)(nil)

type Listner struct {
	hostname string

	logger    *slog.Logger
	store     store.BlocktxStore
	receiveCh chan<- *bcnet.BlockMessage

	blockGroup *multicast.Group[*bcnet.BlockMessage]
}

func NewMcastListner(l *slog.Logger, addr string, network wire.BitcoinNet, store store.BlocktxStore, receiveCh chan<- *bcnet.BlockMessage) *Listner {
	hostname, _ := os.Hostname()

	listner := Listner{
		logger:    l.With("module", "mcast-listner"),
		hostname:  hostname,
		store:     store,
		receiveCh: receiveCh,
	}

	listner.blockGroup = multicast.NewGroup[*bcnet.BlockMessage](l, &listner, addr, multicast.Read, network)
	return &listner
}

func (l *Listner) Connect() bool {
	return l.blockGroup.Connect()
}

func (l *Listner) Disconnect() {
	l.blockGroup.Disconnect()
}

// OnReceive should be fire & forget
func (l *Listner) OnReceiveFromMcast(msg wire.Message) {
	if msg.Command() == wire.CmdBlock {
		blockMsg, ok := msg.(*bcnet.BlockMessage)
		if !ok {
			l.logger.Error("Block msg receive", slog.Any("err", ErrUnableToCastWireMessage))
			return
		}

		// TODO: move it to mediator or smth
		// lock block for the current instance to process
		hash := blockMsg.Hash

		l.logger.Info("Received BLOCK msg from multicast group", slog.String("hash", hash.String()))
		processedBy, err := l.store.SetBlockProcessing(context.Background(), hash, l.hostname, lockTime, maxBlocksInProgress)
		if err != nil {
			if errors.Is(err, store.ErrBlockProcessingMaximumReached) {
				l.logger.Debug("block processing maximum reached", slog.String("hash", hash.String()), slog.String("processed_by", processedBy))
				return
			} else if errors.Is(err, store.ErrBlockProcessingInProgress) {
				l.logger.Debug("block processing already in progress", slog.String("hash", hash.String()), slog.String("processed_by", processedBy))
				return
			}

			l.logger.Error("failed to set block processing", slog.String("hash", hash.String()), slog.String("err", err.Error()))
			return
		}

		l.receiveCh <- blockMsg
	}
	// ignore other messages
}

// OnSend should be fire & forget
func (l *Listner) OnSendToMcast(_ wire.Message) {
	// ignore
}
