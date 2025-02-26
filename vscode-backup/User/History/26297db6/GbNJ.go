package entities

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

type StrategyType string

const (
	LongShort   StrategyType = "longShort"
	LongNeutral StrategyType = "longNeutral"
	USDT        StrategyType = "usdt"
)

type ChainType string

const (
	BTC ChainType = "btc"
	ETH ChainType = "eth"
)

type Signal struct {
	ID        string     `db:"id" json:"id"`
	Chain     string     `db:"chain" json:"chain"`
	Signal    SignalType `db:"signal" json:"signal"`
	Strategy  string     `db:"strategy" json:"strategy"`
	TradeTime time.Time  `db:"trade_time" json:"trade_time"`
	CreatedAt null.Time  `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt null.Time  `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt null.Time  `db:"deleted_at" json:"deleted_at,omitempty"`
}

type OrderStatusType string

const (
	Open      OrderStatusType = "open"
	Closed    OrderStatusType = "closed"
	Staged    OrderStatusType = "staged"
	Cancelled OrderStatusType = "cancelled"
	Error     OrderStatusType = "error"
)

type Order struct {
	ID            string          `db:"id" json:"id"`
	ClientOrderId string          `db:"client_order_id" json:"client_order_id"`
	Strategy      string          `db:"strategy" json:"strategy"`
	Status        OrderStatusType `db:"status" json:"status"`
	CurrencyPair  string          `db:"currency_pair" json:"currency_pair"`
	AvgPrice      string          `db:"avg_price" json:"avg_price"`
	ExecutedQty   string          `db:"executed_qty" json:"executed_qty"`
	FinishedAt    null.Time       `db:"finished_at" json:"finished_at"`
	CreatedAt     null.Time       `db:"created_at" json:"created_at"`
	UpdatedAt     null.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt     null.Time       `db:"deleted_at" json:"deleted_at"`
}

type ContractPosition struct {
	Exchange         string  `json:"exchange"`
	CurrencyPair     string  `json:"currency_pair"`
	Side             string  `json:"side"`
	Quantity         float64 `json:"quantity"`
	QuantityCurrency string  `json:"quantity_currency"`
	EntryPrice       float64 `json:"entry_price"`
}

type Chain struct {
	ID         string
	Type       ChainType
	Strategies []Strategy
}

type Strategy struct {
	ID   string
	Type StrategyType
	Name string
}
