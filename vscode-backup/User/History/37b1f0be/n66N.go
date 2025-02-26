package signalsvc

import (
	"context"
	"fmt"
	"strconv"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

func (ss *signalService) tradeCoinMStrategy(
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

	// get market price
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

	// get current position from exchange
	wl.Debug("strategy is CoinM, calculating trade amount")
	position, err := ss.exdal.GetPositionForStrategy(
		ctx,
		wl,
		strategy,
		currencyPair,
	)
	if err != nil {
		return err
	}

	// get balance from exchange
	balance, err := ss.exdal.GetBalanceForStrategy(
		ctx,
		wl,
		coinroutesapi.CurrencyType(currencyPair),
		strategy,
	)
	if err != nil {
		return err
	}

	// calculate if we are long/short/neutral
	currentSide, err := ss.calculateMarginSide(ctx, wl, position, balance)
	if err != nil {
		return err
	}

	// if we are in desired state, no-op
	if currentSide.IsEquivalent(entities.SignalType(desiredSide)) {
		wl.Infof("no-op: already in desired state: %s == %s", currentSide, desiredSide)
		// log signal and no-op
		err = ss.insertLatestSignalTradedForStrategy(ctx, wl, chain, signal, strategy)
		if err != nil {
			return err
		}
		return nil // no-op
	}

	wl.Debugf("current side: %s, desired side: %s", currentSide, desiredSide)

	// calculate trade amount
	amountToTrade := 0.0

	switch currentSide {
	case entities.Long:
		switch desiredSide {
		case coinroutesapi.Neut:
			amountToTrade = balance
			desiredSide = coinroutesapi.Sell
		case coinroutesapi.Short:
		case coinroutesapi.Sell:
			amountToTrade = 2 * balance
		default:
			return fmt.Errorf("unexpected desired side found for long state:%s", desiredSide)
		}

	case entities.Neutral:
		switch desiredSide {
		case coinroutesapi.Long:
		case coinroutesapi.Buy:
			amountToTrade = position.Quantity / chainMarketPrice
		case coinroutesapi.Short:
		case coinroutesapi.Sell:
			amountToTrade = balance + (balance-position.UnrealizedPnl)/3
		default:
			return fmt.Errorf("unexpected desired side found for neutral state:%s", desiredSide)
		}

	case entities.Short:
		switch desiredSide {
		case coinroutesapi.Long:
		case coinroutesapi.Buy:
			amountToTrade = position.Quantity / chainMarketPrice
		case coinroutesapi.Neut:
			desiredSide = coinroutesapi.Buy
			amountToTrade = (position.Quantity / chainMarketPrice) - (balance + balance - position.UnrealizedPnl)
		default:
			return fmt.Errorf("unexpected desired side found for short state:%s", desiredSide)
		}

	default:
		return fmt.Errorf("invalid current side found:%s", currentSide)
	}

	// create order
	orderPayload := coinroutesapi.ClientOrderCreateRequest{
		OrderType:          coinroutesapi.SmartPost,
		OrderStatus:        coinroutesapi.Open,
		Aggression:         coinroutesapi.Neutral,
		CurrencyPair:       currencyPair,
		Quantity:           strconv.FormatFloat(amountToTrade, 'f', 10, 64),
		Side:               desiredSide,
		Strategy:           strategy.Name,
		UseFundingCurrency: false, // false for coinM
		// EndOffset:          tradeTTL,
		// IntervalLength:     intLength,
		// IsTwap:             false,
	}

	// make trade
	wl.Debugf("about to create order with payload: %+v", orderPayload)
	resp, err := ss.cc.CreateClientOrders(ctx, &orderPayload)
	if err != nil {
		return err
	}

	if resp.ClientOrderId == "" {
		return fmt.Errorf("no client_order_id found in response: %+v", resp)
	}

	wl.Infof("order placed with coinroutes: %+v", resp)

	//  record signal in log
	err = ss.insertLatestSignalTradedForStrategy(ctx, wl, chain, signal, strategy)
	if err != nil {
		wl.Error(err)
	}

	// record trade in order table
	err = ss.insertNewOrder(ctx, wl, resp, strategy, signal)
	if err != nil {
		wl.Error(err)
	}

	// log trade to BQ
	err = ss.dl.Log(ctx, wl, strategy.Name, signal, resp)
	if err != nil {
		wl.Error(err)
	}

	return nil
}
