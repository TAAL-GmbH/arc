package metamorph_test

import (
	"log/slog"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/bitcoin-sv/arc/internal/metamorph"
	"github.com/bitcoin-sv/arc/internal/metamorph/metamorph_api"
	"github.com/bitcoin-sv/arc/internal/metamorph/mocks"
	"github.com/bitcoin-sv/arc/internal/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMissingInputsZMQI(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// make and configure a mocked ZMQI
	mockedZMQI := &mocks.ZMQIMock{
		SubscribeFunc: func(s string, stringsCh chan []string) error {
			if s != "invalidtx" {
				return nil
			}
			event := make([]string, 0)
			event = append(event, "invalidtx")
			event = append(event, "7b2266726f6d426c6f636b223a2066616c73652c22736f75726365223a2022703270222c2261646472657373223a20223132372e302e302e313a3135323234222c226e6f64654964223a2037303139352c2274786964223a202234616531643230396131616165326134616137303365326164646166393133356634613162316364306438373032303033376561353631396434393566373137222c2273697a65223a203139322c22686578223a2022303130303030303030316133376534386661616133613438353966363233313631353035336534323930396532383734646537643738346666376133303430303030303030303030303030313030303030303662343833303435303232313030653530373933303264636632626336326635326265653862333031626565313666373538396336343934613638383337313065343631633439383334656266333032323037663933663361636563626566626465373965343763653334336330663763653731373866323338313166316262303566356234323763313931323630373431343132313033386662386464386534336664663664353036333339623335623563326234396435303930646538613637343239376437346463333838663766333164646631346666666666666666303133373033303030303030303030303030313937366139313435306331663839393031336364353030363231666131366533613932326334663035616261346164383861633030303030303030222c226973496e76616c6964223a20747275652c22697356616c69646174696f6e4572726f72223a2066616c73652c2269734d697373696e67496e70757473223a20747275652c226973446f75626c655370656e644465746563746564223a2066616c73652c2269734d656d706f6f6c436f6e666c6963744465746563746564223a2066616c73652c2269734e6f6e46696e616c223a2066616c73652c22697356616c69646174696f6e54696d656f75744578636565646564223a2066616c73652c2269735374616e646172645478223a20747275652c2272656a656374696f6e436f6465223a20302c2272656a656374696f6e526561736f6e223a2022222c22636f6c6c6964656457697468223a205b5d2c2272656a656374696f6e54696d65223a2022323032332d31312d31335431333a33393a32365a227d")
			event = append(event, "2459")
			stringsCh <- event
			return nil
		},
	}

	statuses := make(chan *metamorph.PeerTxMessage, 1)
	zmqURL, err := url.Parse("https://some-url.com")
	require.NoError(t, err)
	zmq := metamorph.NewZMQ(zmqURL, statuses, logger)
	err = zmq.Start(mockedZMQI)
	require.NoError(t, err)
	status := <-statuses

	assert.Equal(t, status.Status, metamorph_api.Status_SEEN_IN_ORPHAN_MEMPOOL)
}

func TestInvalidTxZMQI(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// make and configure a mocked ZMQI
	mockedZMQI := &mocks.ZMQIMock{
		SubscribeFunc: func(s string, stringsCh chan []string) error {
			if s != "hashtx2" {
				return nil
			}
			event := make([]string, 0)
			event = append(event, "hashtx2")
			event = append(event, testdata.TX1Hash.String())
			event = append(event, "2459")
			stringsCh <- event
			return nil
		},
	}

	statuses := make(chan *metamorph.PeerTxMessage, 1)
	zmqURL, err := url.Parse("https://some-url.com")
	require.NoError(t, err)
	zmq := metamorph.NewZMQ(zmqURL, statuses, logger)
	err = zmq.Start(mockedZMQI)
	require.NoError(t, err)
	var status *metamorph.PeerTxMessage
	select {
	case status = <-statuses:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for status")
	}

	assert.Equal(t, status.Status, metamorph_api.Status_ACCEPTED_BY_NETWORK)
}
