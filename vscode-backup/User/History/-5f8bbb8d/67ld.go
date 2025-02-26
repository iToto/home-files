package signalsvc

import (
	"context"
	"fmt"
	"strconv"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

func (ss *signalService) tradeCoinMStrategyV2(
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
		false,
	)
	if err != nil {
		return err
	}

	b := balance
	q := position.UsdValue
	totalExposure := (b + q) / b

	var currentSide entities.SignalType

	// 0 we are neutral
	// >0 we are long
	// < 0 we are short
	// 0.12, we are 12% exposed
	// -0.09 we are 9% exposed short
	switch {
	case totalExposure == 0:
		currentSide = entities.Neutral
	case totalExposure > 0:
		currentSide = entities.Long
	case totalExposure < 0:
		currentSide = entities.Short
	}

	var targetExposure float64
	var existingExposure float64
	var tradeAmount float64

	if strategy.TradeStrategy == entities.Compound {
		// Order to send compounding:
		// Target exposure: (B) X L (set this to 0 if the signal is neutral)
		// Existing exposure: (B + Q)
		// Order to send: (B) X L -  (B + Q)
		if currentSide == entities.Neutral {
			targetExposure = 0
		} else {
			targetExposure = balance * float64(strategy.Leverage)
		}

		existingExposure = b + q
	} else {
		// Order to send when fixed
		// Target exposure: Fixed amount X leverage, if long, 0 if neutral
		// Existing exposure: (B + Q)
		// Order to send: Fixed amount X L -  (B + Q)
		if currentSide == entities.Neutral {
			targetExposure = 0
		} else {
			targetExposure = strategy.FixedTradeAmount.ValueOrZero() * float64(strategy.Leverage)
		}

		existingExposure = b + q
	}

	tradeAmount = targetExposure - existingExposure

	// no-op if tradeAmount is zero
	if tradeAmount == 0 {
		wl.Infof("no-op: trade amount is zero")
		// log signal and no-op
		err = ss.insertLatestSignalTradedForStrategy(ctx, wl, chain, signal, strategy)
		if err != nil {
			return err
		}
		return nil // no-op
	}
	// debug values
	wl.Debugf("coinm-debug: chainMarketPrice: %f", chainMarketPrice)
	wl.Debugf("coinm-debug: currentSide: %s", currentSide)
	wl.Debugf("coinm-debug: b: %f", b)
	wl.Debugf("coinm-debug: q: %d", q)
	wl.Debugf("coinm-debug: totalExposure: %f", totalExposure)
	wl.Debugf("coinm-debug: targetExposure: %f", targetExposure)
	wl.Debugf("coinm-debug: existingExposure: %f", existingExposure)
	wl.Debugf("coinm-debug: tradeAmount: %f", tradeAmount)

	// create order
	orderPayload := coinroutesapi.ClientOrderCreateRequest{
		OrderType:          coinroutesapi.SmartPost,
		OrderStatus:        coinroutesapi.Open,
		Aggression:         coinroutesapi.Neutral,
		CurrencyPair:       currencyPair,
		Quantity:           strconv.FormatFloat(tradeAmount, 'f', 10, 64),
		Side:               desiredSide,
		Strategy:           strategy.Name,
		UseFundingCurrency: true, // true to use USD
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
