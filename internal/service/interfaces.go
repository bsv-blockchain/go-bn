package service

import (
	"context"
	"fmt"
)

// RPC interface with a rpc server.
type RPC interface {
	Do(ctx context.Context, method string, out interface{}, args ...interface{}) error
}

type request struct {
	method string
	args   []interface{}
}

// Key returns a unique key for the request based on its method and arguments.
func (r request) Key() string {
	return fmt.Sprintf("%s|%s", r.method, r.args)
}
