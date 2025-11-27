package bn

import (
	"context"

	"github.com/bsv-blockchain/go-bt/v2"

	imodels "github.com/bsv-blockchain/go-bn/internal/models"
	"github.com/bsv-blockchain/go-bn/models"
)

// TransactionClient interfaces interaction with the transaction sub commands on a bitcoin node.
type TransactionClient interface {
	AddToConfiscationTransactionWhitelist(ctx context.Context, funds []models.ConfiscationTransactionDetails) (*models.AddToConfiscationTransactionWhitelistResponse, error)
	AddToConsensusBlacklist(ctx context.Context, funds []models.Fund) (*models.AddToConsensusBlacklistResponse, error)
	CreateRawTransaction(ctx context.Context, utxos bt.UTXOs, params models.ParamsCreateRawTransaction) (*bt.Tx, error)
	FundRawTransaction(ctx context.Context, tx *bt.Tx,
		opts *models.OptsFundRawTransaction) (*models.FundRawTransaction, error)
	RawTransaction(ctx context.Context, txID string) (*bt.Tx, error)
	SignRawTransaction(ctx context.Context, tx *bt.Tx,
		opts *models.OptsSignRawTransaction) (*models.SignedRawTransaction, error)
	SendRawTransaction(ctx context.Context, tx *bt.Tx, opts *models.OptsSendRawTransaction) (string, error)
	SendRawTransactions(ctx context.Context,
		params ...models.ParamsSendRawTransactions) (*models.SendRawTransactionsResponse, error)
}

// NewTransactionClient returns a client only capable of interfacing with the transaction sub commands
// on a bitcoin node.
func NewTransactionClient(oo ...BitcoinClientOptFunc) TransactionClient {
	return NewNodeClient(oo...)
}

// CreateRawTransaction creates a raw transaction from the given UTXOs and parameters.
func (c *client) CreateRawTransaction(ctx context.Context, utxos bt.UTXOs,
	params models.ParamsCreateRawTransaction,
) (*bt.Tx, error) {
	params.SetIsMainnet(c.isMainnet)
	var resp string
	if err := c.rpc.Do(ctx, "createrawtransaction", &resp, c.argsFor(&params, utxos.NodeJSON())...); err != nil {
		return nil, err
	}
	return bt.NewTxFromString(resp)
}

// FundRawTransaction funds a raw transaction with the given options.
func (c *client) FundRawTransaction(ctx context.Context, tx *bt.Tx,
	opts *models.OptsFundRawTransaction,
) (*models.FundRawTransaction, error) {
	resp := imodels.InternalFundRawTransaction{FundRawTransaction: &models.FundRawTransaction{}}
	return resp.FundRawTransaction, c.rpc.Do(ctx, "fundrawtransaction", &resp, c.argsFor(opts, tx.String())...)
}

// RawTransaction retrieves a raw transaction by its ID.
func (c *client) RawTransaction(ctx context.Context, txID string) (*bt.Tx, error) {
	var resp bt.Tx
	return &resp, c.rpc.Do(ctx, "getrawtransaction", &resp, txID, true)
}

// SignRawTransaction signs a raw transaction with the given options.
func (c *client) SignRawTransaction(ctx context.Context, tx *bt.Tx,
	opts *models.OptsSignRawTransaction,
) (*models.SignedRawTransaction, error) {
	var resp imodels.InternalSignRawTransaction
	return resp.SignedRawTransaction, c.rpc.Do(ctx, "signrawtransaction", &resp, c.argsFor(opts, tx.String())...)
}

// SendRawTransaction sends a raw transaction to the network and returns the transaction ID.
func (c *client) SendRawTransaction(ctx context.Context, tx *bt.Tx,
	opts *models.OptsSendRawTransaction,
) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "sendrawtransaction", &resp, c.argsFor(opts, tx.String())...)
}

// SendRawTransactions sends multiple raw transactions and returns their responses.
func (c *client) SendRawTransactions(ctx context.Context,
	params ...models.ParamsSendRawTransactions,
) (*models.SendRawTransactionsResponse, error) {
	var resp models.SendRawTransactionsResponse
	return &resp, c.rpc.Do(ctx, "sendrawtransactions", &resp, params)
}

// AddToConsensusBlacklist adds funds to the consensus blacklist.
func (c *client) AddToConsensusBlacklist(ctx context.Context, funds []models.Fund) (*models.AddToConsensusBlacklistResponse, error) {
	var resp models.AddToConsensusBlacklistResponse
	req := models.AddToConsensusBlacklistArgs{Funds: funds}
	return &resp, c.rpc.Do(ctx, "addToConsensusBlacklist", &resp, req)
}

// AddToConfiscationTransactionWhitelist adds confiscation transactions to the whitelist.
func (c *client) AddToConfiscationTransactionWhitelist(ctx context.Context, confiscationTransactions []models.ConfiscationTransactionDetails) (*models.AddToConfiscationTransactionWhitelistResponse, error) {
	var resp models.AddToConfiscationTransactionWhitelistResponse
	req := models.AddToConfiscationTxIdWhitelistArgs{
		ConfiscationTransactions: confiscationTransactions,
	}
	return &resp, c.rpc.Do(ctx, "addToConfiscationTxidWhitelist", &resp, req)
}
