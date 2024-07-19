package broadcast

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"sync"

	"github.com/bitcoin-sv/arc/cmd/broadcaster-cli/helper"
	"github.com/bitcoin-sv/arc/internal/broadcaster"
	"github.com/bitcoin-sv/arc/internal/woc_client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:   "broadcast",
	Short: "Submit transactions to ARC",
	RunE: func(cmd *cobra.Command, args []string) error {

		rateTxsPerSecond, err := helper.GetInt("rate")
		if err != nil {
			return err
		}
		if rateTxsPerSecond == 0 {
			return errors.New("rate must be a value greater than 0")
		}

		batchSize, err := helper.GetInt("batchsize")
		if err != nil {
			return err
		}
		if batchSize == 0 {
			return errors.New("batch size must be a value greater than 0")
		}

		limit, err := helper.GetInt64("limit")
		if err != nil {
			return err
		}

		fullStatusUpdates, err := helper.GetBool("fullStatusUpdates")
		if err != nil {
			return err
		}

		isTestnet, err := helper.GetBool("testnet")
		if err != nil {
			return err
		}

		callbackURL, err := helper.GetString("callback")
		if err != nil {
			return err
		}

		callbackToken, err := helper.GetString("callbackToken")
		if err != nil {
			return err
		}

		authorization, err := helper.GetString("authorization")
		if err != nil {
			return err
		}

		keySets, err := helper.GetKeySets()
		if err != nil {
			return err
		}

		miningFeeSat, err := helper.GetInt("miningFeeSatPerKb")
		if err != nil {
			return err
		}
		if miningFeeSat == 0 {
			return errors.New("no mining fee was given")
		}

		arcServer, err := helper.GetString("apiURL")
		if err != nil {
			return err
		}
		if arcServer == "" {
			return errors.New("no api URL was given")
		}

		wocApiKey, err := helper.GetString("wocAPIKey")
		if err != nil {
			return err
		}

		logger := helper.GetLogger()

		client, err := helper.CreateClient(&broadcaster.Auth{
			Authorization: authorization,
		}, arcServer)
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}

		wg := &sync.WaitGroup{}

		_, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg.Add(1)

		wocClient := woc_client.New(woc_client.WithAuth(wocApiKey), woc_client.WithLogger(logger))

		rateBroadcaster, err := broadcaster.NewMultiKeyRateBroadcaster(logger, client, keySets, wocClient, isTestnet, broadcaster.WithFees(miningFeeSat), broadcaster.WithCallback(callbackURL, callbackToken), broadcaster.WithFullstatusUpdates(fullStatusUpdates), broadcaster.WithBatchSize(batchSize))
		if err != nil {
			return fmt.Errorf("failed to create rate broadcaster: %v", err)
		}

		err = rateBroadcaster.Start(rateTxsPerSecond, limit)
		if err != nil {
			return fmt.Errorf("failed to start rate broadcaster: %v", err)
		}

		go func() {
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt) // Signal from Ctrl+C
			<-signalChan

			cancel()
		}()

		wg.Wait()

		return nil
	},
}

func init() {
	var err error

	Cmd.Flags().Int("rate", 10, "Transactions per second to be rate broad casted per key set")
	err = viper.BindPFlag("rate", Cmd.Flags().Lookup("rate"))
	if err != nil {
		log.Fatal(err)
	}

	Cmd.Flags().Int("batchsize", 10, "Size of batches to submit transactions")
	err = viper.BindPFlag("batchsize", Cmd.Flags().Lookup("batchsize"))
	if err != nil {
		log.Fatal(err)
	}

	Cmd.Flags().Int("limit", 0, "Limit to number of transactions to be submitted after which broadcaster will stop per key set, default: no limit")
	err = viper.BindPFlag("limit", Cmd.Flags().Lookup("limit"))
	if err != nil {
		log.Fatal(err)
	}
}
