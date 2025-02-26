package signalsvc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"

	"github.com/guregu/null"
	"github.com/oklog/ulid"
)

func (ss *signalService) DisableStrategyAndGoNeutral(
	ctx context.Context,
	wl wlog.Logger,
	strategyname string,
) (*entities.Strategy, error) {
	strategy, err := ss.strategyDAL.GetStrategyByName(ctx, wl, strategyname)
	if err != nil {
		wl.Error(fmt.Errorf("unable to find strategy with error: %w", err))
		return nil, err
	}

	if strategy.Enabled {
		// disable strategy to prevent future activity
		wl.Debug("disabling strategy")
		err = ss.strategyDAL.DisableStrategyByName(ctx, wl, strategyname)
		if err != nil {
			wl.Error(fmt.Errorf("unable to disable strategy with error: %w", err))
			return nil, err
		}
		strategy.Enabled = false // update our copy w/o having to re-query DB
	}

	// delete history for strategy so that when we re-enable it will be as though
	// the strategy was added for the first time
	err = ss.deleteSignalStrategyHistory(ctx, wl, strategy.Name)
	if err != nil {
		if errors.Is(err, ErrNoSignalHistory) {
			wl.Info("unable to delete signal-strategy history as no record was found")
		} else {
			wl.Error(err)
		}
	}

	wl.Debug("history for strategy deleted")

	if strategy.Margin == entities.CoinD || strategy.Margin == entities.CoinM {
		// in this case we don't yet want to trade neutral as we need to change calculation
		wl.Info("strategy successfully disabled")
		return strategy, nil
	}

	chain := strings.ToLower(string(strategy.CurrencyPair[0:3]))
	// create dud neutral signal
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %s", err)
	}

	signal := &entities.Signal{
		ID:        id.String(),
		Chain:     chain,
		Signal:    entities.Neutral,
		Strategy:  strategy.Name,
		TradeTime: time.Now(),
		CreatedAt: null.NewTime(time.Now(), true),
	}

	cp := coinroutesapi.CurrencyPairType(strategy.CurrencyPair)

	position, err := ss.exdal.GetPositionForStrategy(
		ctx,
		wl,
		strategy,
		cp,
	)
	if err != nil {
		wl.Debugf("could not get position for strategy %w", err)
		return nil, err
	}
	// check for neutral state (no contracts)
	if position == nil {
		// no-op as we are already in desired state of neutral
		wl.Info("no-op as no positions found so already in neutral state")
		wl.Info("strategy successfully disabled")
		return strategy, nil
	}

	// trade neutral
	err = ss.tradeNeutral(
		ctx,
		wl,
		entities.ChainType(chain),
		strategy,
		signal,
		position,
		cp,
	)
	if err != nil {
		wl.Error(fmt.Errorf("error when attempting to go neutral: %w", err))
		return nil, err
	}

	wl.Info("strategy successfully disabled")
	return strategy, nil
}
