package signalsvc

import (
	"context"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

func (ss *signalService) GetStrategies(
	ctx context.Context,
	wl wlog.Logger,
) ([]*entities.Strategy, error) {
	var strategies []*entities.Strategy
	return strategies, nil
}
