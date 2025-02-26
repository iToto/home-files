// Package tradelogger will handle logging trades to BQ
package tradelogger

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/wingocard/braavos/internal/wlog"
	"google.golang.org/api/option"
)

// Archiver archives transactions to bigquery
type Archiver struct {
	client   *bigquery.Client
	inserter *bigquery.Inserter
}

// NewArchiver creates a new publisher adapter for transactions
func NewArchiver(wl wlog.Logger) (*Archiver, error) {
	config, err := newConfigFromEnvironment()
	if err != nil {
		return nil, err
	}
	return NewArchiverWithConfig(wl, config)
}

// NewArchiverWithConfig creates a new publisher adapter for transactions
func NewArchiverWithConfig(
	wl wlog.Logger,
	config *Config,
) (*Archiver, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}

	ctx := context.Background()
	var client *bigquery.Client
	var err error
	if config.ServiceAccountFile != "" {
		client, err = bigquery.NewClient(ctx,
			config.GCPprojectID,
			option.WithCredentialsFile(config.ServiceAccountFile))
	} else {
		client, err = bigquery.NewClient(ctx, config.GCPprojectID)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to create bigquery inserter client: %v", err)
	}

	ins := client.Dataset(config.DataSet).Table(config.Table).Inserter()

	return &Archiver{
		client:   client,
		inserter: ins,
	}, nil
}

// Close closes the underlying bigquery client
func (a *Archiver) Close() error {
	return a.client.Close()
}
