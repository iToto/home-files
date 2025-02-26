package signalsvc

import (
	"context"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

func (ss *signalService) calculateMarginSide(
	ctx context.Context,
	wl wlog.Logger,
	position *entities.ContractPosition,
	balance float64,
) (string, error) {
	return "", nil
}
