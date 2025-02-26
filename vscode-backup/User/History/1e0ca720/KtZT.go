package signalsvc

import (
	"context"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

func (ss *signalService) CreateStrategy(
	ctx context.Context,
	wl wlog.Logger,
	strategy *entities.Strategy,
) (*entities.Strategy, error) {

}
