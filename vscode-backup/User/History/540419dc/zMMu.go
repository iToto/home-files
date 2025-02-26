package handler

import (
	"time"

	"github.com/guregu/null"
)

type SignalType string

const (
	Short SignalType = "short"
	Sell  SignalType = "sell"
	Null  SignalType = "null"
	Long  SignalType = "long"
	Buy   SignalType = "buy"
)

type Signal struct {
	ID        string     `db:"id" json:"id"`
	Chain     string     `db:"chain" json:"chain"`
	Signal    SignalType `db:"signal" json:"signal"`
	TradeTime time.Time  `db:"trade_time" json:"trade_time"`
	CreatedAt null.Time  `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt null.Time  `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt null.Time  `db:"deleted_at" json:"deleted_at,omitempty"`
}

// TradeLog is the struct used to log a trade in our DW
type TradeLog struct {
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
