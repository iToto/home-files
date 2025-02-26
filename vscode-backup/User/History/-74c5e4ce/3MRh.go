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

const (
	acceptedSignalVersion = 2
)

func (ss *signalService) CreateSignal(
	ctx context.Context,
	wl wlog.Logger,
	signal *entities.SignalSource,
) (*entities.SignalSource, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %s", err)
	}

	signal.ID = id.String()

	// validate fields

	// only allow signal version 2
	if signal.SignalVersion != acceptedSignalVersion {
		return nil, ErrClientBadVersion
	}

	// ensure type is valid
	if signal.Type != entities.ETH && signal.Type != entities.BTC {
		return nil, ErrClientBadSignalType
	}

	err = ss.signalDAL.CreateSignalSource(ctx, wl, signal)
	if err != nil {
		wl.Infof("could not create signal with error: %w", err)
		return nil, err
	}

	return signal, nil
}
