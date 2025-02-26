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
	wl.Debug("strategy is CoinD, calculating trade amount")
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
	// get current side
	var positionQuantity float64
	var currentSide entities.SignalType

	if position == nil {
		positionQuantity = 0
	} else {
		positionQuantity = position.Quantity
	}

	if positionQuantity == 0 {
		currentSide = entities.Neutral
	} else if position.Side == string(entities.Long) || position.Side == string(entities.Buy) {
		currentSide = entities.Buy
	} else if position.Side == string(entities.Short) || position.Side == string(entities.Sell) {
		currentSide = entities.Sell
	} else {
		wl.Debugf("coind-debug: invalid position side found:%s", position.Side)
		return fmt.Errorf("invalid position side found:%s", position.Side)
	}

	wl.Debugf("coind-debug: currentSide: %s", currentSide)
	wl.Debugf("coind-debug: coinValue: %f", coinValue)
	wl.Debugf("coind-debug: multiplier: %d", multiplier)
	wl.Debugf("coind-debug: accountDollarValue: %f", accountDollarValue)
	wl.Debugf("coind-debug: desiredSide: %s", desiredSide)

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

	// If we are neutral and get a long, send an order for the number of BTC that appear in “Balance”
	// 	If you can only send order in $ value, multiple Balance by Market Price and send that
	// If we are long and get a neutral, send short orders for the $ value of the Quantity

	var tradeAmount float64

	if currentSide.IsEquivalent(entities.Neutral) && desiredSide == coinroutesapi.Buy {
		// set tradeAmount for the number of BTC that appear in “Balance”
		tradeAmount = balance
	} else if currentSide.IsEquivalent(entities.Long) && desiredSide == coinroutesapi.Neut {
		// set tradeAmount for the $ value of the Quantity
		if position.QuantityCurrency == "USD" {
			tradeAmount = positionQuantity / chainMarketPrice
		} else {
			tradeAmount = positionQuantity
		}
	} else {
		// unhandled state transition
		return fmt.Errorf("unhandled state transition: %s -> %s", currentSide, desiredSide)
	}

	// absolute value
	tradeAmount = math.Abs(tradeAmount)
	dollarAmount := tradeAmount * chainMarketPrice

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
	wl.Debugf("coind-debug: chainMarketPrice: %f", chainMarketPrice)
	wl.Debugf("coind-debug: currentSide: %s", currentSide)
	wl.Debugf("coind-debug: balance: %f", balance)
	wl.Debugf("coind-debug: position.qty: %f", positionQuantity)
	wl.Debugf("coind-debug: chainMarketPrice: %f", chainMarketPrice)
	wl.Debugf("coind-debug: dollarAmount: %f", dollarAmount)
	wl.Debugf("coind-debug: tradeAmount: %f", tradeAmount)

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
