package signalsvc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

func (ss *signalService) tradeSpotStrategy(
	ctx context.Context,
	wl wlog.Logger,
	chain entities.ChainType,
	strategy *entities.Strategy,
	desiredSide coinroutesapi.SideType,
	signal *entities.Signal,
) error {
	currencyPair := coinroutesapi.CurrencyPairType(strategy.CurrencyPair)
	var chainMarketPrice float64
	var err error

	switch chain {
	case entities.BTC:
		chainMarketPrice, err = ss.btcPrice.GetPrice()
		if err != nil {
			return fmt.Errorf("unable to get market price: %w", err)
		}
	case entities.ETH:
		chainMarketPrice, err = ss.ethPrice.GetPrice()
		if err != nil {
			return fmt.Errorf("unable to get market price: %w", err)
		}
	default:
		return fmt.Errorf("invalid chain found:%s", chain)
	}

	// get current position (neutral, long, short)
	wl.Debug("strategy is Spot, calculating trade amount")
	position, err := ss.exdal.GetPositionForStrategy(
		ctx,
		wl,
		strategy,
		currencyPair,
	)
	if err != nil {
		return err
	}

	// check for neutral state (no contracts)
	if position == nil {
		// neutral state: just need to make position order
		wl.Debug("found neutral state")
		err := ss.tradePosition(
			ctx,
			wl,
			chain,
			strategy,
			signal,
			desiredSide,
			position,
			currencyPair,
			chainMarketPrice,
		)
		if err != nil {
			if errors.Is(err, ErrNoOpSignal) {
				// log signal and no-op
				err = ss.insertLatestSignalTradedForStrategy(ctx, wl, chain, signal, strategy)
				if err != nil {
					return err
				}
				return nil // no-op
			}
			return fmt.Errorf("could not trade into position from neutral: %w", err)
		}
		return nil

	} // end of starting neutral

	// starting in position: make two orders (go-neutral and then fixed price)

	currentSide := coinroutesapi.SideType(strings.ToLower(position.Side))

	if currentSide.IsEquivalent(desiredSide) {
		wl.Infof("position is already in desired state: %s = %s",
			position.Side,
			desiredSide,
		)
		// log signal and no-op
		err = ss.insertLatestSignalTradedForStrategy(ctx, wl, chain, signal, strategy)
		if err != nil {
			return err
		}
		return nil // no-op
	}

	// if desiredSide is long, we trade position (go long)
	switch desiredSide {
	case coinroutesapi.Long, coinroutesapi.Buy:
		wl.Debug("desired side is long")
		// if desiredSide is long, we trade position (go long)
		// long means buying BTC or eth depending which coin the strategy is for
		err := ss.tradePosition(
			ctx,
			wl,
			chain,
			strategy,
			signal,
			desiredSide,
			position,
			currencyPair,
			chainMarketPrice,
		)
		if err != nil {
			if errors.Is(err, ErrNoOpSignal) {
				// log signal and no-op
				err = ss.insertLatestSignalTradedForStrategy(ctx, wl, chain, signal, strategy)
				if err != nil {
					return err
				}
				return nil // no-op
			}
			return fmt.Errorf("could not trade into position from neutral: %w", err)
		}
		return nil
	case coinroutesapi.Short, coinroutesapi.Neut:
		// if desiredSide is short, we trade neutral
		// if desiredSide is neutral, we trade neutral
		wl.Debug("desired side is short")
		// neutral means buying usdt.
		// currencyPair =
		err := ss.tradeNeutral(ctx, wl, chain, strategy, signal, position, currencyPair)
		if err != nil {
			return fmt.Errorf("could not go neutral: %w", err)
		}
	}

	return nil
}
