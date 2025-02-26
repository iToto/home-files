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
		str coinroutesapi.StrategyType,
		chain coinroutesapi.CurrencyType,
	) (float64, error)

	// GetPositionQuantityForStrategy will get the quantity of the current
	// position for a strategy from a given exchange
	GetPositionQuantityForStrategy(
		ctx context.Context,
		wl wlog.Logger,
		str coinroutesapi.StrategyType,
		currencyPair coinroutesapi.CurrencyPairType,
	) (float64, error)

	// GetPositionForStrategy will get the the current position for a strategy
	// from a given exchange
	GetPositionForStrategy(
		ctx context.Context,
		wl wlog.Logger,
		str coinroutesapi.StrategyType,
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
	str coinroutesapi.StrategyType,
	chain coinroutesapi.CurrencyType,
) (float64, error) {
	var amt float64
	balances, err := es.cc.GetBalances(ctx, str)
	wl.Debugf("balances received: %+v", balances)
	if err != nil {
		return amt, err
	}

	// parse balances to find chain and set quantity
	for _, v := range *balances {
		// handle USDT strategies
		if str == coinroutesapi.ETHUSDT {
			chain = coinroutesapi.USDT
		}
		if strings.ToLower(v.Currency) == string(chain) &&
			v.Exchange == "binancefutures" {
			amt, err = strconv.ParseFloat(v.Amount, 32)
			if err != nil {
				return amt, err
			}
			return amt, nil
		}
	}

	return amt, fmt.Errorf("could not find balance")
}

func (es *exchangeService) GetPositionQuantityForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	str coinroutesapi.StrategyType,
	currencyPair coinroutesapi.CurrencyPairType,
) (float64, error) {
	positions, err := es.cc.GetPositions(ctx, str)
	var qty float64
	if err != nil {
		return qty, err
	}

	wl.Debugf("found positions: %+v", positions)
	if len(*positions) == 0 {
		wl.Infof("no-op: no positions found")
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
	str coinroutesapi.StrategyType,
	currencyPair coinroutesapi.CurrencyPairType,
) (*entities.ContractPosition, error) {
	positions, err := es.cc.GetPositions(ctx, str)
	if err != nil {
		return nil, err
	}

	wl.Debugf("found positions: %+v", positions)
	if len(*positions) == 0 {
		wl.Infof("no-op: no positions found")
		return nil, nil // FIXME: make custom error type here for no-op
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

			cp := entities.ContractPosition{
				Exchange:         v.Exchange,
				CurrencyPair:     v.CurrencyPair,
				Side:             v.Side,
				Quantity:         math.Abs(qty),
				QuantityCurrency: v.QuantityCurrency,
				EntryPrice:       ep,
			}

		}
	}

	return nil, fmt.Errorf("unable to find positions")
}
