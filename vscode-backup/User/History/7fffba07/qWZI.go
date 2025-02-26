package handler

import (
	"context"
	"net/http"
	"strings"
	"yield-mvp/internal/balancelogger"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/exchangeDAL"
	"yield-mvp/internal/strategyDAL"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/coinroutespriceconsumer"
	"yield-mvp/pkg/render"
)

func GetBalance(
	wl wlog.Logger,
	cc *coinroutesapi.Client,
	btcConsumer *coinroutespriceconsumer.Consumer,
	ethConsumer *coinroutespriceconsumer.Consumer,
	strategyDal strategyDAL.DAL,
	bl *balancelogger.DataLogger,
	exDAL exchangeDAL.DAL,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "balance")

		// get all strategies to calculate balance for
		strategies, err := strategyDal.GetActiveStrategies(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		for _, strat := range strategies {
			wl3 := wlog.WithStrategy(wl, strat.Name)

			chain := strings.ToLower(string(strat.CurrencyPair[0:3]))

			err := logAccountBalance(
				ctx,
				wl3,
				chain,
				strat,
				coinroutesapi.CurrencyPairType(strat.CurrencyPair),
				bl,
				cc,
				btcConsumer,
				ethConsumer,
				exDAL)
			if err != nil {
				wl3.Infof("unable to log balances to bigquery: %s", err)
				continue
			}
		}
	}
}

func logAccountBalance(
	ctx context.Context,
	wl wlog.Logger,
	chain string,
	strategy *entities.Strategy,
	currencyPair coinroutesapi.CurrencyPairType,
	bl *balancelogger.DataLogger,
	cr *coinroutesapi.Client,
	btcConsumer *coinroutespriceconsumer.Consumer,
	ethConsumer *coinroutespriceconsumer.Consumer,
	exDAL exchangeDAL.DAL,
) error {
	wl.Debug("logging account balance")

	balance, err := exDAL.GetBalanceForStrategy(
		ctx,
		wl,
		coinroutesapi.CurrencyType(chain),
		strategy,
		true,
	)
	if err != nil {
		return err
	}
	contracts, err := exDAL.GetPositionQuantityForStrategy(
		ctx,
		wl,
		strategy,
		currencyPair,
	)
	if err != nil {
		return err
	}

	var price float64

	switch chain {
	case string(entities.BTC):
		price, err = btcConsumer.GetPrice()
		if err != nil {
			return err
		}
	case string(entities.ETH):
		price, err = ethConsumer.GetPrice()
		if err != nil {
			return err
		}
	}

	bl.Log(
		ctx,
		wl,
		chain,
		strategy.Name,
		float32(balance),
		float32(price),
		float32(contracts),
	)
	return nil
}
