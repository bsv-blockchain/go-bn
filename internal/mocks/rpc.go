package mocks

import "context"

// MockRPC mock rpc client.
type MockRPC struct {
	DoFunc func(ctx context.Context, method string, out interface{}, args ...interface{}) error
}

// Do executes the RPC method with the provided context, method name, output interface, and arguments.
func (m *MockRPC) Do(ctx context.Context, method string, out interface{}, args ...interface{}) error {
	if m.DoFunc == nil {
		panic("DoFunc not assigned in this test")
	}

	return m.DoFunc(ctx, method, out, args...)
}
