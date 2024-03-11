package test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bitcoin-sv/arc/api"
	"github.com/bitcoin-sv/arc/api/handler"

	"github.com/bitcoin-sv/arc/internal/metamorph/metamorph_api"
	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/bitcoinsv/bsvutil"
	"github.com/libsv/go-bc"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/unlocker"
	"github.com/stretchr/testify/require"
)

const (
	feeSat = 10

	arcEndpoint      = "http://arc:9090/"
	v1Tx             = "v1/tx"
	v1Txs            = "v1/txs"
	arcEndpointV1Tx  = arcEndpoint + v1Tx
	arcEndpointV1Txs = arcEndpoint + v1Txs
)

type Response struct {
	BlockHash   string `json:"blockHash"`
	BlockHeight int    `json:"blockHeight"`
	ExtraInfo   string `json:"extraInfo"`
	Status      int    `json:"status"`
	Timestamp   string `json:"timestamp"`
	Title       string `json:"title"`
	TxStatus    string `json:"txStatus"`
	Txid        string `json:"txid"`
}

type TxStatusResponse struct {
	BlockHash   string      `json:"blockHash"`
	BlockHeight int         `json:"blockHeight"`
	ExtraInfo   interface{} `json:"extraInfo"` // It could be null or any type, so we use interface{}
	Timestamp   string      `json:"timestamp"`
	TxStatus    string      `json:"txStatus"`
	Txid        string      `json:"txid"`
	MerklePath  string      `json:"merklePath"`
}

func TestMain(m *testing.M) {
	info, err := bitcoind.GetInfo()
	if err != nil {
		log.Fatalf("failed to get info: %v", err)
	}

	log.Printf("current block height: %d", info.Blocks)

	os.Exit(m.Run())
}

func createTx(privateKey string, address string, utxo NodeUnspentUtxo, fee ...uint64) (*bt.Tx, error) {
	tx := bt.NewTx()

	// Add an input using the first UTXO
	utxoTxID := utxo.Txid
	utxoVout := utxo.Vout
	utxoSatoshis := uint64(utxo.Amount * 1e8) // Convert BTC to satoshis
	utxoScript := utxo.ScriptPubKey

	err := tx.From(utxoTxID, utxoVout, utxoScript, utxoSatoshis)
	if err != nil {
		return nil, fmt.Errorf("failed adding input: %v", err)
	}

	// Add an output to the address you've previously created
	recipientAddress := address

	var feeValue uint64
	if len(fee) > 0 {
		feeValue = fee[0]
	} else {
		feeValue = 20 // Set your default fee value here
	}
	amountToSend := uint64(30) - feeValue // Example value - 0.009 BTC (taking fees into account)

	recipientScript, err := bscript.NewP2PKHFromAddress(recipientAddress)
	if err != nil {
		return nil, fmt.Errorf("failed converting address to script: %v", err)
	}

	err = tx.PayTo(recipientScript, amountToSend)
	if err != nil {
		return nil, fmt.Errorf("failed adding output: %v", err)
	}

	// Sign the input

	wif, err := bsvutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode WIF: %v", err)
	}

	// Extract raw private key bytes directly from the WIF structure
	privateKeyDecoded := wif.PrivKey.Serialize()

	pk, _ := bec.PrivKeyFromBytes(bsvec.S256(), privateKeyDecoded)
	unlockerGetter := unlocker.Getter{PrivateKey: pk}
	err = tx.FillAllInputs(context.Background(), &unlockerGetter)
	if err != nil {
		return nil, fmt.Errorf("sign failed: %v", err)
	}

	return tx, nil
}

func TestBatchChainedTxs(t *testing.T) {
	tt := []struct {
		name string
	}{
		{
			name: "submit batch of chained transactions - ext format",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			address, privateKey := getNewWalletAddress(t)

			generate(t, 100)

			t.Logf("generated address: %s", address)

			sendToAddress(t, address, 0.001)

			txID := sendToAddress(t, address, 0.02)
			t.Logf("sent 0.02 BSV to: %s", txID)

			hash := generate(t, 1)
			t.Logf("generated 1 block: %s", hash)

			utxos := getUtxos(t, address)
			require.True(t, len(utxos) > 0, "No UTXOs available for the address")

			txs, err := createTxChain(privateKey, utxos[0], 30)
			require.NoError(t, err)

			arcBody := make([]api.TransactionRequest, len(txs))
			for i, tx := range txs {
				arcBody[i] = api.TransactionRequest{
					RawTx: hex.EncodeToString(tx.ExtendedBytes()),
				}
			}

			payLoad, err := json.Marshal(arcBody)
			require.NoError(t, err)

			buffer := bytes.NewBuffer(payLoad)

			// Send POST request
			req, err := http.NewRequest("POST", arcEndpointV1Txs, buffer)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}

			t.Logf("submitting batch of %d chained txs", len(txs))
			postBatchRequest(t, client, req)

			time.Sleep(1 * time.Second) // give ARC time to perform the status update on DB

			// repeat request to ensure response remains the same
			t.Logf("re-submitting batch of %d chained txs", len(txs))
			postBatchRequest(t, client, req)
		})
	}
}

func postBatchRequest(t *testing.T, client *http.Client, req *http.Request) {
	httpResp, err := client.Do(req)
	require.NoError(t, err)
	defer httpResp.Body.Close()

	require.Equal(t, http.StatusOK, httpResp.StatusCode)

	b, err := io.ReadAll(httpResp.Body)
	require.NoError(t, err)

	bodyResponse := make([]Response, 0)
	err = json.Unmarshal(b, &bodyResponse)
	require.NoError(t, err)

	for i, txResponse := range bodyResponse {
		require.NoError(t, err)
		require.Equalf(t, "SEEN_ON_NETWORK", txResponse.TxStatus, "status of tx %d in chain not as expected", i)
	}
}

func createTxChain(privateKey string, utxo0 NodeUnspentUtxo, length int) ([]*bt.Tx, error) {

	batch := make([]*bt.Tx, length)

	utxoTxID := utxo0.Txid
	utxoVout := uint32(utxo0.Vout)
	utxoSatoshis := uint64(utxo0.Amount * 1e8)
	utxoScript := utxo0.ScriptPubKey
	utxoAddress := utxo0.Address

	for i := 0; i < length; i++ {
		tx := bt.NewTx()

		err := tx.From(utxoTxID, utxoVout, utxoScript, utxoSatoshis)
		if err != nil {
			return nil, fmt.Errorf("failed adding input: %v", err)
		}

		amountToSend := utxoSatoshis - feeSat

		recipientScript, err := bscript.NewP2PKHFromAddress(utxoAddress)
		if err != nil {
			return nil, fmt.Errorf("failed converting address to script: %v", err)
		}

		err = tx.PayTo(recipientScript, amountToSend)
		if err != nil {
			return nil, fmt.Errorf("failed adding output: %v", err)
		}

		// Sign the input

		wif, err := bsvutil.DecodeWIF(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode WIF: %v", err)
		}

		// Extract raw private key bytes directly from the WIF structure
		privateKeyDecoded := wif.PrivKey.Serialize()

		pk, _ := bec.PrivKeyFromBytes(bsvec.S256(), privateKeyDecoded)
		unlockerGetter := unlocker.Getter{PrivateKey: pk}
		err = tx.FillAllInputs(context.Background(), &unlockerGetter)
		if err != nil {
			return nil, fmt.Errorf("sign failed: %v", err)
		}

		batch[i] = tx

		utxoTxID = tx.TxID()
		utxoVout = 0
		utxoSatoshis = amountToSend
		utxoScript = utxo0.ScriptPubKey
	}

	return batch, nil
}

func TestPostCallbackToken(t *testing.T) {
	tt := []struct {
		name string
	}{
		{
			name: "post transaction with callback url and token",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			address, privateKey := getNewWalletAddress(t)

			generate(t, 100)

			t.Logf("generated address: %s", address)

			sendToAddress(t, address, 0.001)

			txID := sendToAddress(t, address, 0.02)
			t.Logf("sent 0.02 BSV to: %s", txID)

			hash := generate(t, 1)
			t.Logf("generated 1 block: %s", hash)

			utxos := getUtxos(t, address)
			require.True(t, len(utxos) > 0, "No UTXOs available for the address")

			tx, err := createTx(privateKey, address, utxos[0])
			require.NoError(t, err)

			arcClient, err := api.NewClientWithResponses(arcEndpoint)
			require.NoError(t, err)

			ctx := context.Background()

			hostname, err := os.Hostname()
			require.NoError(t, err)

			waitForStatus := api.WaitForStatus(metamorph_api.Status_SEEN_ON_NETWORK)
			params := &api.POSTTransactionParams{
				XWaitForStatus: &waitForStatus,
				XCallbackUrl:   handler.PtrTo(fmt.Sprintf("http://%s:9000/callback", hostname)),
				XCallbackToken: handler.PtrTo("1234"),
			}

			arcBody := api.POSTTransactionJSONRequestBody{
				RawTx: hex.EncodeToString(tx.ExtendedBytes()),
			}

			var response *api.POSTTransactionResponse
			response, err = arcClient.POSTTransactionWithResponse(ctx, params, arcBody)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, response.StatusCode())
			require.NotNil(t, response.JSON200)
			require.Equal(t, "SEEN_ON_NETWORK", response.JSON200.TxStatus)

			callbackReceivedChan := make(chan *api.TransactionStatus, 2)
			errChan := make(chan error, 2)

			expectedAuthHeader := "Bearer 1234"
			srv := &http.Server{Addr: ":9000"}
			defer func() {
				t.Log("shutting down callback listener")
				if err = srv.Shutdown(context.TODO()); err != nil {
					t.Fatal("failed to shut down server")
				}
			}()

			iterations := 0
			http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {

				defer func() {
					err := req.Body.Close()
					if err != nil {
						t.Log("failed to close body")
					}
				}()

				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					errChan <- err
				}

				var status api.TransactionStatus
				err = json.Unmarshal(bodyBytes, &status)
				if err != nil {
					errChan <- err
				}

				if expectedAuthHeader != req.Header.Get("Authorization") {
					errChan <- fmt.Errorf("auth header %s not as expected %s", expectedAuthHeader, req.Header.Get("Authorization"))
				}

				// Let ARC send the callback 2 times. First one fails.
				if iterations == 0 {
					t.Log("callback received, responding bad request")

					err = respondToCallback(w, false)
					if err != nil {
						t.Fatalf("Failed to respond to callback: %v", err)
					}

					callbackReceivedChan <- &status

					iterations++
					return
				}

				t.Log("callback received, responding success")

				err = respondToCallback(w, true)
				if err != nil {
					t.Fatalf("Failed to respond to callback: %v", err)
				}
				callbackReceivedChan <- &status
			})

			go func(server *http.Server) {
				t.Log("starting callback server")
				err = server.ListenAndServe()
				if err != nil {
					return
				}
			}(srv)

			generate(t, 10)

			time.Sleep(5 * time.Second) // give ARC time to perform the status update on DB

			var statusResponse *api.GETTransactionStatusResponse
			statusResponse, err = arcClient.GETTransactionStatusWithResponse(ctx, response.JSON200.Txid)
			seenOnNetworkReceived := false

			for i := 0; i <= 1; i++ {
				t.Logf("callback iteration %d", i)
				select {
				case callback := <-callbackReceivedChan:
					t.Logf(*statusResponse.JSON200.TxStatus)
					t.Logf(*callback.TxStatus)
					if *callback.TxStatus == "SEEN_ON_NETWORK" {
						seenOnNetworkReceived = true
						continue
					}
					require.NotNil(t, statusResponse)
					require.NotNil(t, statusResponse.JSON200)
					require.NotNil(t, callback)
					require.Equal(t, statusResponse.JSON200.Txid, callback.Txid)
					require.Equal(t, *statusResponse.JSON200.BlockHeight, *callback.BlockHeight)
					require.Equal(t, *statusResponse.JSON200.BlockHash, *callback.BlockHash)
					require.Equal(t, "MINED", *callback.TxStatus)
					require.NotNil(t, statusResponse.JSON200.MerklePath)
					_, err = bc.NewBUMPFromStr(*statusResponse.JSON200.MerklePath)
					require.NoError(t, err)

				case err := <-errChan:
					t.Fatalf("callback received - failed to parse callback %v", err)
				case <-time.NewTicker(time.Second * 15).C:
					t.Fatal("callback not received")
				}
			}
			require.Equal(t, false, seenOnNetworkReceived)
		})
	}
}

func postTxWithHeadersChecksStatus(t *testing.T, client *api.ClientWithResponses, tx *bt.Tx, expectedStatus string, skipFeeValidation bool, skipTxValidation bool) {

	ctx := context.Background()
	waitForStatus := api.WaitForStatus(metamorph_api.Status_SEEN_ON_NETWORK)

	var skipFeeValidationPtr *bool
	if skipFeeValidation {
		skipFeeValidationPtr = handler.PtrTo(true)
	}

	var skipTxValidationPtr *bool
	if skipTxValidation {
		skipTxValidationPtr = handler.PtrTo(true)
	}
	params := &api.POSTTransactionParams{
		XWaitForStatus:     &waitForStatus,
		XSkipFeeValidation: skipFeeValidationPtr,
		XSkipTxValidation:  skipTxValidationPtr,
	}

	arcBody := api.POSTTransactionJSONRequestBody{
		RawTx: hex.EncodeToString(tx.ExtendedBytes()),
	}

	var response *api.POSTTransactionResponse
	response, err := client.POSTTransactionWithResponse(ctx, params, arcBody)
	require.NoError(t, err)
	fmt.Println("Response Transaction with Zero fee:", response)
	fmt.Println("Response Transaction with Zero fee:", response.JSON200)

	require.Equal(t, http.StatusOK, response.StatusCode())
	require.NotNil(t, response.JSON200)
	require.Equalf(t, expectedStatus, response.JSON200.TxStatus, "status of response: %s does not match expected status: %s for tx ID %s", response.JSON200.TxStatus, expectedStatus, tx.TxID())
}

func TestPostSkipFee(t *testing.T) {
	tt := []struct {
		name string
	}{
		{
			name: "post transaction with skip fee",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			address, privateKey := getNewWalletAddress(t)

			generate(t, 100)

			t.Logf("generated address: %s", address)

			sendToAddress(t, address, 0.001)

			txID := sendToAddress(t, address, 0.02)
			t.Logf("sent 0.02 BSV to: %s", txID)

			hash := generate(t, 1)
			t.Logf("generated 1 block: %s", hash)

			utxos := getUtxos(t, address)
			require.True(t, len(utxos) > 0, "No UTXOs available for the address")

			customFee := uint64(0)

			tx, err := createTx(privateKey, address, utxos[0], customFee)
			require.NoError(t, err)

			fmt.Println("Transaction with Zero fee:", tx)

			url := arcEndpoint

			arcClient, err := api.NewClientWithResponses(url)
			require.NoError(t, err)

			postTxWithHeadersChecksStatus(t, arcClient, tx, "SEEN_ON_NETWORK", true, false)

		})
	}
}

func TestPostSkipTxValidation(t *testing.T) {
	tt := []struct {
		name string
	}{
		{
			name: "post transaction with skip fee",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			address, privateKey := getNewWalletAddress(t)

			generate(t, 100)

			t.Logf("generated address: %s", address)

			sendToAddress(t, address, 0.001)

			txID := sendToAddress(t, address, 0.02)
			t.Logf("sent 0.02 BSV to: %s", txID)

			hash := generate(t, 1)
			t.Logf("generated 1 block: %s", hash)

			utxos := getUtxos(t, address)
			require.True(t, len(utxos) > 0, "No UTXOs available for the address")

			customFee := uint64(0)

			tx, err := createTx(privateKey, address, utxos[0], customFee)
			require.NoError(t, err)

			fmt.Println("Transaction with Zero fee:", tx)

			url := arcEndpoint

			arcClient, err := api.NewClientWithResponses(url)
			require.NoError(t, err)

			postTxWithHeadersChecksStatus(t, arcClient, tx, "SEEN_ON_NETWORK", false, true)

		})
	}
}

func respondToCallback(w http.ResponseWriter, success bool) error {
	resp := make(map[string]string)
	if success {
		resp["message"] = "Success"
		w.WriteHeader(http.StatusOK)
	} else {
		resp["message"] = "Bad Request"
		w.WriteHeader(http.StatusBadRequest)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		return err
	}
	return nil
}

func Test_E2E_Success(t *testing.T) {

	tx := createTxHexStringExtended(t)
	jsonPayload := fmt.Sprintf(`{"rawTx": "%s"}`, hex.EncodeToString(tx.ExtendedBytes()))

	// Send POST request
	req, err := http.NewRequest("POST", arcEndpointV1Tx, strings.NewReader(jsonPayload))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	txID := postSingleRequest(t, client, req)

	time.Sleep(1 * time.Second) // give ARC time to perform the status update on DB

	// repeat request to ensure response remains the same
	txIDRepeat := postSingleRequest(t, client, req)
	require.Equal(t, txID, txIDRepeat)

	// Check transaction status
	statusUrl := fmt.Sprintf("%s/%s", arcEndpointV1Tx, txID)
	statusResp, err := http.Get(statusUrl)
	require.NoError(t, err)
	defer statusResp.Body.Close()

	var statusResponse TxStatusResponse
	require.NoError(t, json.NewDecoder(statusResp.Body).Decode(&statusResponse))
	require.Equalf(t, "SEEN_ON_NETWORK", statusResponse.TxStatus, "Expected txStatus to be 'SEEN_ON_NETWORK' for tx id %s", txID)

	t.Logf("Transaction status: %s", statusResponse.TxStatus)

	generate(t, 10)

	time.Sleep(20 * time.Second)

	statusResp, err = http.Get(statusUrl)
	require.NoError(t, err)
	defer statusResp.Body.Close()

	require.NoError(t, json.NewDecoder(statusResp.Body).Decode(&statusResponse))

	require.Equal(t, "MINED", statusResponse.TxStatus, "Expected txStatus to be 'MINED'")
	t.Logf("Transaction status: %s", statusResponse.TxStatus)

	// Check Merkle path
	t.Logf("BUMP: %s", statusResponse.MerklePath)

	bump, err := bc.NewBUMPFromStr(statusResponse.MerklePath)
	require.NoError(t, err)

	jsonB, err := json.Marshal(bump)
	require.NoError(t, err)
	t.Logf("BUMPjson: %s", string(jsonB))

	root, err := bump.CalculateRootGivenTxid(tx.TxID())
	require.NoError(t, err)

	blockRoot := getBlockRootByHeight(t, statusResponse.BlockHeight)
	require.Equal(t, blockRoot, root)

}

func postSingleRequest(t *testing.T, client *http.Client, req *http.Request) string {
	httpResp, err := client.Do(req)
	require.NoError(t, err)
	defer httpResp.Body.Close()

	require.Equal(t, http.StatusOK, httpResp.StatusCode)

	var response Response
	require.NoError(t, json.NewDecoder(httpResp.Body).Decode(&response))
	require.Equal(t, "SEEN_ON_NETWORK", response.TxStatus)

	return response.Txid
}

func TestPostTx_Success(t *testing.T) {
	tx := createTxHexStringExtended(t) // This is a placeholder for the method to create a valid transaction string.
	jsonPayload := fmt.Sprintf(`{"rawTx": "%s"}`, hex.EncodeToString(tx.ExtendedBytes()))
	resp, err := postTx(t, jsonPayload, nil) // no extra headers
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostTx_BadRequest(t *testing.T) {
	jsonPayload := `{"rawTx": "invalidHexData"}` // intentionally malformed
	resp, err := postTx(t, jsonPayload, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected 400 Bad Request but got: %d", resp.StatusCode)
}

func TestPostTx_MalformedTransaction(t *testing.T) {
	data, err := os.ReadFile("./fixtures/malformedTxHexString.txt")
	require.NoError(t, err)

	txHexString := string(data)
	jsonPayload := fmt.Sprintf(`{"rawTx": "%s"}`, txHexString)
	resp, err := postTx(t, jsonPayload, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected 400 Bad Request but got: %d", resp.StatusCode)
}

func TestPostTx_BadRequestBodyFormat(t *testing.T) {
	improperPayload := `{"transaction": "fakeData"}`

	resp, err := postTx(t, improperPayload, nil) // Using the helper function for the single tx endpoint
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected 400 Bad Request but got: %d", resp.StatusCode)
}

func postTx(t *testing.T, jsonPayload string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", arcEndpointV1Tx, strings.NewReader(jsonPayload))
	if err != nil {
		t.Fatalf("Error creating HTTP request: %s", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	return client.Do(req)
}

func createTxHexStringExtended(t *testing.T) *bt.Tx {
	address, privateKey := getNewWalletAddress(t)

	generate(t, 100)
	t.Logf("generated address: %s", address)

	sendToAddress(t, address, 0.001)

	txID := sendToAddress(t, address, 0.02)
	t.Logf("sent 0.02 BSV to: %s", txID)

	hash := generate(t, 1)
	t.Logf("generated 1 block: %s", hash)

	utxos := getUtxos(t, address)
	require.True(t, len(utxos) > 0, "No UTXOs available for the address")

	tx, err := createTx(privateKey, address, utxos[0])
	require.NoError(t, err)

	return tx
}
