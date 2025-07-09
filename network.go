package bn

import (
	"context"
	"log"

	"github.com/bsv-blockchain/go-bn/internal"
	"github.com/bsv-blockchain/go-bn/models"
)

// NodeAdd enums.
const (
	NodeAddOneTry internal.NodeAddType = "onetry"
	NodeAddRemove internal.NodeAddType = "remove"
	NodeAddAdd    internal.NodeAddType = "add"
)

// BanAction enums.
const (
	BanActionAdd    internal.BanAction = "add"
	BanActionRemove internal.BanAction = "remove"
)

// NetworkClient interfaces interaction with the network sub commands on a bitcoin node.
type NetworkClient interface {
	Ping(ctx context.Context) error
	AddNode(ctx context.Context, node string, command internal.NodeAddType) error
	ClearBanned(ctx context.Context) error
	DisconnectNode(ctx context.Context, params models.ParamsDisconnectNode) error
	NodeInfo(ctx context.Context, opts *models.OptsNodeInfo) ([]*models.NodeInfo, error)
	ConnectionCount(ctx context.Context) (uint64, error)
	ExcessiveBlock(ctx context.Context) (*models.ExcessiveBlock, error)
	NetworkTotals(ctx context.Context) (*models.NetworkTotals, error)
	NetworkInfo(ctx context.Context) (*models.NetworkInfo, error)
	PeerInfo(ctx context.Context) ([]*models.PeerInfo, error)
	ListBanned(ctx context.Context) ([]*models.BannedSubnet, error)
	SetBan(ctx context.Context, subnet string, action internal.BanAction, opts *models.OptsSetBan) error
	SetBlockMaxSize(ctx context.Context, size uint64) (string, error)
	SetExcessiveBlock(ctx context.Context, size uint64) (string, error)
	SetNetworkActive(ctx context.Context, enabled bool) error
	SetTxPropagationFrequency(ctx context.Context, frequency uint64) error
}

// NewNetworkClient returns a client only capable of interfacing with the network sub commands on a bitcoin node.
func NewNetworkClient(oo ...BitcoinClientOptFunc) NetworkClient {
	return NewNodeClient(oo...)
}

// Ping checks the connection to the node by sending a ping command.
func (c *client) Ping(ctx context.Context) error {
	return c.rpc.Do(ctx, "ping", nil)
}

// AddNode adds, removes, or tries to connect to a node.
func (c *client) AddNode(ctx context.Context, node string, command internal.NodeAddType) error {
	return c.rpc.Do(ctx, "addnode", nil, node, command)
}

// ClearBanned clears the list of banned IPs or subnets.
func (c *client) ClearBanned(ctx context.Context) error {
	return c.rpc.Do(ctx, "clearbanned", nil)
}

// DisconnectNode disconnects a node based on the provided parameters.
func (c *client) DisconnectNode(ctx context.Context, params models.ParamsDisconnectNode) error {
	return c.rpc.Do(ctx, "disconnectnode", nil, params.Args()...)
}

// NodeInfo retrieves information about nodes in the network.
func (c *client) NodeInfo(ctx context.Context, opts *models.OptsNodeInfo) ([]*models.NodeInfo, error) {
	var resp []*models.NodeInfo
	return resp, c.rpc.Do(ctx, "getaddednodeinfo", &resp, c.argsFor(opts)...)
}

// ConnectionCount returns the number of connections to the node.
func (c *client) ConnectionCount(ctx context.Context) (uint64, error) {
	var resp uint64
	return resp, c.rpc.Do(ctx, "getconnectioncount", &resp)
}

// ExcessiveBlock retrieves the current excessive block size.
func (c *client) ExcessiveBlock(ctx context.Context) (*models.ExcessiveBlock, error) {
	var resp models.ExcessiveBlock
	return &resp, c.rpc.Do(ctx, "getexcessiveblock", &resp)
}

// NetworkTotals retrieves network totals including bytes sent and received.
func (c *client) NetworkTotals(ctx context.Context) (*models.NetworkTotals, error) {
	var resp models.NetworkTotals
	return &resp, c.rpc.Do(ctx, "getnettotals", &resp)
}

// NetworkInfo retrieves general information about the network.
func (c *client) NetworkInfo(ctx context.Context) (*models.NetworkInfo, error) {
	var resp models.NetworkInfo
	return &resp, c.rpc.Do(ctx, "getnetworkinfo", &resp)
}

// PeerInfo retrieves information about connected peers.
func (c *client) PeerInfo(ctx context.Context) ([]*models.PeerInfo, error) {
	var resp []*models.PeerInfo
	return resp, c.rpc.Do(ctx, "getpeerinfo", &resp)
}

// ListBanned retrieves a list of banned IPs or subnets.
func (c *client) ListBanned(ctx context.Context) ([]*models.BannedSubnet, error) {
	var resp []*models.BannedSubnet
	return resp, c.rpc.Do(ctx, "listbanned", &resp)
}

// SetBan sets a ban on a specific IP or subnet with the specified action.
func (c *client) SetBan(ctx context.Context, subnet string, action internal.BanAction, opts *models.OptsSetBan) error {
	return c.rpc.Do(ctx, "setban", nil, c.argsFor(opts, subnet, action)...)
}

// SetBlockMaxSize sets the maximum block size for the node.
func (c *client) SetBlockMaxSize(ctx context.Context, size uint64) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "setblockmaxsize", &resp, size)
}

// SetExcessiveBlock sets the excessive block size for the node.
func (c *client) SetExcessiveBlock(ctx context.Context, size uint64) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "setexcessiveblock", &resp, size)
}

// SetNetworkActive enables or disables network activity for the node.
func (c *client) SetNetworkActive(ctx context.Context, enabled bool) error {
	return c.rpc.Do(ctx, "setnetworkactive", nil, enabled)
}

// SetTxPropagationFrequency sets the frequency of transaction propagation in the network.
func (c *client) SetTxPropagationFrequency(_ context.Context, _ uint64) error {
	log.Fatal("not implemented")
	return nil
}
