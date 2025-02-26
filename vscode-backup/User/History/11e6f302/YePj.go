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

	oas.onAir = onAir

	return nil

}

func (oas *onAirService) GetOnAirStatus(
	ctx context.Context,
	wl wlog.Logger,
) (entities.OnAirStatus, error) {

	return entities.OnAirStatus{}, nil
}

func (oas *onAirService) ToggleOnAirStatus(
	ctx context.Context,
	wl wlog.Logger,
) (entities.OnAirStatus, error) {
	return entities.OnAirStatus{}, nil
}
