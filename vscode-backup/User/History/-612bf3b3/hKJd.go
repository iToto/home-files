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

func (ss *signalService) tradeUSDTStrategy(
	ctx context.Context,
	wl wlog.Logger,
	chain entities.ChainType,
	strategy *entities.Strategy,
	desiredSide coinroutesapi.SideType,
	signal *entities.Signal,
) error {
	crStrategy := coinroutesapi.SupportedStrategy(strategy.Name)
	var currencyPair coinroutesapi.CurrencyPairType
	var chainMarketPrice float64

	switch chain {
	case entities.BTC:
		if strategy.Exchange == entities.FTX {
			currencyPair = coinroutesapi.USDBTCPerpetual
		} else {
			currencyPair = coinroutesapi.USDTBTCPerpetual
		}

		chainMarketPrice, err := ss.btcPrice.GetPrice()
		if err != nil {
			return fmt.Errorf("unable to get market price: %w", err)
		}
	case entities.ETH:
		currencyPair = coinroutesapi.USDTETHPerpetual
		chainMarketPrice, err := ss.ethPrice.GetPrice()
		if err != nil {
			return fmt.Errorf("unable to get market price: %w", err)
		}
	default:
		return fmt.Errorf("invalid chain found: %s", chain)
	}

	// get current position (neutral, long, short)
	wl.Debug("strategy is USDT, calculating trade amount")
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

	wl.Debug("changing position, creating two orders")

	// go neutral
	err = ss.tradeNeutral(ctx, wl, chain, crStrategy, signal, position, currencyPair)
	if err != nil {
		return fmt.Errorf("could not go neutral: %w", err)
	}

	// make position trade
	err = ss.tradePosition(
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
		return fmt.Errorf("could not trade into position: %w", err)
	}

	return nil
}
