package coinroutesapi

import "database/sql"

type ExchangeAccountResponse struct {
	Name                string   `json:"name"`
	Exchange            string   `json:"exchange"`
	Strategies          []string `json:"strategies"`
	Enabled             bool     `json:"enabled"`
	MakerFee            string   `json:"maker_fee"`
	TakerFee            string   `json:"taker_fee"`
	ClientID            string   `json:"client_id"`
	LastUpdatedBalances string   `json:"last_updated_balances"`
}

type CadeAlgoParams struct {
	Side                string `json:"side,omitempty"`
	Quantity            string `json:"quantity,omitempty"`
	SweepPct            string `json:"sweep_pct,omitempty"`
	Aggression          string `json:"aggression,omitempty"`
	MaxPostSize         string `json:"max_post_size,omitempty"`
	MaxOrderCount       string `json:"max_order_count,omitempty"`
	PostExchangeCount   string `json:"post_exchange_count,omitempty"`
	MarketDataExchanges string `json:"market_data_exchanges,omitempty"`
}

type CurrencyPairType string

const (
	ETHInversePerpetual CurrencyPairType = "ETH-USD.IPERP"
	BTCInversePerpetual CurrencyPairType = "BTC-USD.IPERP"
)

type OrderType string

const (
	SmartPost OrderType = "smart post"
	Sweep     OrderType = "sweep"
	SmartStop OrderType = "smart stop"
	Spread    OrderType = "spread"
	CadeOther OrderType = "cade other"
	Pov       OrderType = "pov"
)

type OrderStatusType string

const (
	Open      OrderStatusType = "open"
	Staged    OrderStatusType = "staged"
	Cancelled OrderStatusType = "cancelled"
	Error     OrderStatusType = "error"
	Closed    OrderStatusType = "closed"
)

type AgressionType string

const (
	Neutral    AgressionType = "neutral"
	Aggressive AgressionType = "aggressive"
	Passive    AgressionType = "passive"
	MarketPeg  AgressionType = "market peg"
)

type SideType string

const (
	Buy  SideType = "buy"
	Sell SideType = "sell"
	Na   SideType = "na"
)

type StrategyType string

const (
	Default    StrategyType = "default"
	YieldChain StrategyType = "yieldchain"
)

type ClientOrderCreateRequest struct {
	OrderType    OrderType
	OrderStatus  OrderStatusType
	Aggression   AgressionType
	MaxPostSize  string
	CurrencyPair CurrencyPairType
	Quantity     string
	Side         SideType
	Strategy     StrategyType
}

type ClientOrderCreateResponse struct {
	ID                           int            `json:"id,omitempty"`
	UID                          string         `json:"uid,omitempty"`
	Owner                        int            `json:"owner,omitempty"`
	CreatedAt                    string         `json:"created_at,omitempty"`
	ConfirmedAt                  string         `json:"confirmed_at,omitempty"`
	UpdatedAt                    string         `json:"updated_at,omitempty"`
	TaskUpdatedAt                string         `json:"task_updated_at,omitempty"`
	FinishedAt                   string         `json:"finished_at,omitempty"`
	OrderStatusState             string         `json:"order_status_state,omitempty"`
	ErrorText                    string         `json:"error_text,omitempty"`
	ErrorCode                    string         `json:"error_code,omitempty"`
	ExecutedQty                  string         `json:"executed_qty,omitempty"`
	AvgPrice                     string         `json:"avg_price,omitempty"`
	ExchangeAccounts             []string       `json:"exchange_accounts,omitempty"`
	NetConsideration             string         `json:"net_consideration,omitempty"`
	GrossConsideration           string         `json:"gross_consideration,omitempty"`
	ExecutedFeeCost              string         `json:"executed_fee_cost,omitempty"`
	RemainingQty                 string         `json:"remaining_qty,omitempty"`
	HasOpenChildOrders           bool           `json:"has_open_child_orders,omitempty"`
	NetExecutedQuantity          string         `json:"net_executed_quantity,omitempty"`
	ExchangeSummaryStats         []string       `json:"exchange_summary_stats,omitempty"`
	FeeSummaryStats              []string       `json:"fee_summary_stats,omitempty"`
	ParentOrderName              string         `json:"parent_order_name,omitempty"`
	TaskId                       string         `json:"task_id,omitempty"`
	ParentOrderSlug              string         `json:"parent_order_slug,omitempty"`
	PctExecuted                  string         `json:"pct_executed,omitempty"`
	Sequence                     int            `json:"sequence,omitempty"`
	TraderName                   string         `json:"trader_name,omitempty"`
	StopTriggered                string         `json:"stop_triggered,omitempty"`
	LastUpdate                   string         `json:"last_update,omitempty"`
	AvgFeePrice                  *string        `json:"avg_fee_price,omitempty"`
	LegSummaryStats              []string       `json:"leg_summary_stats,omitempty"`
	LegFeeSummaryStats           []string       `json:"leg_fee_summary_stats,omitempty"`
	LegExchangeSummaryStats      []string       `json:"leg_exchange_summary_stats,omitempty"`
	OrderEvents                  []string       `json:"order_events,omitempty"`
	CadeAlgoParams               CadeAlgoParams `json:"cade_algo_params,omitempty"`
	TcaReport                    *string        `json:"tca_report,omitempty"`
	OrderType                    string         `json:"order_type,omitempty"`
	OrderStatus                  string         `json:"order_status,omitempty"`
	Aggression                   string         `json:"aggression,omitempty"`
	MovePercentage               *string        `json:"move_percentage,omitempty"`
	PercentLimitPrice            *string        `json:"percent_limit_price,omitempty"`
	MaxPostSize                  string         `json:"max_post_size,omitempty"`
	IsTwap                       bool           `json:"is_twap,omitempty"`
	EndTime                      *string        `json:"end_time,omitempty"`
	IntervalLength               *string        `json:"interval_length,omitempty"`
	TolerancePct                 *string        `json:"tolerance_pct,omitempty"`
	SweepPct                     string         `json:"sweep_pct,omitempty"`
	RepriceSeconds               *string        `json:"reprice_seconds,omitempty"`
	RepriceRandomizationPct      *string        `json:"reprice_randomization_pct,omitempty"`
	CurrencyPair                 string         `json:"currency_pair,omitempty"`
	Quantity                     string         `json:"quantity,omitempty"`
	LimitPrice                   sql.NullString `json:"limit_price,omitempty"`
	Side                         string         `json:"side,omitempty"`
	Exchanges                    []string       `json:"exchanges,omitempty"`
	MarketDataExchanges          []string       `json:"market_data_exchanges,omitempty"`
	ClientOrderId                string         `json:"client_order_id,omitempty"`
	Strategy                     string         `json:"strategy,omitempty"`
	PostExchangeCount            int            `json:"post_exchange_count,omitempty"`
	MaxOrderCount                int            `json:"max_order_count,omitempty"`
	MarginLeverage               sql.NullString `json:"margin_leverage,omitempty"`
	MarginBorrowPct              sql.NullString `json:"margin_borrow_pct,omitempty"`
	ParentOrder                  sql.NullString `json:"parent_order,omitempty"`
	IgnoreFilterLevels           sql.NullString `json:"ignore_filter_levels,omitempty"`
	UseFundingCurrency           sql.NullBool   `json:"use_funding_currency,omitempty"`
	StopPrice                    sql.NullString `json:"stop_price,omitempty"`
	EndOffset                    sql.NullString `json:"end_offset,omitempty"`
	QuantityCurrency             string         `json:"quantity_currency,omitempty"`
	SweepTo                      sql.NullString `json:"sweep_to,omitempty"`
	Notes                        sql.NullString `json:"notes,omitempty"`
	VolumeConstraintPct          sql.NullString `json:"volume_constraint_pct,omitempty"`
	VolumeConstraintIntervalSecs sql.NullString `json:"volume_constraint_interval_secs,omitempty"`
	VolumeConstraintMaxPostSize  sql.NullString `json:"volume_constraint_max_post_size,omitempty"`
	MinPostSize                  sql.NullString `json:"min_post_size,omitempty"`
}
