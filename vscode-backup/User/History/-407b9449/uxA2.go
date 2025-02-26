package signalsvc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

func (ss *signalService) GetAndProcessSignal(
	ctx context.Context,
	wl wlog.Logger,
	strats []string,
	chain entities.ChainType,
) (*entities.Signal, error) {

	signalResp, err := ss.sc.GetBTCSignal(ctx)
	if err != nil {
		return nil, err
	}

	if signalResp == nil {
		err = errors.New("received empty signal from api")
		return nil, err
	}

	wl.Infof("signal received: %+v", signalResp)

	// parse last signal to see if we have an operation to make
	signal, err := ss.checkForSignal(ctx, wl, signalResp)
	if err != nil {
		return nil, err
	}

	if signal.Signal == entities.Null {
		// FIXME: This is just returning as a no-op. Should make this more relevant
		return signal, nil
	}

	wl.Debugf("calculated signal: %+v", signal)

	// execute signal on all strategies
	for _, v := range strats {
		strategy := coinroutesapi.StrategyType(v)
		if err := strategy.Validate(); err != nil {
			wl.Info("skipping invalid strategy")
			continue
		}

		wl.Infof("handling strategy: %s", v)
		wl = wlog.WithStrategy(wl, v)
		var resp *coinroutesapi.ClientOrderCreateResponse

		// parse signal from signal API and create proper payload for CoinRoutes
		payload, err := ss.parseSignalAndCreateTradePayload(
			ctx,
			wl,
			signal,
			coinroutesapi.BTC,
			strategy,
		)
		if err != nil {
			wl.Error(fmt.Errorf("could not create payload for signal: %w", err))
			continue
		}

		// if we have a payload, then signal told us to make a trade
		if payload == nil {
			wl.Info("no payload, skipping")
			continue
		}

		// set amount
		switch payload.Side {
		case coinroutesapi.Buy:
			// drop contracts (positions)
			position, err := ss.getCoinRoutesPositionForStrategy(ctx, wl, strategy, cc, payload.CurrencyPair)
			if err != nil {
				wl.Error(err)
				continue
			}

			payload.Quantity = strconv.FormatFloat(position, 'f', 10, 32)

		case coinroutesapi.Sell:
			// short available balance (create contracts)
			balance, err := ss.getCoinRoutesBalanceForStrategy(
				ctx,
				wl,
				strategy,
				coinroutesapi.BTC)
			if err != nil {
				wl.Error(err)
				continue
			}

			// add multiple for long strategy
			if strategy == coinroutesapi.BTCLongShort {
				wl.Debug("multiplying amount for LongShort")

				balance = 2 * balance // 2x leverage for long-short
			}

			payload.Quantity = strconv.FormatFloat(balance, 'f', 10, 32)

		}

		// check for zero amount
		orderFloat := 0.0
		if payload.Quantity != "" {
			orderFloat, err = strconv.ParseFloat(payload.Quantity, 32)
			if err != nil {
				return nil, err
			}
		}

		if orderFloat == 0 {
			wl.Infof("no-op: quantity is zero value: %+v", payload.Quantity)
			return nil, err
		}

		// apply buffer to trade amount and overwrite payload
		wl.Debugf(
			"applying buffer of %f. Original: %f Actual: %f",
			tradeAmountBuffer,
			orderFloat,
			orderFloat*tradeAmountBuffer,
		)
		orderFloat *= tradeAmountBuffer
		payload.Quantity = strconv.FormatFloat(orderFloat, 'f', 10, 32)

		// Make Trade
		wl.Debugf("about to create order with payload: %+v", payload)
		resp, err = ss.cc.CreateClientOrders(ctx, payload)
		if err != nil {
			wl.Error(fmt.Errorf("unable to place order with error: %w", err))
			continue
		}

		//  record trade in log
		err = ss.insertLatestSignalTradedForChain(ctx, wl, chain, signal)
		if err != nil {
			return nil, err
		}

		if resp.ClientOrderId == "" {
			wl.Error(fmt.Errorf("no client_order_id found in response: %+v", resp))
			continue
		}

		wl.Infof("order placed with coinroutes: %+v", resp)

		// record trade in order table
		err = insertNewOrder(ctx, wl, resp, db)
		if err != nil {
			wl.Error(fmt.Errorf("unable to insert order into table: %w", err))
			continue
		}

		// log trade to BQ
		err = ss.dl.Log(ctx, wl, string(strategy), signal, resp)
		if err != nil {
			wl.Error(fmt.Errorf("unable to log trade to bigquery: %w", err))
			continue
		}
	}
}
