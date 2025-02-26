package wlog

import validation "github.com/go-ozzo/ozzo-validation/v4"

// Config holds the configuration options for a Logger.
type Config struct {
	// The minimum log level to print
	MinLogLevel string `env:"MIN_LOG_LEVEL"`
	// Enables pretty log printing if true
	PrettyLogs bool `env:"PRETTY_LOGS"`
}

// Validate makes sure the configuration is valid.
// It returns an error when the configuration is not valid.
func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.MinLogLevel, validation.In("debug", "info", "error")),
	)
}
