package bn_test

import (
	"context"
	"net/http"
	"testing"

	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"

	"github.com/bsv-blockchain/go-bn"
	"github.com/bsv-blockchain/go-bn/internal/config"
	"github.com/bsv-blockchain/go-bn/internal/mocks"
	"github.com/bsv-blockchain/go-bn/internal/service"
	"github.com/bsv-blockchain/go-bn/models"
	"github.com/bsv-blockchain/go-bn/testing/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUtilClientSignMessageWithPrivKey tests the SignMessageWithPrivKey method of the UtilClient.
func TestUtilClientSignMessageWithPrivKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testFile   string
		pk         *primitives.PrivateKey
		msg        string
		expRequest models.Request
		expMsg     string
		expErr     error
	}{
		"successful request": {
			testFile: "signmessagewithprivkey",
			pk: func() *primitives.PrivateKey {
				wifKey, err := primitives.PrivateKeyFromWif("cW9n4pgq9MqqGD8Ux5cwpgJAJ1VzPvZgskbCEmK1QmWUicejRFQn")
				assert.NoError(t, err)
				return wifKey
			}(),
			expRequest: models.Request{
				JSONRpc: "1.0",
				ID:      "go-bn",
				Method:  "signmessagewithprivkey",
				Params:  []interface{}{"cW9n4pgq9MqqGD8Ux5cwpgJAJ1VzPvZgskbCEmK1QmWUicejRFQn", "hello"},
			},
			msg:    "hello",
			expMsg: "IL4oekQr7n8+u6QWCvZ+jMFhRz/zMMq4wfBvXhh+eP/zVzknU+IteOsEwyGguMnN/m7BvtOdf5b9JofdI4jEktI=",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svr, cls := util.TestServer(t, &test.expRequest, test.testFile)
			defer cls()

			r := service.NewRPC(&config.RPC{
				Host: svr.URL,
			}, &http.Client{})

			c := bn.NewUtilClient(
				bn.WithHost(svr.URL),
				bn.WithCustomRPC(&mocks.MockRPC{
					DoFunc: func(ctx context.Context, method string, out interface{}, args ...interface{}) error {
						assert.Equal(t, "signmessagewithprivkey", method)
						assert.Len(t, args, 2)
						assert.Equal(t, test.pk.WifPrefix(byte(primitives.TestNet)), args[0])
						assert.Equal(t, test.msg, args[1])

						return r.Do(ctx, method, out, args...)
					},
				}),
			)

			signedMsg, err := c.SignMessageWithPrivKey(context.TODO(), test.pk, test.msg)
			if test.expErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expMsg, signedMsg)
			}
		})
	}
}
