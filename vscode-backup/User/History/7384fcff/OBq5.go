package signalsvc

import (
	"context"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

func (ss *signalService) GetETHSignal(ctx context.Context, wl wlog.Logger, strats []string) (*entities.Signal, error) {

}
