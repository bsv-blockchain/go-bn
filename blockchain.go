package bn

import (
	"context"

	"github.com/bsv-blockchain/go-bc"
	"github.com/bsv-blockchain/go-bt/v2"

	"github.com/bsv-blockchain/go-bn/models"
)

// BlockChainClient interfaces interaction with the blockchain sub commands on a bitcoin node.
type BlockChainClient interface {
	BestBlockHash(ctx context.Context) (string, error)
	BlockHex(ctx context.Context, hash string) (string, error)
	BlockHexByHeight(ctx context.Context, height int) (string, error)
	BlockDecodeHeader(ctx context.Context, hash string) (*models.BlockDecodeHeader, error)
	BlockDecodeHeaderByHeight(ctx context.Context, height int) (*models.BlockDecodeHeader, error)
	Block(ctx context.Context, hash string) (*models.Block, error)
	BlockByHeight(ctx context.Context, height int) (*models.Block, error)
	ChainInfo(ctx context.Context) (*models.ChainInfo, error)
	BlockCount(ctx context.Context) (uint32, error)
	BlockHash(ctx context.Context, height int) (string, error)
	BlockHeader(ctx context.Context, hash string) (*models.BlockHeader, error)
	BlockHeaderHex(ctx context.Context, hash string) (string, error)
	BlockStats(ctx context.Context, hash string, fields ...string) (*models.BlockStats, error)
	BlockStatsByHeight(ctx context.Context, height int, fields ...string) (*models.BlockStats, error)
	ChainTips(ctx context.Context) ([]*models.ChainTip, error)
	ChainTxStats(ctx context.Context, opts *models.OptsChainTxStats) (*models.ChainTxStats, error)
	Difficulty(ctx context.Context) (float64, error)
	InvalidateBlock(ctx context.Context, blockHash string) error
	MerkleProof(ctx context.Context, blockHash, txID string, opts *models.OptsMerkleProof) (*bc.MerkleProof, error)
	LegacyMerkleProof(ctx context.Context, txID string,
		opts *models.OptsLegacyMerkleProof) (*models.LegacyMerkleProof, error)
	RawMempool(ctx context.Context) (models.MempoolTxs, error)
	RawMempoolIDs(ctx context.Context) ([]string, error)
	RawNonFinalMempool(ctx context.Context) ([]string, error)
	MempoolEntry(ctx context.Context, txID string) (*models.MempoolEntry, error)
	MempoolAncestors(ctx context.Context, txID string) (models.MempoolTxs, error)
	MempoolAncestorIDs(ctx context.Context, txID string) ([]string, error)
	MempoolDescendants(ctx context.Context, txID string) (models.MempoolTxs, error)
	MempoolDescendantIDs(ctx context.Context, txID string) ([]string, error)
	Output(ctx context.Context, txID string, n int, opts *models.OptsOutput) (*models.Output, error)
	OutputSetInfo(ctx context.Context) (*models.OutputSetInfo, error)
	PreciousBlock(ctx context.Context, blockHash string) error
	PruneChain(ctx context.Context, height int) (uint32, error)
	CheckJournal(ctx context.Context) (*models.JournalStatus, error)
	RebuildJournal(ctx context.Context) error
	VerifyChain(ctx context.Context) (bool, error)
	Generate(ctx context.Context, n int, opts *models.OptsGenerate) ([]string, error)
	GenerateToAddress(ctx context.Context, n int, addr string, opts *models.OptsGenerate) ([]string, error)
}

// NewBlockChainClient returns a client only capable of interfacing with the blockchain sub commands on a bitcoin node.
func NewBlockChainClient(oo ...BitcoinClientOptFunc) BlockChainClient {
	return NewNodeClient(oo...)
}

// BestBlockHash returns the hash of the best block in the longest blockchain.
func (c *client) BestBlockHash(ctx context.Context) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getbestblockhash", &resp)
}

// BlockHex returns the raw hex representation of a block given its hash.
func (c *client) BlockHex(ctx context.Context, hash string) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getblock", &resp, hash, models.VerbosityRawBlock)
}

// BlockHexByHeight returns the raw hex representation of a block given its height.
func (c *client) BlockHexByHeight(ctx context.Context, height int) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getblockbyheight", &resp, height, models.VerbosityRawBlock)
}

// BlockDecodeHeader returns the decoded block header for a given block hash.
func (c *client) BlockDecodeHeader(ctx context.Context, hash string) (*models.BlockDecodeHeader, error) {
	resp := models.BlockDecodeHeader{BlockHeader: models.BlockHeader{BlockHeader: &bc.BlockHeader{}}}
	return &resp, c.rpc.Do(ctx, "getblock", &resp, hash, models.VerbosityDecodeHeader)
}

// BlockDecodeHeaderByHeight returns the decoded block header for a given block height.
func (c *client) BlockDecodeHeaderByHeight(ctx context.Context, height int) (*models.BlockDecodeHeader, error) {
	resp := models.BlockDecodeHeader{BlockHeader: models.BlockHeader{BlockHeader: &bc.BlockHeader{}}}
	return &resp, c.rpc.Do(ctx, "getblockbyheight", &resp, height, models.VerbosityDecodeHeader)
}

// Block returns the block data for a given block hash.
func (c *client) Block(ctx context.Context, hash string) (*models.Block, error) {
	resp := models.Block{BlockHeader: models.BlockHeader{BlockHeader: &bc.BlockHeader{}}}
	return &resp, c.rpc.Do(ctx, "getblock", &resp, hash, models.VerbosityDecodeTransactions)
}

// BlockByHeight returns the block data for a given block height.
func (c *client) BlockByHeight(ctx context.Context, height int) (*models.Block, error) {
	resp := models.Block{BlockHeader: models.BlockHeader{BlockHeader: &bc.BlockHeader{}}}
	return &resp, c.rpc.Do(ctx, "getblockbyheight", &resp, height, models.VerbosityDecodeTransactions)
}

// ChainInfo returns information about the current state of the blockchain.
func (c *client) ChainInfo(ctx context.Context) (*models.ChainInfo, error) {
	var resp models.ChainInfo
	return &resp, c.rpc.Do(ctx, "getblockchaininfo", &resp)
}

// BlockCount returns the number of blocks in the longest blockchain.
func (c *client) BlockCount(ctx context.Context) (uint32, error) {
	var resp uint32
	return resp, c.rpc.Do(ctx, "getblockcount", &resp)
}

// BlockHash returns the hash of the block at a given height.
func (c *client) BlockHash(ctx context.Context, height int) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getblockhash", &resp, height)
}

// BlockHeader returns the block header for a given block hash.
func (c *client) BlockHeader(ctx context.Context, hash string) (*models.BlockHeader, error) {
	resp := models.BlockHeader{BlockHeader: &bc.BlockHeader{}}
	return &resp, c.rpc.Do(ctx, "getblockheader", &resp, hash, true)
}

// BlockHeaderHex returns the raw hex representation of a block header given its hash.
func (c *client) BlockHeaderHex(ctx context.Context, hash string) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getblockheader", &resp, hash, false)
}

// BlockStats returns statistics about a block given its hash.
func (c *client) BlockStats(ctx context.Context, hash string, fields ...string) (*models.BlockStats, error) {
	var resp models.BlockStats
	return &resp, c.rpc.Do(ctx, "getblockstats", &resp, hash, fields)
}

// BlockStatsByHeight returns statistics about a block given its height.
func (c *client) BlockStatsByHeight(ctx context.Context, height int, fields ...string) (*models.BlockStats, error) {
	var resp models.BlockStats
	return &resp, c.rpc.Do(ctx, "getblockstatsbyheight", &resp, height, fields)
}

// ChainTips returns information about the current chain tips.
func (c *client) ChainTips(ctx context.Context) ([]*models.ChainTip, error) {
	var resp []*models.ChainTip
	return resp, c.rpc.Do(ctx, "getchaintips", &resp)
}

// ChainTxStats returns statistics about the transactions in the blockchain.
func (c *client) ChainTxStats(ctx context.Context, opts *models.OptsChainTxStats) (*models.ChainTxStats, error) {
	var resp models.ChainTxStats
	return &resp, c.rpc.Do(ctx, "getchaintxstats", &resp, c.argsFor(opts)...)
}

// Difficulty returns the current difficulty of the blockchain.
func (c *client) Difficulty(ctx context.Context) (float64, error) {
	var resp float64
	return resp, c.rpc.Do(ctx, "getdifficulty", &resp)
}

// InvalidateBlock invalidates a block in the blockchain, effectively removing it from the chain.
func (c *client) InvalidateBlock(ctx context.Context, blockHash string) error {
	return c.rpc.Do(ctx, "invalidateblock", nil, blockHash)
}

// MempoolEntry returns detailed information about a transaction in the mempool.
func (c *client) MempoolEntry(ctx context.Context, txID string) (*models.MempoolEntry, error) {
	var resp models.MempoolEntry
	return &resp, c.rpc.Do(ctx, "getmempoolentry", &resp, txID)
}

// RawMempool returns the raw mempool transactions.
func (c *client) RawMempool(ctx context.Context) (models.MempoolTxs, error) {
	var resp models.MempoolTxs
	return resp, c.rpc.Do(ctx, "getrawmempool", &resp, true)
}

// RawMempoolIDs returns the IDs of the transactions in the raw mempool.
func (c *client) RawMempoolIDs(ctx context.Context) ([]string, error) {
	var resp []string
	return resp, c.rpc.Do(ctx, "getrawmempool", &resp, false)
}

// RawNonFinalMempool returns the non-final transactions in the raw mempool.
func (c *client) RawNonFinalMempool(ctx context.Context) ([]string, error) {
	var resp []string
	return resp, c.rpc.Do(ctx, "getrawnonfinalmempool", &resp)
}

// MempoolAncestors returns the ancestor transactions of a given transaction in the mempool.
func (c *client) MempoolAncestors(ctx context.Context, txID string) (models.MempoolTxs, error) {
	var resp models.MempoolTxs
	return resp, c.rpc.Do(ctx, "getmempoolancestors", &resp, txID, true)
}

// MempoolAncestorIDs returns the IDs of the ancestor transactions of a given transaction in the mempool.
func (c *client) MempoolAncestorIDs(ctx context.Context, txID string) ([]string, error) {
	var resp []string
	return resp, c.rpc.Do(ctx, "getmempoolancestors", &resp, txID, false)
}

// MempoolDescendants returns the descendant transactions of a given transaction in the mempool.
func (c *client) MempoolDescendants(ctx context.Context, txID string) (models.MempoolTxs, error) {
	var resp models.MempoolTxs
	return resp, c.rpc.Do(ctx, "getmempooldescendants", &resp, txID, true)
}

// MempoolDescendantIDs returns the IDs of the descendant transactions of a given transaction in the mempool.
func (c *client) MempoolDescendantIDs(ctx context.Context, txID string) ([]string, error) {
	var resp []string
	return resp, c.rpc.Do(ctx, "getmempooldescendants", &resp, txID, false)
}

// MerkleProof returns the merkle proof for a transaction in a block.
func (c *client) MerkleProof(ctx context.Context, blockHash, txID string,
	opts *models.OptsMerkleProof,
) (*bc.MerkleProof, error) {
	var resp bc.MerkleProof
	return &resp, c.rpc.Do(ctx, "getmerkleproof2", &resp, c.argsFor(opts, blockHash, txID)...)
}

// LegacyMerkleProof returns the legacy merkle proof for a transaction.
func (c *client) LegacyMerkleProof(ctx context.Context, txID string,
	opts *models.OptsLegacyMerkleProof,
) (*models.LegacyMerkleProof, error) {
	var resp models.LegacyMerkleProof
	return &resp, c.rpc.Do(ctx, "getmerkleproof", &resp, c.argsFor(opts, txID)...)
}

// Output returns the output details for a specific transaction ID and output index.
func (c *client) Output(ctx context.Context, txID string, n int, opts *models.OptsOutput) (*models.Output, error) {
	resp := models.Output{Output: &bt.Output{}}
	return &resp, c.rpc.Do(ctx, "gettxout", &resp, c.argsFor(opts, txID, n)...)
}

// OutputSetInfo returns information about the current output set.
func (c *client) OutputSetInfo(ctx context.Context) (*models.OutputSetInfo, error) {
	var resp models.OutputSetInfo
	return &resp, c.rpc.Do(ctx, "gettxoutsetinfo", &resp)
}

// PreciousBlock marks a block as precious, preventing it from being pruned.
func (c *client) PreciousBlock(ctx context.Context, blockHash string) error {
	return c.rpc.Do(ctx, "preciousblock", nil, blockHash)
}

// PruneChain prunes the blockchain up to a specified height, removing old blocks.
func (c *client) PruneChain(ctx context.Context, height int) (uint32, error) {
	var resp uint32
	return resp, c.rpc.Do(ctx, "pruneblockchain", &resp, height)
}

// CheckJournal checks the status of the journal, which is used for transaction recovery.
func (c *client) CheckJournal(ctx context.Context) (*models.JournalStatus, error) {
	var resp models.JournalStatus
	return &resp, c.rpc.Do(ctx, "checkjournal", &resp)
}

// RebuildJournal rebuilds the journal, which is used for transaction recovery.
func (c *client) RebuildJournal(ctx context.Context) error {
	return c.rpc.Do(ctx, "rebuildjournal", nil)
}

// VerifyChain verifies the integrity of the blockchain, checking for any inconsistencies.
func (c *client) VerifyChain(ctx context.Context) (bool, error) {
	var resp bool
	return resp, c.rpc.Do(ctx, "verifychain", &resp)
}

// Generate generates a specified number of blocks with optional parameters.
func (c *client) Generate(ctx context.Context, n int, opts *models.OptsGenerate) ([]string, error) {
	var resp []string
	return resp, c.rpc.Do(ctx, "generate", &resp, c.argsFor(opts, n)...)
}

// GenerateToAddress generates a specified number of blocks and sends the rewards to a specified address.
func (c *client) GenerateToAddress(ctx context.Context, n int, addr string,
	opts *models.OptsGenerate,
) ([]string, error) {
	var resp []string
	return resp, c.rpc.Do(ctx, "generatetoaddress", &resp, c.argsFor(opts, n, addr)...)
}
