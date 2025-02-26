package coinroutesapi

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

type CurrencyType string

const (
	BTC CurrencyType = "btc"
	ETH CurrencyType = "eth"
)

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
	OrderType    OrderType        `json:"order_type"`
	OrderStatus  OrderStatusType  `json:"order_status"`
	Aggression   AgressionType    `json:"aggression"`
	MaxPostSize  string           `json:"max_post_size"`
	CurrencyPair CurrencyPairType `json:"currency_pair"`
	Quantity     string           `json:"quantity"`
	Side         SideType         `json:"side"`
	Strategy     StrategyType     `json:"strategy"`
}

type PositionResponse struct {
	Account            string `json:"account"`
	Exchange           string `json:"exchange"`
	PositionId         string `json:"position_id"`
	CurrencyPair       string `json:"currency_pair"`
	Side               string `json:"side"`
	Quantity           string `json:"quantity"`
	QuantityCurrency   string `json:"quantity_currency"`
	EntryPrice         string `json:"entry_price"`
	CollateralUsed     string `json:"collateral_used"`
	CollateralCurrency string `json:"collateral_currency"`
	LiquidationPrice   string `json:"liquidation_price"`
	UnrealizedPnl      string `json:"unrealized_pnl"`
}

type CurrencyBalanceResponse struct {
	Currency string
	Amount   string
	Account  string
	Exchange string
}

type ClientOrderCreateResponse struct {
	ID                           int            `json:"id"`
	UID                          string         `json:"uid"`
	Owner                        int            `json:"owner"`
	CreatedAt                    string         `json:"created_at"`
	ConfirmedAt                  string         `json:"confirmed_at"`
	UpdatedAt                    string         `json:"updated_at"`
	TaskUpdatedAt                string         `json:"task_updated_at"`
	FinishedAt                   string         `json:"finished_at"`
	OrderStatusState             string         `json:"order_status_state"`
	ErrorText                    string         `json:"error_text"`
	ErrorCode                    string         `json:"error_code"`
	ExecutedQty                  string         `json:"executed_qty"`
	AvgPrice                     string         `json:"avg_price"`
	ExchangeAccounts             []string       `json:"exchange_accounts"`
	NetConsideration             string         `json:"net_consideration"`
	GrossConsideration           string         `json:"gross_consideration"`
	ExecutedFeeCost              string         `json:"executed_fee_cost"`
	RemainingQty                 string         `json:"remaining_qty"`
	HasOpenChildOrders           bool           `json:"has_open_child_orders"`
	NetExecutedQuantity          string         `json:"net_executed_quantity"`
	ExchangeSummaryStats         []string       `json:"exchange_summary_stats"`
	FeeSummaryStats              []string       `json:"fee_summary_stats"`
	ParentOrderName              string         `json:"parent_order_name"`
	TaskId                       string         `json:"task_id"`
	ParentOrderSlug              string         `json:"parent_order_slug"`
	PctExecuted                  string         `json:"pct_executed"`
	Sequence                     int            `json:"sequence"`
	TraderName                   string         `json:"trader_name"`
	StopTriggered                string         `json:"stop_triggered"`
	LastUpdate                   string         `json:"last_update"`
	AvgFeePrice                  *string        `json:"avg_fee_price"`
	LegSummaryStats              []string       `json:"leg_summary_stats"`
	LegFeeSummaryStats           []string       `json:"leg_fee_summary_stats"`
	LegExchangeSummaryStats      []string       `json:"leg_exchange_summary_stats"`
	OrderEvents                  []string       `json:"order_events"`
	CadeAlgoParams               CadeAlgoParams `json:"cade_algo_params"`
	TcaReport                    *string        `json:"tca_report"`
	OrderType                    string         `json:"order_type"`
	OrderStatus                  string         `json:"order_status"`
	Aggression                   string         `json:"aggression"`
	MovePercentage               *string        `json:"move_percentage"`
	PercentLimitPrice            *string        `json:"percent_limit_price"`
	MaxPostSize                  string         `json:"max_post_size"`
	IsTwap                       bool           `json:"is_twap"`
	EndTime                      *string        `json:"end_time"`
	IntervalLength               *string        `json:"interval_length"`
	TolerancePct                 *string        `json:"tolerance_pct"`
	SweepPct                     string         `json:"sweep_pct"`
	RepriceSeconds               *string        `json:"reprice_seconds"`
	RepriceRandomizationPct      *string        `json:"reprice_randomization_pct"`
	CurrencyPair                 string         `json:"currency_pair"`
	Quantity                     string         `json:"quantity"`
	LimitPrice                   *string        `json:"limit_price"`
	Side                         string         `json:"side"`
	Exchanges                    []string       `json:"exchanges"`
	MarketDataExchanges          []string       `json:"market_data_exchanges"`
	ClientOrderId                string         `json:"client_order_id"`
	Strategy                     string         `json:"strategy"`
	PostExchangeCount            int            `json:"post_exchange_count"`
	MaxOrderCount                int            `json:"max_order_count"`
	MarginLeverage               *string        `json:"margin_leverage"`
	MarginBorrowPct              *string        `json:"margin_borrow_pct"`
	ParentOrder                  *string        `json:"parent_order"`
	IgnoreFilterLevels           *string        `json:"ignore_filter_levels"`
	UseFundingCurrency           bool           `json:"use_funding_currency"`
	StopPrice                    *string        `json:"stop_price"`
	EndOffset                    *string        `json:"end_offset"`
	QuantityCurrency             string         `json:"quantity_currency"`
	SweepTo                      *string        `json:"sweep_to"`
	Notes                        *string        `json:"notes"`
	VolumeConstraintPct          *string        `json:"volume_constraint_pct"`
	VolumeConstraintIntervalSecs *string        `json:"volume_constraint_interval_secs"`
	VolumeConstraintMaxPostSize  *string        `json:"volume_constraint_max_post_size"`
	MinPostSize                  *string        `json:"min_post_size"`
}
