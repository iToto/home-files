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
	Chain       string     `json:"chain"`
	Signal      SignalType `json:"signal"`
	CurrentTime time.Time  `json:"current_time"`
	CreatedAt   null.Time  `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   null.Time  `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt   null.Time  `db:"deleted_at" json:"deleted_at,omitempty"`
}
