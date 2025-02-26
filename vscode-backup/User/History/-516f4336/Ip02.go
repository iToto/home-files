package signalsvc

import (
	"context"
	"fmt"
	"strconv"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
)

func (ss *signalService) tradeNeutral(
	ctx context.Context,
	wl wlog.Logger,
	chain entities.ChainType,
	strategy string,
	signal *entities.Signal,
	position *entities.ContractPosition,
	currencyPair coinroutesapi.CurrencyPairType,
) error {
	// to go neutral, side must be invert of position
	desiredSide := coinroutesapi.SideType(position.Side).GetInverseSide()
	if desiredSide == coinroutesapi.Na {
		return fmt.Errorf("could not get inverted side to go neutral: %s", position.Side)
	}

	// make order for inverse position
	neutralOrder := coinroutesapi.ClientOrderCreateRequest{
		OrderType:          coinroutesapi.SmartPost,
		OrderStatus:        coinroutesapi.Open,
		Aggression:         coinroutesapi.Neutral,
		CurrencyPair:       currencyPair,
		Quantity:           strconv.FormatFloat(position.Quantity, 'f', 10, 64),
		Side:               desiredSide,
		Strategy:           strategy,
		UseFundingCurrency: false, // needs to be in coin
		// EndOffset:          tradeTTL,
		// IntervalLength:     intLength,
		// IsTwap:             false,
	}
	// make neutral trade
	wl.Debugf("about to create neutral order with payload: %+v", neutralOrder)
	neutralResp, err := ss.cc.CreateClientOrders(ctx, &neutralOrder)
	if err != nil {
		return err
	}

	if neutralResp.ClientOrderId == "" {
		return fmt.Errorf("no client_order_id found in response: %+v", neutralResp)
	}

	wl.Infof("neutral order placed with coinroutes: %+v", neutralResp)

	// record neutral trade in order table
	err = ss.insertNewOrder(ctx, wl, neutralResp)
	if err != nil {
		wl.Error(err)
	}

	// log trade to BQ
	err = ss.dl.Log(ctx, wl, string(strategy), signal, neutralResp)
	if err != nil {
		wl.Error(err)
	}

	return nil
}
