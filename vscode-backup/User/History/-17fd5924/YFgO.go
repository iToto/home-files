// Package tradelogger will handle logging trades to BQ
package tradelogger

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"

	"github.com/guregu/null"
	"github.com/oklog/ulid"
)

// this must match the schema in BQ
type record struct {
	ID           string    `db:"id" json:"id"`
	Exchange     string    `db:"exchange" json:"exchange"`
	Chain        string    `db:"chain" json:"chain"`
	SignalID     string    `db:"signal_id" json:"signal_id"`
	TradeOrderID string    `db:"trade_order_id" json:"trade_order_id"` //CoinRoutes: ClientOrderID
	Quantity     string    `db:"quantity" json:"quantity"`
	Side         string    `db:"side" json:"side"`
	CurrencyPair string    `db:"currency_pair" json:"currency_pair"`
	CreatedAt    null.Time `db:"created_at" json:"created_at"`
	UpdatedAt    null.Time `db:"updated_at" json:"updated_at"`
	DeletedAt    null.Time `db:"deleted_at" json:"deleted_at"`
}

// Archive writes the given event to BigQuery.
func (a *DataLogger) Log(ctx context.Context, wl wlog.Logger, t *coinroutesapi.ClientOrderCreateResponse) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return fmt.Errorf("error generating id: %s", err)
	}

	r := &record{}

	return a.inserter.Put(ctx, r)
}
