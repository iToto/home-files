package entities

import (
	"reflect"
	"time"

	"github.com/guregu/null"
)

type SignalType string

const (
	Short   SignalType = "short"
	Sell    SignalType = "sell"
	Null    SignalType = "null"
	Long    SignalType = "long"
	Buy     SignalType = "buy"
	Neutral SignalType = "neutral"
)

// IsEquivalent will compare if desired and current state are similar (Buy/Long, Sell/Short)
func (s SignalType) IsEquivalent(desired SignalType) bool {
	if s == Buy || s == Long {
		if desired == Buy || desired == Long {
			return true
		}
		return false
	}

	if s == Sell || s == Short {
		if desired == Sell || desired == Short {
			return true
		}
		return false
	}

	if s == Neutral {
		return desired == Neutral
	}

	// unknown type
	return false
}

func (s SignalType) IsValid() bool {
	switch s {
	case Buy, Sell, Null, Long, Short, Neutral:
		return true
	default:
		return false
	}
}

type StrategyType string

// FIXME: This is incorrect as USDT and USD are not strategy types
const (
	LongShort   StrategyType = "longShort"
	LongNeutral StrategyType = "longNeutral"
	USDT        StrategyType = "usdt"
	USD         StrategyType = "usd"
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

type SignalSource struct {
	ID            string    `db:"id" json:"id"`
	Type          ChainType `db:"type" json:"type"`
	IP            string    `db:"ip" json:"ip"`
	Enabled       bool      `db:"enabled" json:"enabled"`
	SignalVersion int64     `db:"signal_version" json:"signal_version"`
	CreatedAt     null.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt     null.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt     null.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type SignalLogV1 struct {
	ID            string    `bigquery:"id"`
	Chain         string    `bigquery:"chain"`
	IP            string    `bigquery:"ip"`
	Signal        string    `bigquery:"signal"`
	LastData      time.Time `bigquery:"last_data"`
	CurrentTime   time.Time `bigquery:"current_time"`
	LastTrade     string    `bigquery:"last_trade"`
	LastTradeTime time.Time `bigquery:"last_trade_time"`
	CreatedAt     time.Time `bigquery:"created_at"`
}

type SignalLogV2 struct {
	ID                  string    `bigquery:"id"`
	Chain               string    `bigquery:"chain"`
	IP                  string    `bigquery:"ip"`
	FetchResultStatus   string    `bigquery:"fetch_result_status"`
	FetchType           string    `bigquery:"fetch_type"`
	StrategyState       string    `bigquery:"strategy_state"`
	StrategyVersion     string    `bigquery:"strategy_version"`
	LastChecked         time.Time `bigquery:"last_checked"`
	LastTradeSignal     string    `bigquery:"last_trade_signal"`
	LastTradeSignalTime time.Time `bigquery:"last_trade_signal_time"`
	CreatedAt           time.Time `bigquery:"created_at"`
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
	ID            string          `db:"id" json:"id,omitempty"`
	ClientOrderId string          `db:"client_order_id" json:"client_order_id,omitempty"`
	Strategy      string          `db:"strategy" json:"strategy,omitempty"`
	Status        OrderStatusType `db:"status" json:"status,omitempty"`
	CurrencyPair  string          `db:"currency_pair" json:"currency_pair,omitempty"`
	AvgPrice      string          `db:"avg_price" json:"avg_price,omitempty"`
	ExecutedQty   string          `db:"executed_qty" json:"executed_qty,omitempty"`
	Side          string          `db:"side" json:"side,omitempty"`
	Coin          ChainType       `db:"coin" json:"coin,omitempty"`
	SignalID      string          `db:"signal_id" json:"signal_id,omitempty"`
	Signal        SignalType      `db:"signal" json:"signal,omitempty"`
	FinishedAt    null.Time       `db:"finished_at" json:"finished_at,omitempty"`
	CreatedAt     null.Time       `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt     null.Time       `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt     null.Time       `db:"deleted_at" json:"deleted_at,omitempty"`
}

type ContractPosition struct {
	Exchange         string  `json:"exchange"`
	CurrencyPair     string  `json:"currency_pair"`
	Side             string  `json:"side"`
	Quantity         float64 `json:"quantity"`
	QuantityCurrency string  `json:"quantity_currency"`
	EntryPrice       float64 `json:"entry_price"`
	UnrealizedPnl    float64 `json:"unrealized_pnl"`
}

type Chain struct {
	ID         string
	Type       ChainType
	Strategies []Strategy
}

type ExchangeType string

const (
	Binance ExchangeType = "binancefutures"
	FTX     ExchangeType = "ftx"
)

type MarginType string

const (
	USDTM MarginType = "usdt"
	USDM  MarginType = "usdm"
	CoinM MarginType = "coin"
)

type LeverageAmount int

const (
	OneX LeverageAmount = 1
	TwoX LeverageAmount = 2
)

type TradeStrategyType string

const (
	Fixed    TradeStrategyType = "fixed"
	Compound TradeStrategyType = "compound"
)

type CurrencyPairType string

const (
	ETHInversePerpetual CurrencyPairType = "ETH-USD.IPERP"
	USDTETHPerpetual    CurrencyPairType = "ETH-USDT.PERP"
	USDETHPerpetual     CurrencyPairType = "ETH-USD.PERP"
	BTCInversePerpetual CurrencyPairType = "BTC-USD.IPERP"
	USDTBTCPerpetual    CurrencyPairType = "BTC-USDT.PERP"
	USDBTCPerpetual     CurrencyPairType = "BTC-USD.PERP"
)

func (t CurrencyPairType) String() string {
	return string(t)
}

type Strategy struct {
	ID               string            `db:"id" json:"id"`
	Enabled          bool              `db:"enabled" json:"enabled"`
	UserID           string            `db:"user_id" json:"user_id"`
	SignalSourceID   string            `db:"signal_source_id" json:"signal_source_id"`
	Type             StrategyType      `db:"type" json:"type"`
	Name             string            `db:"name" json:"name"`
	Exchange         ExchangeType      `db:"exchange" json:"exchange"`
	Margin           MarginType        `db:"margin" json:"margin"`
	Leverage         LeverageAmount    `db:"leverage" json:"leverage"`
	FixedTradeAmount null.Float        `db:"fixed_trade_amount" json:"fixed_trade_amount"`
	TradeStrategy    TradeStrategyType `db:"trade_strategy" json:"trade_strategy"`
	CurrencyPair     CurrencyPairType  `db:"currency_pair" json:"currency_pair"`
	CreatedAt        null.Time         `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt        null.Time         `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt        null.Time         `db:"deleted_at" json:"deleted_at,omitempty"`
}

type SignalStrategies struct {
	Signal     *SignalSource
	Chain      ChainType
	Strategies []*Strategy
}

type User struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt null.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt null.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt null.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type OrderReportRecord struct {
	Strategy    string `json:"strategy,omitempty"`
	Coin        string `json:"coin,omitempty"`
	SignalID    string `json:"signal_id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	Direction   string `json:"direction,omitempty"`
	Signal      string `json:"signal,omitempty"`
	AvgPrice    string `json:"avg_price,omitempty"`
	ExecutedQty string `json:"executed_qty,omitempty"`
	PNL         string `json:"pnl,omitempty"`
}

func (orr *OrderReportRecord) GetHeaders() []string {
	var o OrderReportRecord
	t := reflect.TypeOf(o)

	for _, name := range columns {

	}
}
