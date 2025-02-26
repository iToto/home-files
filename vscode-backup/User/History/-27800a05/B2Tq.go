package signalsvc

import (
	"context"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

func (ss *signalService) GetSignals(
	ctx context.Context,
	wl wlog.Logger,
) ([]*entities.SignalSource, error) {
	wl.Debug("get signals")
	var signals []*entities.SignalSource

	signals, err := ss.signalDAL.GetSignalSources(ctx, wl)
	if err != nil {
		return nil, err
	}

	return signals, nil
}
