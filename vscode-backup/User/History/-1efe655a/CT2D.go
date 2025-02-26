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
	signalStrategy *entities.SignalStrategies,
) error {
	// if we don't have any strategies, no-op signal
	if len(signalStrategy.Strategies) < 1 {
		wl1.Info("found no strategies, skipping")
		return ErrNoOpSignal
	}
	var signalResp *signalapi.SignalResponseV2
	var err error

	signalResp, err = ss.sc.GetSignalFromIPV2(ctx, wl1, *signalStrategy.Signal)
	if err != nil {
		return ErrSignalClient
	}

	if signalResp == nil {
		err = errors.New("received empty signal from api")
		return err
	}

	// check for empty string signal (no-op)
	if signalResp.LastTradeSignal == "" || signalResp.StrategyState == "" {
		wl1.Info("found empty signal, cannot operate")
		return ErrNoOpSignal
	}

	wl1.Infof("signal received: %+v", signalResp)

	// execute signal on all strategies
	for _, currentStrategy := range signalStrategy.Strategies {
		wl := wlog.WithStrategy(wl1, currentStrategy.Name)

		// parse last signal to see if we have an operation to make for this strategy
		signal, err := ss.checkForSignal(ctx, wl, signalResp, currentStrategy)
		if err != nil {
			// check for DB error
			if errors.Is(err, ErrDBConnection) {
				// return this DB error up the stack so that it can be properly handled
				// we want to stop processing on DB errors to prevent a cascade of errors
				return err
			}

			if errors.Is(err, ErrNoOpSignal) {
				// this signal cannot be processed, return up the stack so it can be handled
				return err
			}

			if errors.Is(err, ErrNoSignalHistory) {
				// simply log this signal and no-op (we will trade on the next signal)
				// record processed signal in log

				// normalize signal to buy/sell
				if signal.Signal == entities.Long {
					signal.Signal = entities.Buy
				}

				if signal.Signal == entities.Short {
					signal.Signal = entities.Sell
				}

				err = ss.insertLatestSignalTradedForStrategy(ctx, wl, signalStrategy.Chain, signal, currentStrategy)
				if err != nil {
					wl.Error(err)
				}
				continue

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

		// convert long to buy and short to sell
		if desiredSide == coinroutesapi.Long {
			desiredSide = coinroutesapi.Buy
			signal.Signal = entities.Buy
		}

		if desiredSide == coinroutesapi.Short {
			desiredSide = coinroutesapi.Sell
			signal.Signal = entities.Sell
		}

		if !desiredSide.IsValid() {
			wl.Infof("invalid side from signal: %s, skipping", desiredSide)
			continue
		}

		/********************************
		* check if we want to go neutral
		*********************************/
		if signal.Signal == entities.Neutral {
			wl.Debug("found neutral signal")
			cp := coinroutesapi.CurrencyPairType(currentStrategy.CurrencyPair)

			position, err := ss.exdal.GetPositionForStrategy(
				ctx,
				wl,
				currentStrategy,
				cp,
			)
			if err != nil {
				wl.Debugf("could not get position for strategy %w", err)
				return err
			}
			// check for neutral state (no contracts)
			if position == nil {
				// record processed signal in log
				err = ss.insertLatestSignalTradedForStrategy(ctx, wl, signalStrategy.Chain, signal, currentStrategy)
				if err != nil {
					wl.Error(err)
				}
				// no-op as we are already in desired state of neutral
				wl.Info("no-op as no positions found so already in neutral state")
				continue
			}
			// trade neutral
			err = ss.tradeNeutral(ctx, wl, signalStrategy.Chain, currentStrategy, signal, position, cp)
			if err != nil {
				wl.Error(err)
			}

			// record processed signal in log
			err = ss.insertLatestSignalTradedForStrategy(ctx, wl, signalStrategy.Chain, signal, currentStrategy)
			if err != nil {
				wl.Error(err)
			}
			continue
		}

		// FTX exchange runs similar to USDT, so if we are using FTX, handle accordingly
		// NB: Order matters here as a strategy can overlap on exchange and margin.
		// in this case, exchange takes precedence, so we need to check for that first
		if currentStrategy.Exchange == entities.FTX {
			err := ss.tradeUSDTStrategy(ctx, wl, signalStrategy.Chain, currentStrategy, desiredSide, signal)
			if err != nil {
				wl.Error(fmt.Errorf("error while processing FTX strategy: %w", err))
				continue
			}
			continue
		}

		// check for USDT strategy
		if currentStrategy.Margin == entities.USDTM {
			wl.Debug("processing USDM strategy")
			err := ss.tradeUSDTStrategy(ctx, wl, signalStrategy.Chain, currentStrategy, desiredSide, signal)
			if err != nil {
				wl.Error(fmt.Errorf("error while processing USDT strategy: %w", err))
				continue
			}
			continue
		}

		if currentStrategy.Margin == entities.USDM {
			wl.Debug("processing USDM strategy")
			err := ss.tradeMarginStrategy(ctx, wl, signalStrategy.Chain, currentStrategy, desiredSide, signal)
			if err != nil {
				wl.Error(fmt.Errorf("error while processing USDM strategy: %w", err))
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
			coinroutesapi.CurrencyType(signalStrategy.Chain),
			currentStrategy.Name,
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
		amount, err := ss.calculateOrderQuantity(ctx, wl, signalStrategy.Chain, currentStrategy, payload)
		if err != nil {
			if errors.Is(err, ErrNoOpCurrentPosition) {
				wl.Debugf("no-op on current position, skipping")
				//  record signal in log so we don't try again
				err = ss.insertLatestSignalTradedForStrategy(ctx, wl, signalStrategy.Chain, signal, currentStrategy)
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
		err = ss.insertLatestSignalTradedForStrategy(ctx, wl, signalStrategy.Chain, signal, currentStrategy)
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
		err = ss.insertNewOrder(ctx, wl, resp, currentStrategy, signal)
		if err != nil {
			wl.Error(fmt.Errorf("unable to insert order into table: %w", err))
			continue
		}

		// log trade to BQ
		err = ss.dl.Log(ctx, wl, string(currentStrategy.Name), signal, resp)
		if err != nil {
			wl.Infof("unable to log trade to bigquery: %w", err)
			continue
		}
	}
	return nil
}
