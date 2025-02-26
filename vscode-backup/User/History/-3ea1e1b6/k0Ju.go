package onair

import (
	"context"
	"on-air/internal/entities"
	"on-air/internal/wlog"
	"time"

	"github.com/guregu/null"
)

type SVC interface {
	SetOnAirStatus(ctx context.Context, wl wlog.Logger, onAir entities.OnAirStatus) (entities.OnAirStatus, error)
	GetOnAirStatus(ctx context.Context, wl wlog.Logger) (entities.OnAirStatus, error)
	ToggleOnAirStatus(ctx context.Context, wl wlog.Logger) (entities.OnAirStatus, error)
}

type onAirService struct {
	// add any dependencies here (DB, Client, etc.)
	onAir entities.OnAirStatus
}

func New(
// pass in any dependencies
) (SVC, error) {

	// initiatlize onAir
	initOnAir := entities.OnAirStatus{
		IsOnAir:     false,
		LastUpdated: null.TimeFrom(time.Now()),
		LastOnAir:   null.Time{},
	}

	return &onAirService{onAir: initOnAir}, nil
}
