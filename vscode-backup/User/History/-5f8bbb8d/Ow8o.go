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

func (ss *signalService) calculateCoinValue(
	ctx context.Context,
	wl wlog.Logger,
	strategy *entities.Strategy,
	position *entities.ContractPosition,
	marketPrice float64,
	balance float64,
) (float64, error) {
	if position == nil {
		// no position, return 0
		return 0, nil
	}

	// (position.Quantity / position.Entry Price / strategy.Account Leverage) + balance.Amount (Balance) + position.Unrealized P&L
	if strategy.AccountLeverage.IsZero() {
		return 0, fmt.Errorf("strategy.AccountLeverage is zero")
	}

	// (Quantity / Market Price / Account Leverage) + Balance

	return balance, nil
}

func (ss *signalService) calculateAccountDollarValue(
	ctx context.Context,
	wl wlog.Logger,
	position *entities.ContractPosition,
	balance *[]entities.StrategyBalance,
	coinValue float64,
	marketPrice float64,
) (float64, error) {
	// TODO: Coin Value X Market Price
	// Question: which property from CR are we getting the market price from?

	if coinValue == 0 {
		return 0, nil
	}

	return coinValue * marketPrice, nil // TODO: confirm this is correct
}

func (ss *signalService) calculateMultiplier(
	ctx context.Context,
	wl wlog.Logger,
	position *entities.ContractPosition,
) (int8, error) {
	var multiplier int8

	if position == nil {
		return 1, nil
	}

	// If “Side” is long or empty, value is 1, if “Side” is short, value is -1
	if coinroutesapi.SideType(position.Side).IsEquivalent(coinroutesapi.Buy) ||
		coinroutesapi.SideType(position.Side).IsEquivalent(coinroutesapi.Neut) {
		multiplier = 1
	} else if coinroutesapi.SideType(position.Side).IsEquivalent(coinroutesapi.Sell) {
		multiplier = -1
	} else {
		// assuming empty here
		multiplier = 1
	}

	return multiplier, nil
}

func (ss *signalService) calculateSignalMultiplier(
	ctx context.Context,
	wl wlog.Logger,
	desiredSide coinroutesapi.SideType,
) (int8, error) {
	var signalMultiplier int8

	// This is 1 if signal is long or neutral, -1 if signal is short
	// Note state is different from side, side can be short but we state can be long,
	// if the number of coins in our account is worth more than the short value of the contracts

	if desiredSide == coinroutesapi.Buy || desiredSide == coinroutesapi.Neut {
		signalMultiplier = 1
	} else if desiredSide == coinroutesapi.Sell {
		signalMultiplier = -1
	} else {
		return 0, fmt.Errorf("invalid desired side found:%s", desiredSide)
	}
	return signalMultiplier, nil
}

// - this is  a number we can use to determine our position
func (ss *signalService) calculatePositionPercent(
	ctx context.Context,
	wl wlog.Logger,
	position *entities.ContractPosition,
	multiplier int8,
	accountDollarValue float64,
) (float64, error) {
	var positionQuantity float64

	if position == nil {
		positionQuantity = 0
	} else {
		positionQuantity = position.Quantity
	}

	if positionQuantity == 0 && accountDollarValue == 0 {
		return 0, nil
	}

	// (position.Quantity X Multiplier + Account $ Value) / Account $ Value
	// debug values
	wl.Debugf("**** positionQuantity:%f", positionQuantity)
	wl.Debugf("**** multiplier:%d", multiplier)
	wl.Debugf("**** accountDollarValue:%f", accountDollarValue)
	return (positionQuantity*float64(multiplier) + accountDollarValue) / accountDollarValue, nil
}

func (ss *signalService) calculateDesiredExposure(
	ctx context.Context,
	wl wlog.Logger,
	strategy *entities.Strategy,
	signalMultiplier int8,
	desiredSide coinroutesapi.SideType,
	accountDollarValue float64,
) (float64, error) {
	// If signal is neutral then it is 0
	// If signal is long or short, then it is then:
	// strategy.fixed trade amount X strategy.Trade Leverage X Signal Multiplier

	// debug values
	// wl.Debugf("**** strategy.FixedTradeAmount.ValueOrZero():%f", strategy.FixedTradeAmount.ValueOrZero())
	// wl.Debugf("**** strategy.Leverage:%d", strategy.Leverage)
	// wl.Debugf("**** signalMultiplier:%d", signalMultiplier)
	// wl.Debugf("**** desiredSide:%s", desiredSide)

	if desiredSide == coinroutesapi.Neut {
		return 0, nil
	}

	if strategy.TradeStrategy == entities.Fixed {
		wl.Debug("fixed trade strategy")
		if strategy.FixedTradeAmount.IsZero() {
			return 0, fmt.Errorf("fixed trade amount is zero")
		}

		if desiredSide != coinroutesapi.Sell && desiredSide != coinroutesapi.Buy {
			return 0, fmt.Errorf("invalid desired side found:%s", desiredSide)
		}

		return strategy.FixedTradeAmount.ValueOrZero() * float64(strategy.Leverage) * float64(signalMultiplier), nil

	} else {
		wl.Debug("compound trade strategy")
		if accountDollarValue == 0 {
			return 0, fmt.Errorf("account dollar value is zero")
		}

		if desiredSide != coinroutesapi.Sell && desiredSide != coinroutesapi.Buy {
			return 0, fmt.Errorf("invalid desired side found:%s", desiredSide)
		}

		return accountDollarValue * float64(strategy.Leverage) * float64(signalMultiplier), nil
	}

}

func (ss *signalService) calculateCurrentContractExposure(
	ctx context.Context,
	wl wlog.Logger,
	position *entities.ContractPosition,
	multiplier int8,
) (float64, error) {
	// position.Quantity X Multiplier
	if position == nil {
		return 0, nil
	}

	return position.Quantity * float64(multiplier), nil
}

func (ss *signalService) calculateCurrentSide(
	ctx context.Context,
	wl wlog.Logger,
	position *entities.ContractPosition,
	positionPercent float64,
) (entities.SignalType, error) {
	var positionQuantity float64

	if position == nil {
		positionQuantity = 0
	} else {
		positionQuantity = position.Quantity
	}

	// If position.Quantity is 0, we are long (we hold no short positions and we have BTC in our account, so we are long)
	if positionQuantity == 0 || position == nil {
		return entities.Buy, nil
	}

	// If position.side is long (quantity can be anything), we are long.
	// This means we hold btc in the account and have long contracts on top of it, so we are actually long with leverage
	if position.Side == string(entities.Buy) {
		return entities.Buy, nil
	}

	// If “Position %” >10% we are long, <-10% we are short. Otherwise neutral
	if positionPercent > 0.15 {
		return entities.Buy, nil
	} else if positionPercent < -0.15 {
		return entities.Sell, nil
	} else {
		return entities.Neutral, nil
	}

}

func (ss *signalService) calculateTradeAmount(
	ctx context.Context,
	wl wlog.Logger,
	desiredExposure float64,
	currentContractExposure float64,
	accountDollarValue float64,
	marketPrice float64,
) (float64, error) {
	// Desire exposure - Current contract exposure - Account $ Value
	return (desiredExposure - currentContractExposure - accountDollarValue) / marketPrice, nil
}
