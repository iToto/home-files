package signalsvc

import (
	"context"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

func (ss *signalService) GetStrategies(
	ctx context.Context,
	wl wlog.Logger,
	enabledFilterPresent bool,
	enabled bool,
) ([]*entities.Strategy, error) {
	wl.Debug("get strategies")
	var strategies []*entities.Strategy

	if enabledFilterPresent && enabled {
		wl.Debugf("get active strategies")
		strategies, err := ss.strategyDAL.GetActiveStrategies(ctx, wl)
		if err != nil {
			return nil, err
		}

		return strategies, nil
	}

	wl.Debugf("get all strategies")

	strategies, err := ss.strategyDAL.GetStrategies(ctx, wl)
	if err != nil {
		return nil, err
	}

	return strategies, nil
}
