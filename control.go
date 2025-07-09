package bn

import (
	"context"
	"time"

	"github.com/bsv-blockchain/go-bn/models"
)

// ControlClient interfaces interaction with the control sub commands on a bitcoin node.
type ControlClient interface {
	ActiveZMQNotifications(ctx context.Context) ([]*models.ZMQNotification, error)
	DumpParams(ctx context.Context) ([]string, error)
	Info(ctx context.Context) (*models.Info, error)
	MemoryInfo(ctx context.Context) (*models.MemoryInfo, error)
	Settings(ctx context.Context) (*models.Settings, error)
	Stop(ctx context.Context) error
	Uptime(ctx context.Context) (time.Duration, error)
}

// NewControlClient returns a client only capable of interfacing with the control sub commands on a bitcoin node.
func NewControlClient(oo ...BitcoinClientOptFunc) ControlClient {
	return NewNodeClient(oo...)
}

// ActiveZMQNotifications returns a list of active ZMQ notifications on the node.
func (c *client) ActiveZMQNotifications(ctx context.Context) ([]*models.ZMQNotification, error) {
	var resp []*models.ZMQNotification
	return resp, c.rpc.Do(ctx, "activezmqnotifications", &resp)
}

// DumpParams returns the current parameters of the node.
func (c *client) DumpParams(ctx context.Context) ([]string, error) {
	var resp []string
	return resp, c.rpc.Do(ctx, "dumpparameters", &resp)
}

// Info returns general information about the node.
func (c *client) Info(ctx context.Context) (*models.Info, error) {
	var resp models.Info
	return &resp, c.rpc.Do(ctx, "getinfo", &resp)
}

// MemoryInfo returns memory usage information of the node.
func (c *client) MemoryInfo(ctx context.Context) (*models.MemoryInfo, error) {
	var resp models.MemoryInfo
	return &resp, c.rpc.Do(ctx, "getmemoryinfo", &resp)
}

// Settings return the current settings of the node.
func (c *client) Settings(ctx context.Context) (*models.Settings, error) {
	var resp models.Settings
	return &resp, c.rpc.Do(ctx, "getsettings", &resp)
}

// Stop gracefully stops the node.
func (c *client) Stop(ctx context.Context) error {
	return c.rpc.Do(ctx, "stop", nil)
}

// Uptime returns the uptime of the node.
func (c *client) Uptime(ctx context.Context) (time.Duration, error) {
	var resp time.Duration
	return resp, c.rpc.Do(ctx, "uptime", &resp)
}
