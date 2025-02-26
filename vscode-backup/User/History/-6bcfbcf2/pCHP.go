package signalsvc

import (
	"context"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

const LongThreshold = 200

func (ss *signalService) calculateMarginSide(
	ctx context.Context,
	wl wlog.Logger,
	position *entities.ContractPosition,
	balance float64,
) (string, error) {

	// long
	if balance < LongThreshold {
		return "long", nil
	}

	// use calculation to decide if we are neutral or short
	posA := (balance - position.UnrealizedPnl) * (4 / 3) * position.EntryPrice

	// neutral
	neutralLowerBound := balance * .90
	neutralUpperBound := balance * 1.10

	if posA <= neutralUpperBound && posA >= neutralLowerBound {
		return "neutral", nil
	}

	return "short", nil
}
