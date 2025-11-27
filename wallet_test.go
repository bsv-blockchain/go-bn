package bn_test

import (
	"context"
	"net/http"
	"testing"

	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bsv-blockchain/go-bn"
	"github.com/bsv-blockchain/go-bn/internal/config"
	"github.com/bsv-blockchain/go-bn/internal/mocks"
	"github.com/bsv-blockchain/go-bn/internal/service"
	"github.com/bsv-blockchain/go-bn/models"
	"github.com/bsv-blockchain/go-bn/testing/util"
)

// TestWalletClientBalance tests the Balance method of the WalletClient.
func TestWalletClientBalance(t *testing.T) {
	tests := map[string]struct {
		testFile   string
		opts       *models.OptsBalance
		expBalance uint64
		expRequest models.Request
		expErr     error
	}{
		"successful request": {
			testFile:   "balance",
			expBalance: 123455600,
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "getbalance",
			},
		},
		"successful request with opts": {
			testFile:   "balance",
			expBalance: 123455600,
			opts: &models.OptsBalance{
				MinimumConfirmations: 1,
				Account:              "wow",
				IncludeWatchOnly:     true,
			},
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "getbalance",
				Params:  []interface{}{"wow", 1.0, true},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "getbalance", method)
						if test.opts != nil {
							assert.Len(t, args, 3)
						} else {
							assert.Empty(t, args)
						}

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			balance, err := c.Balance(context.TODO(), test.opts)
			if test.expErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expBalance, balance)
			}
		})
	}
}

// TestWalletClientUnconfirmedBalance tests the UnconfirmedBalance method of the WalletClient.
func TestWalletClientUnconfirmedBalance(t *testing.T) {
	tests := map[string]struct {
		testFile   string
		expBalance uint64
		expRequest models.Request
		expErr     error
	}{
		"successful request": {
			testFile:   "balance",
			expBalance: 123455600,
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "getunconfirmedbalance",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "getunconfirmedbalance", method)
						assert.Empty(t, args)

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			balance, err := c.UnconfirmedBalance(context.TODO())
			if test.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expBalance, balance)
			}
		})
	}
}

// TestWalletClientReceivedByAddress tests the ReceivedByAddress method of the WalletClient.
func TestWalletClientReceivedByAddress(t *testing.T) {
	tests := map[string]struct {
		testFile    string
		address     string
		expReceived uint64
		expRequest  models.Request
		expErr      error
	}{
		"successful request": {
			testFile:    "getreceivedbyaddress",
			address:     "mzcEDt2d7QwHazAwD11WWSn8eSCb4gtpSY",
			expReceived: 25000,
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "getreceivedbyaddress",
				Params:  []interface{}{"mzcEDt2d7QwHazAwD11WWSn8eSCb4gtpSY"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "getreceivedbyaddress", method)
						assert.Len(t, args, 1)
						assert.Equal(t, test.address, args[0])

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			balance, err := c.ReceivedByAddress(context.TODO(), test.address)
			if test.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expReceived, balance)
			}
		})
	}
}

// TestWalletClientDumpPrivateKey tests the DumpPrivateKey method of the WalletClient.
func TestWalletClientDumpPrivateKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testFile   string
		address    string
		expPk      *primitives.PrivateKey
		expRequest models.Request
		expErr     error
	}{
		"success request": {
			testFile: "dumpprivkey",
			address:  "mzcEDt2d7QwHazAwD11WWSn8eSCb4gtpSY",
			expRequest: models.Request{
				ID:      "go-bn",
				JSONRpc: "1.0",
				Method:  "dumpprivkey",
				Params:  []interface{}{"mzcEDt2d7QwHazAwD11WWSn8eSCb4gtpSY"},
			},
			expPk: func() *primitives.PrivateKey {
				wifKey, err := primitives.PrivateKeyFromWif("cW9n4pgq9MqqGD8Ux5cwpgJAJ1VzPvZgskbCEmK1QmWUicejRFQn")
				assert.NoError(t, err)
				return wifKey
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "dumpprivkey", method)
						assert.Len(t, args, 1)
						assert.Equal(t, test.address, args[0])

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			wifKey, err := c.DumpPrivateKey(context.TODO(), test.address)
			if test.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expPk, wifKey)
			}
		})
	}
}

// TestWalletClientNewAddress tests the NewAddress method of the WalletClient.
func TestWalletClientNewAddress(t *testing.T) {
	tests := map[string]struct {
		testFile   string
		opts       *models.OptsNewAddress
		expRequest models.Request
		expAddress string
		expArgsLen int
		expErr     error
	}{
		"successful request without opts": {
			testFile:   "getnewaddress",
			expAddress: "mxokrvSv54CTNer4Am8WTutjqJcpGS3Txz",
			expArgsLen: 0,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "getnewaddress",
			},
		},
		"successful request with opts": {
			testFile:   "getnewaddress",
			expAddress: "mxokrvSv54CTNer4Am8WTutjqJcpGS3Txz",
			expArgsLen: 1,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "getnewaddress",
				Params:  []interface{}{"accountname"},
			},
			opts: &models.OptsNewAddress{
				Account: "accountname",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "getnewaddress", method)
						assert.Len(t, args, test.expArgsLen)
						if test.opts != nil {
							assert.Equal(t, test.opts.Account, args[0])
						}

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			address, err := c.NewAddress(context.TODO(), test.opts)
			if test.expErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expAddress, address)
			}
		})
	}
}

// TestWalletClientListAccounts tests the ListAccounts method of the WalletClient.
func TestWalletClientListAccounts(t *testing.T) {
	tests := map[string]struct {
		testFile    string
		opts        *models.OptsListAccounts
		expRequest  models.Request
		expAccounts map[string]uint64
		expArgsLen  int
		expErr      error
	}{
		"successful request without opts": {
			testFile:   "listaccounts",
			expArgsLen: 0,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "listaccounts",
			},
			expAccounts: map[string]uint64{
				"":     8567000000,
				"john": 100000,
				"bob":  100000000,
			},
		},
		"successful request with opts": {
			testFile:   "listaccounts",
			expArgsLen: 2,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "listaccounts",
				Params:  []interface{}{4.0, true},
			},
			expAccounts: map[string]uint64{
				"":     8567000000,
				"john": 100000,
				"bob":  100000000,
			},
			opts: &models.OptsListAccounts{
				MinConf:          4,
				IncludeWatchOnly: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "listaccounts", method)
						assert.Len(t, args, test.expArgsLen)
						if test.opts != nil {
							assert.Equal(t, test.opts.MinConf, args[0])
							assert.Equal(t, test.opts.IncludeWatchOnly, args[1])
						}

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			accounts, err := c.ListAccounts(context.TODO(), test.opts)
			if test.expErr != nil {
				require.Error(t, err)
				require.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expAccounts, accounts)
			}
		})
	}
}

// TestWalletClientMove tests the Move method of the WalletClient.
func TestWalletClientMove(t *testing.T) {
	tests := map[string]struct {
		testFile   string
		opts       *models.OptsMove
		expRequest models.Request
		expResult  bool
		expArgsLen int
		from       string
		to         string
		amount     uint64
		expAmount  float64
		expErr     error
	}{
		"successful request without opts": {
			testFile:   "move",
			amount:     123456789994,
			from:       "john",
			to:         "bob",
			expResult:  true,
			expArgsLen: 3,
			expAmount:  1234.56789994,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "move",
				Params:  []interface{}{"john", "bob", 1234.56789994},
			},
		},
		"successful request with opts": {
			testFile:   "move",
			amount:     123456789994,
			from:       "john",
			to:         "bob",
			expResult:  true,
			expArgsLen: 5,
			expAmount:  1234.56789994,
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "move",
				Params:  []interface{}{"john", "bob", 1234.56789994, "", "oh wow"},
			},
			opts: &models.OptsMove{
				Comment: "oh wow",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewWalletClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "move", method)
						assert.Len(t, args, test.expArgsLen)
						assert.Equal(t, test.from, args[0])
						assert.Equal(t, test.to, args[1])
						assert.InDelta(t, test.expAmount, args[2], 0.0001)
						if test.opts != nil {
							assert.Equal(t, test.opts.Comment, args[4])
						}

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			result, err := c.Move(context.TODO(), test.from, test.to, test.amount, test.opts)
			if test.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expResult, result)
			}
		})
	}
}
