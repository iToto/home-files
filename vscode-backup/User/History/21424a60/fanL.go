package entities

import "github.com/guregu/null"

type OnAirStatus struct {
	IsOnAir     bool
	LastUpdated null.Time
	LastOnAir   null.Time
}
