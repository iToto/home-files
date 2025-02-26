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
	wl.Debug("get strategies")
	var strategies []*entities.Strategy

	strategies, err := ss.strategyDAL.GetStrategies(ctx, wl)
	if err != nil {
		return nil, err
	}

	return strategies, nil
}
