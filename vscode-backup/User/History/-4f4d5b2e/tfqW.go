package signalsvc

import (
	"context"
	"fmt"
	"strconv"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

func (ss *signalService) tradePosition(
	ctx context.Context,
	wl wlog.Logger,
	chain entities.ChainType,
	strategy *entities.Strategy,
	signal *entities.Signal,
	desiredSide coinroutesapi.SideType,
	position *entities.ContractPosition,
	currencyPair coinroutesapi.CurrencyPairType,
) error {
	crStrategy := coinroutesapi.SupportedStrategy(strategy.Name)

	var orderAmount float64

	// handle fixed/compound strategies for position amount
	if position == nil {
		// neutral state, amount should simply be the available balance
		balance, err := ss.exdal.GetBalanceForStrategy(
			ctx,
			wl,
			coinroutesapi.CurrencyType(chain),
			strategy)
		if err != nil {
			return err
		}

		orderAmount = balance

		// if fixed amount and balance < fixed amount, default to balance
		if strategy.TradeStrategy == entities.Fixed {
			if balance > strategy.FixedTradeAmount {
				wl.Debugf("balance greater than trade amount %f > %f",
					balance,
					strategy.FixedTradeAmount,
				)
				orderAmount = strategy.FixedTradeAmount
			}
		}

	} else {
		// position found, amount is either fixed or compound
		switch strategy.TradeStrategy {
		case entities.Fixed:
			if strategy.FixedTradeAmount == 0 {
				return fmt.Errorf(
					"strategy with fixed trade missing trade amount: %f",
					strategy.FixedTradeAmount,
				)
			}
			wl.Debugf("setting order amount to fixed amount: %f", strategy.FixedTradeAmount)
			orderAmount = strategy.FixedTradeAmount

		case entities.Compound:
			if position.Quantity <= 0 {
				return fmt.Errorf("invalid position.quantity: %f", position.Quantity)
			}

			if position.EntryPrice <= 0 {
				return fmt.Errorf("invalid position.entry_price: %f", position.EntryPrice)
			}

			orderAmount = (position.Quantity * position.EntryPrice) + position.UnrealizedPnl
			wl.Debugf("setting order amount to compound amount: (%f * %f) + %f = %f",
				position.Quantity,
				position.EntryPrice,
				position.UnrealizedPnl,
				orderAmount,
			)

		default:
			return fmt.Errorf("missing trade strategy, cannot set amount for position")
		}

	}

	if orderAmount <= 0 {
		wl.Error(fmt.Errorf(("no-op: order amount not properly set: %f", orderAmount)))
		return ErrNoOpSignal
	}

	// apply buffer only to fixed trades
	// apply buffer to trade amount and overwrite payload
	if strategy.TradeStrategy == entities.Fixed {
		wl.Debugf(
			"applying buffer of %f. Original: %f Actual: %f",
			tradeAmountBuffer,
			orderAmount,
			orderAmount*tradeAmountBuffer,
		)
		orderAmount *= tradeAmountBuffer
	}

	// set leverage for strategy
	if strategy.Leverage == entities.TwoX {
		orderAmount *= 2
	}

	// send 1 order at fixed price in USDT
	orderPayload := coinroutesapi.ClientOrderCreateRequest{
		OrderType:          coinroutesapi.SmartPost,
		OrderStatus:        coinroutesapi.Open,
		Aggression:         coinroutesapi.Neutral,
		CurrencyPair:       currencyPair,
		Quantity:           strconv.FormatFloat(orderAmount, 'f', 10, 64),
		Side:               desiredSide,
		Strategy:           crStrategy,
		UseFundingCurrency: true, // needs to be in USDT
		EndOffset:          tradeTTL,
		IntervalLength:     intLength,
		IsTwap:             false,
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
	err = ss.insertNewOrder(ctx, wl, resp)
	if err != nil {
		wl.Error(err)
	}

	// log trade to BQ
	err = ss.dl.Log(ctx, wl, string(crStrategy), signal, resp)
	if err != nil {
		wl.Error(err)
	}

	return nil
}
