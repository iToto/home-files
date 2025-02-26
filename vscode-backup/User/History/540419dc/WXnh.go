package handler

import "time"

type SignalType string

const (
	Short SignalType = "short"
	Sell  SignalType = "sell"
	Null  SignalType = "null"
	Long  SignalType = "long"
	Buy   SignalType = "buy"
)

type Signal struct {
	Chain         string     `json:"chain"`
	Signal        SignalType `json:"signal"`
	LastData      time.Time  `json:"last_data"`
	CurrentTime   time.Time  `json:"current_time"`
	LastTrade     SignalType `json:"last_trade"`
	LastTradeTime time.Time  `json:"last_trade_time"`
}
