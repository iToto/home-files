// Package signallogger will handle logging trades to BQ
package signallogger

import (
	"context"
	"time"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/signalapi"
)

// this must match the schema in BQ
type record struct {
	ID            string    `bigquery:"id"`
	Chain         string    `bigquery:"chain"`
	Signal        string    `bigquery:"signal"`
	LastData      time.Time `bigquery:"last_data"`
	CurrentTime   time.Time `bigquery:"current_time"`
	LastTrade     string    `bigquery:"last_trade"`
	LastTradeTime time.Time `bigquery:"last_trade_time"`
	SignalDelta   bool      `bigquery:"signal_delta"`
	Strategy      string    `bigquery:"strategy"`
	CreatedAt     time.Time `bigquery:"created_at"`
}

// Archive writes the given event to BigQuery.
func (a *DataLogger) Log(
	ctx context.Context,
	wl wlog.Logger,
	delta bool,
	s *signalapi.SignalResponse,
	id string,
	strat string,
) error {
	r := &record{
		ID:            id, // FIXME: We should be passing entities here
		Chain:         s.Chain,
		Signal:        string(s.Signal),
		LastData:      s.LastData,
		CurrentTime:   s.CurrentTime,
		LastTrade:     string(s.LastTrade),
		LastTradeTime: s.LastTradeTime,
		SignalDelta:   delta,
		Strategy:      strat,
		CreatedAt:     time.Now(),
	}

	return a.inserter.Put(ctx, r)
}
