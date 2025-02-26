package entities

import (
	"time"
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

type ChainType string

const (
	BTC ChainType = "btc"
	ETH ChainType = "eth"
)

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
	TLS           bool
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

// V2
// "fetch_result_status": "SUCCESS",
// "fetch_type": "DIRECT",
// "strategy_state": "LONG",
// "strategy_version": "eth-r15",
// "last_checked": "2022-07-20 01:03:03",
// "last_trade_signal": "LONG",
// "last_trade_signal_time": "2022-07-16 16:15:04"
