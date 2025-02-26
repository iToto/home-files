package onair

import (
	"context"
	"on-air/internal/entities"
	"on-air/internal/wlog"
)

func SetOnAirStatus(
	ctx context.Context,
	wl wlog.Logger,
	onAir entities.OnAirStatus,
) error {

}

func GetOnAirStatus(ctx context.Context, wl wlog.Logger) (entities.OnAirStatus, error)

func ToggleOnAirStatus(ctx context.Context, wl wlog.Logger) error
