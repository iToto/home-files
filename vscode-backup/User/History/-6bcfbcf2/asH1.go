package signalsvc

import (
	"context"
	"math"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

const LongThreshold = 200

func (ss *signalService) calculateMarginSide(
	ctx context.Context,
	wl wlog.Logger,
	position *entities.ContractPosition,
	balance float64,
) (entities.SignalType, error) {

	// long
	if math.Abs(position.Quantity) < LongThreshold {
		return entities.Long, nil
	}

	// use calculation to decide if we are neutral or short
	posA := (balance - position.UnrealizedPnl) * (4 / 3) * position.EntryPrice

	// (0.4839345100 - 0.02138496) * (4/3) * 1721.620 = 1,061.7794083613

	// neutral
	neutralLowerBound := math.Abs(position.Quantity) * .90  // 1070.00000 * .90 = 963.00000
	neutralUpperBound := math.Abs(position.Quantity) * 1.10 // 1070.00000 *

	if posA <= neutralUpperBound && posA >= neutralLowerBound {
		return entities.Neutral, nil
	}

	// assuming if we are neither long or neutral, we are short
	return entities.Short, nil
}
