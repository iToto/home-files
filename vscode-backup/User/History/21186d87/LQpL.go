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

	// validate fields
	err = strategy.Exchange.Validate()
	if err != nil {
		wl.Infof("invalid exchange: %w", err)
		return nil, err
	}

	// ensure name isn't already taken
	existing, err := ss.strategyDAL.GetStrategyByName(ctx, wl, strategy.Name)
	if err == nil && existing.Name == strategy.Name {
		wl.Info("cannot create strategy with existing name")
		return nil, fmt.Errorf("cannot create strategy with existing name")
	}

	// signal_source
	_, err = ss.signalDAL.GetSignalSourceByID(ctx, wl, strategy.SignalSourceID)
	if err != nil {
		wl.Infof("could not find signal: %w", err)
		return nil, err
	}

	// user
	_, err = ss.userDAL.GetUserByID(ctx, wl, strategy.UserID)
	if err != nil {
		wl.Infof("could not find user: %w", err)
		return nil, err
	}

	// verify fixed trade amount
	if strategy.TradeStrategy == entities.Fixed && strategy.FixedTradeAmount.IsZero() {
		return nil, fmt.Errorf("required fixed trade amount for fixed strategies")
	}

	err = ss.strategyDAL.CreateStrategy(ctx, wl, strategy)
	if err != nil {
		wl.Infof("could not create strategy with error: %w", err)
		return nil, err
	}

	return strategy, nil

}
