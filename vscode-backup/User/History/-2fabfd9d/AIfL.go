// Package exchangeDAL will handle all interfacing with trading with exchanges
package exchangeDAL

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

type DAL interface {
	// GetBalanceForStrategy will get the balance for a given strategy from an exchange
	GetBalanceForStrategy(
		ctx context.Context,
		wl wlog.Logger,
		chain coinroutesapi.CurrencyType,
		strategy *entities.Strategy,
	) (float64, error)

	// GetRawBalanceForStrategy will get and return the raw balance response from the exchange
	GetRawBalanceForStrategy(
		ctx context.Context,
		wl wlog.Logger,
		chain coinroutesapi.CurrencyType,
		strategy *entities.Strategy,
	) (*[]entities.StrategyBalance, error)

	// GetPositionQuantityForStrategy will get the quantity of the current
	// position for a strategy from a given exchange
	GetPositionQuantityForStrategy(
		ctx context.Context,
		wl wlog.Logger,
		strategy *entities.Strategy,
		currencyPair coinroutesapi.CurrencyPairType,
	) (float64, error)

	// GetPositionForStrategy will get the the current position for a strategy
	// from a given exchange
	GetPositionForStrategy(
		ctx context.Context,
		wl wlog.Logger,
		strategy *entities.Strategy,
		currencyPair coinroutesapi.CurrencyPairType,

	) (*entities.ContractPosition, error)
}

type exchangeService struct {
	cc *coinroutesapi.Client
}

func New(c *coinroutesapi.Client) (DAL, error) {
	return &exchangeService{
		cc: c,
	}, nil
}

func (es *exchangeService) GetBalanceForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	chain coinroutesapi.CurrencyType,
	strategy *entities.Strategy,
) (float64, error) {
	var amt float64
	balances, err := es.cc.GetBalances(ctx, strategy.Name)
	wl.Debugf("balances received: %+v", balances)
	if err != nil {
		return amt, err
	}

	// parse balances to find chain and set quantity
	for _, v := range *balances {
		// handle USDTM Margin (Binance)
		if strategy.Margin == entities.USDTM {
			if strategy.Exchange == entities.DYDX {
				chain = coinroutesapi.USD // DYDX uses USD for balance currency
			} else if strategy.Exchange == entities.Vertex {
				chain = coinroutesapi.USDC
			} else {
				chain = coinroutesapi.USDT
			}
		}

		// handle USDM Margin (FTX)
		if strategy.Margin == entities.USDM {
			chain = coinroutesapi.USD
		}

		// handle coin Margin
		if strategy.Margin == entities.CoinM || strategy.Margin == entities.CoinD {
			switch strategy.CurrencyPair {
			case entities.ETHInversePerpetual:
				chain = coinroutesapi.ETH
			case entities.BTCInversePerpetual:
				chain = coinroutesapi.BTC
			default:
				return amt, fmt.Errorf("invalid currency pair for coin margin")
			}
		}

		// handle spot Margin
		if strategy.Margin == entities.Spot {
			chain = coinroutesapi.USDT
		}

		if strings.ToLower(v.Currency) == string(chain) &&
			v.Exchange == string(strategy.Exchange) {
			amt, err = strconv.ParseFloat(v.Amount, 32)
			if err != nil {
				return amt, err
			}
			return amt, nil
		}
	}

	return amt, fmt.Errorf("could not find balance")
}

func (es *exchangeService) GetRawBalanceForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	chain coinroutesapi.CurrencyType,
	strategy *entities.Strategy,
) (*[]entities.StrategyBalance, error) {
	balances, err := es.cc.GetBalances(ctx, strategy.Name)
	wl.Debugf("balances received: %+v", balances)
	if err != nil {
		return nil, err
	}

	balancesToReturn := []entities.StrategyBalance{}
	// convert response to StrategyBalance
	for _, v := range *balances {
		b := entities.StrategyBalance{
			Exchange: v.Exchange,
			Currency: v.Currency,
			Amount:   v.Amount,
			Account:  v.Account,
			USDValue: v.USDValue,
			USDPrice: v.USDPrice,
		}
		balancesToReturn = append(balancesToReturn, b)
	}

	return &balancesToReturn, nil

}

func (es *exchangeService) GetPositionQuantityForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	strategy *entities.Strategy,
	currencyPair coinroutesapi.CurrencyPairType,
) (float64, error) {
	positions, err := es.cc.GetPositions(
		ctx,
		strategy.Name,
	)
	var qty float64
	if err != nil {
		return qty, err
	}

	wl.Debugf("found positions: %+v", positions)
	if len(*positions) == 0 {
		wl.Infof("no positions found")
		return qty, nil // FIXME: make custom error type here for no-op
	}

	// parse positions to find same currency pair and set quantity
	for _, v := range *positions {
		wl.Debugf("cur: %s, search: %s", v.CurrencyPair, currencyPair)
		if v.CurrencyPair == string(currencyPair) {
			qty, err := strconv.ParseFloat(v.Quantity, 32)
			if err != nil {
				wl.Error(err)
				continue
			}

			return math.Abs(qty), nil
		}
	}

	return qty, fmt.Errorf("unable to find positions")
}

func (es *exchangeService) GetPositionForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	strategy *entities.Strategy,
	currencyPair coinroutesapi.CurrencyPairType,
) (*entities.ContractPosition, error) {
	// If SPOT strategy, use balance API instead of positions API
	if strategy.Margin == entities.Spot {
		wl.Debug("using balance API for SPOT strategy")
		var amt float64
		balances, err := es.cc.GetBalances(ctx, strategy.Name)
		wl.Debugf("balances received: %+v", balances)
		if err != nil {
			wl.Debugf("error found when getting balances, err: %s", err)
			return nil, err
		}

		// determine what we're looking for
		var currency string
		switch strategy.CurrencyPair {
		case entities.USDTBTC:
			currency = "BTC"
		case entities.USDTETH:
			currency = "ETH"
		default:
			return nil, fmt.Errorf("invalid currency pair for spot margin")
		}

		var cp *entities.ContractPosition

		for _, b := range *balances {
			if b.Currency == currency {
				// For SPOT strategy, we return the value of the coin in USD. 0 if not present
				amt, err = strconv.ParseFloat(b.USDValue, 32)
				if err != nil {
					wl.Infof("error found when parsing balance to float, err: %s", err)
					return nil, err
				}
				if amt < 100 {
					wl.Info("balance is less than $100, no position")
					return nil, nil // no position
				}
				cp = &entities.ContractPosition{
					Exchange:         b.Exchange,
					CurrencyPair:     strategy.CurrencyPair.String(),
					Side:             "long", // long is only possibility with SPOT here
					Quantity:         amt,
					QuantityCurrency: "USD", // always using USD for SPOT
					EntryPrice:       0,     // irrelevant for SPOT
					UnrealizedPnl:    0,     // irrelevant for SPOT
				}
				wl.Debugf("found position for SPOT: %+v", cp)
			}
		}

		return cp, nil
	}

	positions, err := es.cc.GetPositions(
		ctx,
		strategy.Name,
	)
	if err != nil {
		return nil, err
	}

	wl.Debugf("found positions: %+v", positions)
	if len(*positions) == 0 {
		wl.Infof("no positions found")
		return nil, nil // no positions found
	}

	// parse positions to find same currency pair
	for _, v := range *positions {
		wl.Debugf("cur: %s, search: %s", v.CurrencyPair, currencyPair)
		if v.CurrencyPair == string(currencyPair) {
			qty, err := strconv.ParseFloat(v.Quantity, 64)
			if err != nil {
				wl.Error(err)
				continue
			}

			ep, err := strconv.ParseFloat(v.EntryPrice, 64)
			if err != nil {
				wl.Error(err)
				continue
			}

			upnl, err := strconv.ParseFloat(v.UnrealizedPnl, 64)
			if err != nil {
				wl.Error(err)
				continue
			}

			cp := entities.ContractPosition{
				Exchange:         v.Exchange,
				CurrencyPair:     v.CurrencyPair,
				Side:             v.Side,
				Quantity:         math.Abs(qty),
				QuantityCurrency: v.QuantityCurrency,
				EntryPrice:       ep,
				UnrealizedPnl:    upnl,
			}

			return &cp, nil

		}
	}

	return nil, fmt.Errorf("unable to find positions")
}
