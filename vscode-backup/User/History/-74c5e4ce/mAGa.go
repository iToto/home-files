package signalsvc

import (
	"context"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

func (ss *signalService) CreateSignalSource(
	ctx context.Context,
	wl wlog.Logger,
	signal *entities.SignalSource,
) (*entities.SignalSource, error) {

}
