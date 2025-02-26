package entities

import (
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

type SignalApiVersion string

const (
	V1 SignalApiVersion = "v1"
	V2 SignalApiVersion = "v2"
)

type SignalSource struct {
	Type          ChainType
	IP            string
	Enabled       bool
	SignalVersion SignalApiVersion
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
	ID               string
	Type             StrategyType
	Name             string
	Exchange         ExchangeType
	Margin           MarginType
	Leverage         LeverageAmount
	FixedTradeAmount float64
	TradeStrategy    TradeStrategyType
	CurrencyPair     CurrencyPairType
}
