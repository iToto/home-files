package signalsvc

import (
	"context"
	"math/rand"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
)

func (ss *signalService) CreateSignalSource(
	ctx context.Context,
	wl wlog.Logger,
	signal *entities.SignalSource,
) (*entities.SignalSource, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano))) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
}
