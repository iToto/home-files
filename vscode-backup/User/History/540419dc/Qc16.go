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
	Chain       string     `json:"chain"`
	Signal      SignalType `json:"signal"`
	CurrentTime time.Time  `json:"current_time"`
}
