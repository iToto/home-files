package signalsvc

import (
	"context"
	"fmt"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"

	"github.com/guregu/null"
)

func (ss *signalService) UpdateStrategy(
	ctx context.Context,
	wl wlog.Logger,
	strategy *entities.Strategy,
) (*entities.Strategy, error) {
	wl.Debug("updating strategy")

	// validate fields
	if strategy.ID == "" {
		return nil, fmt.Errorf("cannot update strategy without ID")
	}

	// validate fields
	if strategy.Exchange != "" {
		err := strategy.Exchange.Validate()
		if err != nil {
			wl.Infof("invalid exchange: %w", err)
			return nil, err
		}
	}

	// ensure name isn't already taken
	existingName, err := ss.strategyDAL.GetStrategyByName(ctx, wl, strategy.Name)
	if err == nil && existingName.Name == strategy.Name {
		wl.Info("cannot update strategy with existing name")
		return nil, fmt.Errorf("cannot update strategy with existing name")
	}

	existing, err := ss.strategyDAL.GetStrategyByID(ctx, wl, strategy.ID)
	if err != nil {
		return nil, err
	}

	isDifferent := false

	// check if properties differ, and update
	if strategy.Name != "" && existing.Name != strategy.Name {
		existing.Name = strategy.Name
		isDifferent = true
	}

	if strategy.Margin != "" && existing.Margin != strategy.Margin {
		existing.Margin = strategy.Margin
		isDifferent = true
	}

	if strategy.Leverage != 0 && existing.Leverage != strategy.Leverage {
		existing.Leverage = strategy.Leverage
		isDifferent = true
	}

	if strategy.AccountLeverage.ValueOrZero() != 0 && existing.AccountLeverage != strategy.AccountLeverage {
		existing.AccountLeverage = strategy.AccountLeverage
		isDifferent = true
	}

	if strategy.FixedTradeAmount.ValueOrZero() != 0 && existing.FixedTradeAmount != strategy.FixedTradeAmount {
		existing.FixedTradeAmount = strategy.FixedTradeAmount
		isDifferent = true
	}

	// handle FixedTradeAmount null value
	if strategy.FixedTradeAmount.IsZero() && !existing.FixedTradeAmount.IsZero() {
		existing.FixedTradeAmount = strategy.FixedTradeAmount
		isDifferent = true
	}

	if strategy.TradeStrategy != "" && existing.TradeStrategy != strategy.TradeStrategy {
		// verify fixed trade amount
		if strategy.TradeStrategy == entities.Fixed && strategy.FixedTradeAmount.IsZero() {
			return nil, fmt.Errorf("required fixed trade amount for fixed strategies")
		}
		existing.TradeStrategy = strategy.TradeStrategy
		isDifferent = true
	}

	if strategy.CurrencyPair != "" && existing.CurrencyPair != strategy.CurrencyPair {
		existing.CurrencyPair = strategy.CurrencyPair
		isDifferent = true
	}

	if strategy.Exchange != "" && existing.Exchange != strategy.Exchange {
		existing.Exchange = strategy.Exchange
		isDifferent = true
	}

	if strategy.Type != "" && existing.Type != strategy.Type {
		existing.Type = strategy.Type
		isDifferent = true
	}

	if strategy.SignalSourceID != "" && existing.SignalSourceID != strategy.SignalSourceID {
		// signal_source
		_, err = ss.signalDAL.GetSignalSourceByID(ctx, wl, strategy.SignalSourceID)
		if err != nil {
			wl.Infof("could not find signal: %w", err)
			return nil, err
		}

		existing.SignalSourceID = strategy.SignalSourceID
		isDifferent = true
	}

	if !isDifferent {
		// no change, no-op
		return existing, nil
	}

	existing.UpdatedAt = null.NewTime(time.Now(), true)

	err = ss.strategyDAL.UpdateStrategy(ctx, wl, existing)
	if err != nil {
		wl.Infof("could not update strategy with error: %w", err)
		return nil, err
	}

	return existing, nil

}
