package tradelogger

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config holds Publisher Config
type Config struct {
	GCPprojectID string `env:"DATALOG_GCP_PROJECT_ID"`
	DataSet      string `env:"DATALOG__DATASET"`
	Table        string `env:"DATALOG_TABLE"`
	// only used for local development
	ServiceAccountFile string `env:"DATALOG_SERVICE_ACCOUNT_FILE"`
}

// Validate validates that config is valid.
func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.GCPprojectID, validation.Required),
		validation.Field(&c.DataSet, validation.Required),
		validation.Field(&c.Table, validation.Required),
	)
}

// newConfigFromEnvironment returns a new config that gets initialized
// from parsing the environment variables available in the runtime.
func newConfigFromEnvironment() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %s", err)
	}
	return cfg, nil
}
