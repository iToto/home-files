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
	BTC  CurrencyType = "btc"
	ETH  CurrencyType = "eth"
	USDT CurrencyType = "usdt"
	USD  CurrencyType = "usd"
	USDC CurrencyType = "usdc"
)

type CurrencyPairType string

const (
	ETHInversePerpetual CurrencyPairType = "ETH-USD.IPERP"
	USDTETHPerpetual    CurrencyPairType = "ETH-USDT.PERP"
	USDETHPerpetual     CurrencyPairType = "ETH-USD.PERP"
	BTCInversePerpetual CurrencyPairType = "BTC-USD.IPERP"
	USDTBTCPerpetual    CurrencyPairType = "BTC-USDT.PERP"
	USDBTCPerpetual     CurrencyPairType = "BTC-USD.PERP"
	USDTBTC             CurrencyPairType = "BTC-USDT"
	USDTETH             CurrencyPairType = "ETH-USDT"
)

func (t CurrencyPairType) String() string {
	return string(t)
}

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

// FIXME: This should be harmonized with entities.SignalType
type SideType string

const (
	Buy   SideType = "buy"
	Sell  SideType = "sell"
	Na    SideType = "na"
	Long  SideType = "long"
	Short SideType = "short"
	Neut  SideType = "neutral"
)

// IsEquivalent will compare if desired and current state are similar (Buy/Long, Sell/Short, Neut)
func (s SideType) IsEquivalent(desired SideType) bool {
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

	if s == Neut {
		if desired == Neut {
			return true
		}
	}

	// unknown type
	return false
}

func (s SideType) IsValid() bool {
	switch s {
	case Buy, Sell, Na, Long, Short, Neut:
		return true
	default:
		return false
	}
}

func (s SideType) GetInverseSide() SideType {
	switch s {
	case Buy:
		return Sell
	case Long:
		return Sell
	case Sell:
		return Buy
	case Short:
		return Buy
	default:
		return Na
	}
}

type ClientOrderCreateRequest struct {
	OrderType   OrderType       `json:"order_type"`
	OrderStatus OrderStatusType `json:"order_status"`
	Aggression  AgressionType   `json:"aggression"`
	// MaxPostSize  string           `json:"max_post_size"`
	CurrencyPair       CurrencyPairType `json:"currency_pair"`
	Quantity           string           `json:"quantity"`
	Side               SideType         `json:"side"`
	Strategy           string           `json:"strategy"`
	UseFundingCurrency bool             `json:"use_funding_currency"`
	EndOffset          int              `json:"end_offset"`
	IsTwap             bool             `json:"is_twap"`
	IntervalLength     int              `json:"interval_length"`
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
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
	Account  string `json:"account"`
	Exchange string `json:"exchange"`
	USDValue string `json:"usd_value"`
	USDPrice string `json:"usd_price"`
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

type ClientOrderGetResponse struct {
	ClientOrderId       string  `json:"client_order_id"`
	Strategy            string  `json:"strategy"`
	OrderStatus         string  `json:"order_status"`
	CurrencyPair        string  `json:"currency_pair"`
	AvgPrice            string  `json:"avg_price"`
	ExecutedQty         string  `json:"executed_qty"`
	ExecutedFeeCost     string  `json:"executed_fee_cost"`
	NetExecutedQuantity string  `json:"net_executed_quantity"`
	RemainingQty        float64 `json:"remaining_qty"`
	PctExecuted         float64 `json:"pct_executed"`
	TraderName          string  `json:"trader_name"`
	QuantityCurrency    string  `json:"quantity_currency"`
	FinishedAt          string  `json:"finished_at"`
}

var (
	sampleOrderCreateResponse []byte = []byte(`{
		"id": 60,
		"uid": "5d218c7f-734c-4aaa-b0f7-fe4db4bd37df",
		"owner": 9,
		"created_at": "2022-06-13T16:38:01.486415Z",
		"confirmed_at": null,
		"updated_at": "2022-06-13T16:38:02.832911Z",
		"task_updated_at": "2022-06-13T16:38:02.832911Z",
		"finished_at": "2022-06-13T16:38:02.832911Z",
		"order_status_state": "Closed",
		"error_text": null,
		"error_code": "dust",
		"executed_qty": "1020.00000000",
		"avg_price": "1251.08000000",
		"exchange_accounts": [
			"4f9f67b4-88b5-4c09-a15d-8fec268f353c",
			"dc48558e-1671-42df-86f2-f1e298e81947"
		],
		"net_consideration": "1020.00008152",
		"gross_consideration": "1020.00000000",
		"executed_fee_cost": "0.00008152",
		"remaining_qty": 0.0015368385851072673,
		"has_open_child_orders": false,
		"net_executed_quantity": "1020.0000815200",
		"exchange_summary_stats": [],
		"fee_summary_stats": [],
		"parent_order_name": null,
		"task_id": "",
		"parent_order_slug": null,
		"pct_executed": 99.81185387094618,
		"sequence": 8,
		"trader_name": "sam",
		"stop_triggered": null,
		"last_update": "2022-06-13T16:38:02.832911Z",
		"avg_fee_price": 1251.08,
		"leg_summary_stats": [
			{
				"product": "ETH-USD.IPERP",
				"side": "sell",
				"leg_id": null,
				"exchange_filled_qty": "1020",
				"net_executed_quantity": "1020.00008152",
				"exchange_filled_fee": "0",
				"gross_consideration": "1020",
				"net_consideration": "1020.00008152",
				"avg_price": "1251.08",
				"avg_net_price": "1251.07990001"
			}
		],
		"leg_fee_summary_stats": [
			{
				"product": "ETH-USD.IPERP",
				"side": "sell",
				"exchange_fee_currency__slug": "ETH",
				"exchange_filled_qty": "1020",
				"net_executed_quantity": "1020.00008152",
				"exchange_filled_fee": "0.00008152",
				"gross_consideration": "1020",
				"net_consideration": "1020.00008152",
				"avg_price": "1251.08",
				"avg_net_price": "1251.07990001"
			}
		],
		"leg_exchange_summary_stats": [
			{
				"product": "ETH-USD.IPERP",
				"side": "sell",
				"exchange_fee_currency__slug": "ETH",
				"exchange__slug": "binancefutures",
				"exchange_filled_qty": "1020",
				"net_executed_quantity": "1020.00008152",
				"exchange_filled_fee": "0.00008152",
				"gross_consideration": "1020",
				"net_consideration": "1020.00008152",
				"avp_price": "1251.08",
				"avg_net_price": "1251.07990001"
			}
		],
		"order_events": [],
		"cade_algo_params": {
			"side": "sell",
			"quantity": "0.8168324232",
			"sweep_pct": "0.150000",
			"aggression": "neutral",
			"max_order_count": "2",
			"post_exchange_count": "1",
			"market_data_exchanges": "binancefutures,binance"
		},
		"tca_report": null,
		"order_type": "smart post",
		"order_status": "closed",
		"aggression": "neutral",
		"move_percentage": null,
		"percent_limit_price": null,
		"max_post_size": null,
		"is_twap": false,
		"end_time": null,
		"interval_length": null,
		"tolerance_pct": null,
		"sweep_pct": "0.1500",
		"reprice_seconds": null,
		"reprice_randomization_pct": null,
		"currency_pair": "ETH-USD.IPERP",
		"quantity": "0.8168324232",
		"limit_price": null,
		"side": "sell",
		"exchanges": [
			"binance",
			"binancefutures"
		],
		"market_data_exchanges": [
			"binance",
			"binancefutures"
		],
		"client_order_id": "68f76d35-2665-461c-8d15-6bdfb4676fae",
		"strategy": "yttest3",
		"post_exchange_count": 1,
		"max_order_count": 2,
		"margin_leverage": null,
		"margin_borrow_pct": null,
		"parent_order": null,
		"ignore_filter_levels": null,
		"use_funding_currency": false,
		"stop_price": null,
		"end_offset": null,
		"quantity_currency": "ETH",
		"sweep_to": null,
		"notes": null,
		"volume_constraint_pct": null,
		"volume_constraint_interval_secs": null,
		"volume_constraint_max_post_size": null,
		"min_post_size": null
	}`)
)
