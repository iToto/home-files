package signalsvc

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

func (ss *signalService) tradeCoinDStrategy(
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

	// get raw balances from exchange
	balances, err := ss.exdal.GetRawBalanceForStrategy(
		ctx,
		wl,
		coinroutesapi.CurrencyType(currencyPair),
		strategy,
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

	coinValue, err := ss.calculateCoinValue(ctx, wl, strategy, position, chainMarketPrice, balance)
	if err != nil {
		return fmt.Errorf("unable to calculate coin value: %w", err)
	}
	accountDollarValue, err := ss.calculateAccountDollarValue(ctx, wl, position, balances, coinValue, chainMarketPrice)
	if err != nil {
		return fmt.Errorf("unable to calculate account dollar value: %w", err)
	}
	multiplier, err := ss.calculateMultiplier(ctx, wl, position)
	if err != nil {
		return fmt.Errorf("unable to calculate multiplier: %w", err)
	}
	positionPercent, err := ss.calculatePositionPercent(ctx, wl, position, multiplier, accountDollarValue)
	if err != nil {
		return fmt.Errorf("unable to calculate position percent: %w", err)
	}
	// get current side
	currentSide, err := ss.calculateCurrentSide(ctx, wl, position, positionPercent)
	if err != nil {
		return fmt.Errorf("unable to calculate current side: %w", err)
	}

	wl.Debugf("coinm-debug: currentSide: %s", currentSide)
	wl.Debugf("coinm-debug: positionPercent: %f", positionPercent)
	wl.Debugf("coinm-debug: coinValue: %f", coinValue)
	wl.Debugf("coinm-debug: multiplier: %d", multiplier)
	wl.Debugf("coinm-debug: accountDollarValue: %f", accountDollarValue)
	wl.Debugf("coinm-debug: desiredSide: %s", desiredSide)

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

	currentContractExposure, err := ss.calculateCurrentContractExposure(ctx, wl, position, multiplier)
	if err != nil {
		return fmt.Errorf("unable to calculate current contract exposure: %w", err)
	}
	signalMultiplier, err := ss.calculateSignalMultiplier(ctx, wl, desiredSide)
	if err != nil {
		return fmt.Errorf("unable to calculate signal multiplier: %w", err)
	}
	desiredExposure, err := ss.calculateDesiredExposure(ctx, wl, strategy, signalMultiplier, desiredSide, accountDollarValue)
	if err != nil {
		return fmt.Errorf("unable to calculate desired exposure: %w", err)
	}

	tradeAmount, err := ss.calculateTradeAmount(ctx, wl, desiredExposure, currentContractExposure, accountDollarValue, chainMarketPrice)
	if err != nil {
		return fmt.Errorf("unable to calculate trade amount: %w", err)
	}

	if tradeAmount > 0 {
		desiredSide = coinroutesapi.Buy
	} else {
		desiredSide = coinroutesapi.Sell
	}

	// absolute value
	tradeAmount = math.Abs(tradeAmount)

	wl.Debugf("coinm-debug: tradeAmount before rounding: %f", tradeAmount)

	dollarAmount := tradeAmount * chainMarketPrice

	// debug dollar amount
	wl.Debugf("coinm-debug: dollarAmount before rounding: %f", dollarAmount)

	// For BTC round the final order we send to the nearest 100 And for eth round to the nearest 10
	// For Deribit, we need to round the order size to the nearest 10 for both BTC and ETH
	if chain == entities.BTC {
		if strategy.Exchange == entities.Deribit {
			wl.Debug("coinm-debug: rounding to nearest 10 BTC for Deribit")
			dollarAmount = math.Round(dollarAmount/10) * 10
		} else {
			dollarAmount = math.Round(dollarAmount/100) * 100
		}
	} else if chain == entities.ETH {
		dollarAmount = math.Round(dollarAmount/10) * 10
	}

	// debug dollar amount
	wl.Debugf("coinm-debug: dollarAmount after rounding: %f", dollarAmount)

	tradeAmount = dollarAmount / chainMarketPrice

	// debug trade amount
	wl.Debugf("coinm-debug: tradeAmount after rounding: %f", tradeAmount)

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
	wl.Debugf("coinm-debug: currentContractExposure: %f", currentContractExposure)
	wl.Debugf("coinm-debug: signalMultiplier: %d", signalMultiplier)
	wl.Debugf("coinm-debug: desiredExposure: %f", desiredExposure)
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
