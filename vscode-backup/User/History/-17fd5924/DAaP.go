// Package tradelogger will handle logging trades to BQ
package tradelogger

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"

	"github.com/oklog/ulid"
)

// this must match the schema in BQ
type record struct {
	ID           string    `bigquery:"id"`
	Exchange     string    `bigquery:"exchange"`
	Chain        string    `bigquery:"chain"`
	Strategy     string    `bigquery:"strategy"`
	SignalID     string    `bigquery:"signal_id"`
	TradeOrderID string    `bigquery:"trade_order_id"` //CoinRoutes: ClientOrderID
	Quantity     string    `bigquery:"quantity"`
	Side         string    `bigquery:"side"`
	CurrencyPair string    `bigquery:"currency_pair"`
	CreatedAt    time.Time `bigquery:"created_at"`
}

// Archive writes the given event to BigQuery.
func (a *DataLogger) Log(
	ctx context.Context,
	wl wlog.Logger,
	strategy string,
	s *Signal,
	t *coinroutesapi.ClientOrderCreateResponse,
) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return fmt.Errorf("error generating id: %s", err)
	}

	r := &record{
		ID:           id.String(),
		Exchange:     "coinroutes", // FIXME: This should be dynamic when other exchanges are offered
		Chain:        s.Chain,
		Strategy:     strategy,
		SignalID:     s.ID,
		TradeOrderID: t.ClientOrderId,
		Quantity:     t.Quantity,
		Side:         t.Side,
		CurrencyPair: t.CurrencyPair,
		CreatedAt:    time.Now(),
	}

	return a.inserter.Put(ctx, r)
}
