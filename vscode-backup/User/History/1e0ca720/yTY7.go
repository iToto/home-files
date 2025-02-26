package signalsvc

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"

	"github.com/oklog/ulid"
)

func (ss *signalService) CreateStrategy(
	ctx context.Context,
	wl wlog.Logger,
	strategy *entities.Strategy,
) (*entities.Strategy, error) {
	wl.Debug("creating strategy")
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %s", err)
	}

	strategy.ID = id.String()

}
