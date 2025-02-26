// Package signalsvc is the service that handles getting and processing trade signals
package signalsvc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/signalapi"
)

// GetAndProcessSignal will ping for a specified signal, parse it and process it
func (ss *signalService) GetAndProcessSignal(
	ctx context.Context,
	wl wlog.Logger,
	strats []entities.Strategy,
	chain entities.ChainType,
) error {
	var signalResp *signalapi.SignalResponse
	var err error

	switch chain {
	case entities.BTC:
		signalResp, err = ss.sc.GetBTCSignal(ctx)
		if err != nil {
			return err
		}
	case entities.ETH:
		signalResp, err = ss.sc.GetETHSignal(ctx)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported chain: %s", chain)
	}

	if signalResp == nil {
		err = errors.New("received empty signal from api")
		return err
	}

	wl.Infof("signal received: %+v", signalResp)

	// parse last signal to see if we have an operation to make
	signal, err := ss.checkForSignal(ctx, wl, signalResp)
	if err != nil {
		return err
	}

	if signal.Signal == entities.Null {
		// FIXME: This is just returning as a no-op. Should make this more relevant
		return nil
	}

	wl.Debugf("calculated signal: %+v", signal)

	// execute signal on all strategies
	for _, v := range strats {
		strategy := coinroutesapi.StrategyType(v.Name)
		wl = wlog.WithStrategy(wl, string(strategy))
		if err := strategy.Validate(); err != nil {
			wl.Debug("skipping invalid strategy")
			continue
		}

		var resp *coinroutesapi.ClientOrderCreateResponse

		wl.Debug("strategy process start")
		// parse signal from signal API and create proper payload for CoinRoutes
		payload, err := ss.parseSignalAndCreateTradePayload(
			ctx,
			wl,
			signal,
			coinroutesapi.CurrencyType(chain),
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
		amount, err := ss.calculateOrderQuantity(ctx, wl, chain, strategy, payload)
		if err != nil {
			wl.Error(err)
			continue
		}

		payload.Quantity = strconv.FormatFloat(amount, 'f', 10, 32)

		// check for zero amount
		orderFloat := 0.0
		if payload.Quantity != "" {
			orderFloat, err = strconv.ParseFloat(payload.Quantity, 32)
			if err != nil {
				return err
			}
		}

		if orderFloat == 0 {
			// YIE-45 this is an error as we cannot attain desired position
			wl.Error(
				fmt.Errorf(
					"cannot attain position with order quantity having zero value. Side: %s Quantity: %+v",
					payload.Side,
					payload.Quantity,
				))
			return err
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
			return err
		}

		if resp.ClientOrderId == "" {
			wl.Error(fmt.Errorf("no client_order_id found in response: %+v", resp))
			continue
		}

		wl.Infof("order placed with coinroutes: %+v", resp)

		// record trade in order table
		err = ss.insertNewOrder(ctx, wl, resp)
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
	return nil
}
