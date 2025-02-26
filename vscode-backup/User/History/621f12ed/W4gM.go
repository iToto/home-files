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
	wl1 wlog.Logger,
	strats []entities.Strategy,
	chain entities.ChainType,
) error {
	var signalResp *signalapi.SignalResponse
	var err error

	switch chain {
	case entities.BTC:
		signalResp, err = ss.sc.GetBTCSignal(ctx)
		if err != nil {
			return ErrSignalClient
		}
	case entities.ETH:
		signalResp, err = ss.sc.GetETHSignal(ctx)
		if err != nil {
			return ErrSignalClient
		}
	default:
		return fmt.Errorf("unsupported chain: %s", chain)
	}

	if signalResp == nil {
		err = errors.New("received empty signal from api")
		return err
	}

	wl1.Infof("signal received: %+v", signalResp)

	// execute signal on all strategies
	for _, v := range strats {
		coinRoutesStrategyName := coinroutesapi.SupportedStrategy(v.Name)
		wl := wlog.WithStrategy(wl1, string(coinRoutesStrategyName))
		if err := coinRoutesStrategyName.Validate(); err != nil {
			wl.Debug("skipping invalid strategy")
			continue
		}

		// parse last signal to see if we have an operation to make for this strategy
		signal, err := ss.checkForSignal(ctx, wl, signalResp, &v)
		if err != nil {
			// check for DB error
			if errors.Is(err, ErrDBConnection) {
				// return this DB error up the stack so that it can be properly handled
				return err
			}
			wl.Error(fmt.Errorf("error when checking signal for strategy: %w ", err))
			continue
		}

		if signal.Signal == entities.Null {
			// FIXME: This is just returning as a no-op. Should make this more relevant
			continue
		}

		wl.Debugf("calculated signal: %+v", signal)

		var resp *coinroutesapi.ClientOrderCreateResponse

		wl.Debug("strategy process start")

		desiredSide := coinroutesapi.SideType(signal.Signal)
		if !desiredSide.IsValid() {
			wl.Error(fmt.Errorf("invalid side from signal: %s", desiredSide))
			continue
		}

		// check for USDT strategy
		if coinRoutesStrategyName == coinroutesapi.ETHUSDT {
			err := ss.tradeUSDTStrategy(ctx, wl, chain, &v, desiredSide, signal)
			if err != nil {
				wl.Error(fmt.Errorf("error while processing USDT strategy: %w", err))
				continue
			}
			continue
		}

		// FTX exchange runs similar to USDT, so if we are using FTX, handle accordingly
		if v.Exchange == entities.FTX {
			err := ss.tradeUSDTStrategy(ctx, wl, chain, &v, desiredSide, signal)
			if err != nil {
				wl.Error(fmt.Errorf("error while processing FTX strategy: %w", err))
				continue
			}
			continue
		}

		/***********************
		* non USDT strategies
		***********************/

		// parse signal from signal API and create proper payload for CoinRoutes
		payload, err := ss.parseSignalAndCreateTradePayload(
			ctx,
			wl,
			signal,
			coinroutesapi.CurrencyType(chain),
			coinRoutesStrategyName,
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
		amount, err := ss.calculateOrderQuantity(ctx, wl, chain, &v, payload)
		if err != nil {
			if errors.Is(err, ErrNoOpCurrentPosition) {
				//  record signal in log so we don't try again
				err = ss.insertLatestSignalTradedForStrategy(ctx, wl, chain, signal, &v)
				if err != nil {
					wl.Error(err)
				}
				continue
			}
			wl.Error(err)
			continue
		}

		payload.Quantity = strconv.FormatFloat(amount, 'f', 10, 32)

		// check for zero amount
		orderFloat := 0.0
		if payload.Quantity != "" {
			orderFloat, err = strconv.ParseFloat(payload.Quantity, 32)
			if err != nil {
				wl.Error(fmt.Errorf("unable to parse order to float: %w", err))
				continue
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
			continue
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

		//  record signal in log
		err = ss.insertLatestSignalTradedForStrategy(ctx, wl, chain, signal, &v)
		if err != nil {
			wl.Error(err)
			continue
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
		err = ss.dl.Log(ctx, wl, string(coinRoutesStrategyName), signal, resp)
		if err != nil {
			wl.Error(fmt.Errorf("unable to log trade to bigquery: %w", err))
			continue
		}
	}
	return nil
}
