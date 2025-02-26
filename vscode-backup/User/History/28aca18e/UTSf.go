package handler

import (
	"context"
	"net/http"
	"yield-mvp/internal/balancelogger"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/exchangeDAL"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/coinroutespriceconsumer"
)

func GetBalance(
	wl wlog.Logger,
	cc *coinroutesapi.Client,
	btcConsumer *coinroutespriceconsumer.Consumer,
	ethConsumer *coinroutespriceconsumer.Consumer,
	chains []entities.Chain,
	bl *balancelogger.DataLogger,
	exDAL exchangeDAL.DAL,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "balance")
		wl2 := wl
		for _, chain := range chains {
			wl2 = wlog.WithChain(wl, string(chain.Type))
			for _, strat := range chain.Strategies {
				wl2 = wlog.WithStrategy(wl2, strat.Name)
				// log balances post trade
				var currencyPair coinroutesapi.CurrencyPairType

				if chain.Type == entities.BTC {
					currencyPair = coinroutesapi.BTCInversePerpetual
				}

				if chain.Type == entities.ETH {
					currencyPair = coinroutesapi.ETHInversePerpetual
				}

				err := logAccountBalance(
					ctx,
					wl2,
					string(chain.Type),
					coinroutesapi.StrategyType(strat.Name),
					currencyPair,
					bl,
					cc,
					btcConsumer,
					ethConsumer,
					exDAL)
				if err != nil {
					wl.Infof("unable to log balances to bigquery: %s", err)
					continue
				}
			}
		}
	}
}

func logAccountBalance(
	ctx context.Context,
	wl wlog.Logger,
	chain string,
	str coinroutesapi.StrategyType,
	currencyPair coinroutesapi.CurrencyPairType,
	bl *balancelogger.DataLogger,
	cr *coinroutesapi.Client,
	btcConsumer *coinroutespriceconsumer.Consumer,
	ethConsumer *coinroutespriceconsumer.Consumer,
	exDAL exchangeDAL.DAL,
) error {
	wl.Debug("logging account balance")

	balance, err := exDAL.GetBalanceForStrategy(ctx, wl, str, coinroutesapi.CurrencyType(chain))
	if err != nil {
		return err
	}
	contracts, err := exDAL.GetPositionForStrategy(ctx, wl, str, currencyPair)
	if err != nil {
		return err
	}

	switch chain {
	case string(entities.BTC):
		price, err := ethConsumer.GetPrice()
		if err != nil {
			return err
		}

	}

	bl.Log(
		ctx,
		wl,
		chain,
		string(str),
		float32(balance),
		float32(price),
		float32(contracts),
	)
	return nil
}
