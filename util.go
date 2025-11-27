package bn

import (
	"context"

	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"

	"github.com/bsv-blockchain/go-bn/models"
)

// UtilClient interfaces interaction with the util sub commands on a bitcoin node.
type UtilClient interface {
	ClearInvalidTransactions(ctx context.Context) (uint64, error)
	CreateMultiSig(ctx context.Context, n int, keys ...string) (*models.MultiSig, error)
	ValidateAddress(ctx context.Context, address string) (*models.ValidateAddress, error)
	SignMessageWithPrivKey(ctx context.Context, w *primitives.PrivateKey, msg string) (string, error)
	VerifySignedMessage(ctx context.Context, w *primitives.PrivateKey, signature, message string) (bool, error)
}

// NewUtilClient returns a client only capable of interfacing with the util sub commands on a bitcoin node.
func NewUtilClient(oo ...BitcoinClientOptFunc) UtilClient {
	return NewNodeClient(oo...)
}

// ClearInvalidTransactions clears the invalid transactions from the node's memory pool.
func (c *client) ClearInvalidTransactions(ctx context.Context) (uint64, error) {
	var resp uint64
	return resp, c.rpc.Do(ctx, "clearinvalidtransactions", &resp)
}

// CreateMultiSig creates a multi-signature address with the given number of required signatures and public keys.
func (c *client) CreateMultiSig(ctx context.Context, n int, keys ...string) (*models.MultiSig, error) {
	var resp models.MultiSig
	return &resp, c.rpc.Do(ctx, "createmultisig", &resp, n, keys)
}

// SignMessageWithPrivKey signs a message with the given private key (PrivateKey).
func (c *client) SignMessageWithPrivKey(ctx context.Context, pk *primitives.PrivateKey, msg string) (string, error) {
	var resp string
	wif := pk.Wif()
	if !c.isMainnet {
		wif = pk.WifPrefix(byte(primitives.TestNet))
	}
	return resp, c.rpc.Do(ctx, "signmessagewithprivkey", &resp, wif, msg)
}

// ValidateAddress checks if the given address is valid and returns information about it.
func (c *client) ValidateAddress(ctx context.Context, address string) (*models.ValidateAddress, error) {
	var resp models.ValidateAddress
	return &resp, c.rpc.Do(ctx, "validateaddress", &resp, address)
}

// VerifySignedMessage verifies a signed message against the given public key (PrivateKey) and message.
func (c *client) VerifySignedMessage(ctx context.Context, pk *primitives.PrivateKey, signature, message string) (bool, error) {
	var resp bool
	return resp, c.rpc.Do(ctx, "verifymessage", &resp, pk.Wif(), signature, message)
}
