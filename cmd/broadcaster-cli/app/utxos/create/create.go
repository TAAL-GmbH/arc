package create

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/bitcoin-sv/arc/broadcaster"
	"github.com/bitcoin-sv/arc/cmd/broadcaster-cli/helper"
	"github.com/bitcoin-sv/arc/lib/keyset"
	"github.com/bitcoin-sv/arc/lib/woc_client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a UTXO set",
	RunE: func(cmd *cobra.Command, args []string) error {
		outputs := viper.GetInt("outputs")
		satoshisPerOutput := viper.GetUint64("satoshis")

		isTestnet := helper.GetBool("testnet")
		callbackURL := helper.GetString("callback")
		callbackToken := helper.GetString("callbackToken")
		authorization := helper.GetString("authorization")
		keyFile := helper.GetString("keyFile")
		miningFeeSat := helper.GetInt("miningFeeSatPerKb")
		arcServer := helper.GetString("apiURL")

		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

		client, err := helper.CreateClient(&broadcaster.Auth{
			Authorization: authorization,
		}, arcServer)
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}

		var fundingKeySet *keyset.KeySet
		var receivingKeySet *keyset.KeySet

		fundingKeySet, receivingKeySet, err = helper.GetKeySetsKeyFile(keyFile)
		if err != nil {
			return fmt.Errorf("failed to get key sets: %v", err)
		}

		wocClient := woc_client.New()

		preparer, _ := broadcaster.NewRateBroadcaster(logger, client, fundingKeySet, receivingKeySet, &wocClient,
			broadcaster.WithFees(miningFeeSat),
			broadcaster.WithIsTestnet(isTestnet),
			broadcaster.WithCallback(callbackURL, callbackToken),
		)

		err = preparer.CreateUtxos(outputs, satoshisPerOutput)
		if err != nil {
			return fmt.Errorf("failed to create utxos: %v", err)
		}
		return nil
	},
}

func init() {
	var err error

	Cmd.Flags().Int("outputs", 10, "Nr of requested outputs")
	err = viper.BindPFlag("outputs", Cmd.Flags().Lookup("outputs"))
	if err != nil {
		log.Fatal(err)
	}

	Cmd.Flags().Int("satoshis", 1000, "Nr of satoshis per output outputs")
	err = viper.BindPFlag("satoshis", Cmd.Flags().Lookup("satoshis"))
	if err != nil {
		log.Fatal(err)
	}
}
