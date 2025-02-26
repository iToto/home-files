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

}
