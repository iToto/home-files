package onair

import (
	"context"
	"on-air/internal/entities"
	"on-air/internal/wlog"
)

func (oas *onAirService) SetOnAirStatus(
	ctx context.Context,
	wl wlog.Logger,
	onAir entities.OnAirStatus,
) error {

	return nil

}

func GetOnAirStatus(
	ctx context.Context,
	wl wlog.Logger,
) (entities.OnAirStatus, error) {

	return entities.OnAirStatus{}, nil
}

func ToggleOnAirStatus(ctx context.Context, wl wlog.Logger) error {
	return nil
}
