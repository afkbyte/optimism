package dial

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

// L1EthClientWithTimeout attempts to dial the L1 provider using the
// provided URL. If the dial doesn't complete within DefaultTimeout seconds,
// this method will return an error.

// @DEV Replace this with creating a new instance of HTTP POST "github.com/btcsuite/btcd/rpcclient" and returning it
// connCfg := &rpcclient.ConnConfig{
// 	Host:         "regtest.dctrl.wtf/",
// 	User:         "test",
// 	Pass:         "test",
// 	HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
// 	DisableTLS:   false,
// }
func L1EthClientWithTimeout(ctx context.Context, url string, disableHTTP2 bool) (
	*ethclient.Client, error) {

	ctxt, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	if strings.HasPrefix(url, "http") {
		httpClient := new(http.Client)
		if disableHTTP2 {
			log.Info("Disabled HTTP/2 support in L1 eth client")
			httpClient.Transport = &http.Transport{
				TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			}
		}
		//nolint:staticcheck // Geth v1.10.23 uses rpc.DialOptions and rpc.WithClient, but we need to update geth first. Lint is flagged because of global go workspace usage.
		rpcClient, err := rpc.DialHTTPWithClient(url, httpClient)
		if err != nil {
			return nil, err
		}

		return ethclient.NewClient(rpcClient), nil
	}

	return ethclient.DialContext(ctxt, url)
}
