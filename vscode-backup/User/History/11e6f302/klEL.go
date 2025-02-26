package onair

import (
	"context"
	"on-air/internal/entities"
	"on-air/internal/wlog"
	"time"

	"github.com/guregu/null"
)

func (oas *onAirService) SetOnAirStatus(
	ctx context.Context,
	wl wlog.Logger,
	onAir entities.OnAirStatus,
) (entities.OnAirStatus, error) {

	oas.onAir = onAir
	oas.onAir.LastUpdated = null.TimeFrom(time.Now())

	if onAir.IsOnAir {
		oas.onAir.LastOnAir = null.TimeFrom(time.Now())
	}

	wl.Debugf("setting onAir: %v", oas.onAir)

	return oas.onAir, nil

}

func (oas *onAirService) GetOnAirStatus(
	ctx context.Context,
	wl wlog.Logger,
) (entities.OnAirStatus, error) {

	wl.Debugf("getting onAir: %v", oas.onAir)

	return oas.onAir, nil
}

func (oas *onAirService) ToggleOnAirStatus(
	ctx context.Context,
	wl wlog.Logger,
) (entities.OnAirStatus, error) {
	oas.onAir.LastUpdated = null.TimeFrom(time.Now())

	if oas.onAir.IsOnAir {
		wl.Debug("currenty on air, setting to off and setting last on air")
		oas.onAir.LastOnAir = null.TimeFrom(time.Now())
	}

	oas.onAir.IsOnAir = !oas.onAir.IsOnAir
	wl.Debugf("toggling onAir: %v", oas.onAir)
	return oas.onAir, nil
}
