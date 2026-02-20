package config

// RPC config.
type RPC struct {
	Host     string
	Username string
	Password string //nolint:gosec // G117: Required RPC credential field, not a secret value
}
