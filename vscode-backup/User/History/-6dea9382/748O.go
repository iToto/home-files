// Package signalsvc is the service that handles getting and processing trade signals
package signalsvc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/exchangeDAL"
	"yield-mvp/internal/signallogger"
	"yield-mvp/internal/tradelogger"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/signalapi"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid"
)

const (
	twentyMinutes     int     = 1200
	fiveMinutes       int     = 300
	fiveSeconds       int     = 5
	tradeAmountBuffer float64 = 0.98
	usdtTradeAmount   float64 = 500.00
	ftxTradeAmount    float64 = 200.00
	tradeTTL                  = fiveMinutes
	intLength                 = 1
)

type SVC interface {
	// GetAndProcessSignal will ping for a specified signal, parse it and process it
	GetAndProcessSignal(ctx context.Context,
		wl wlog.Logger,
		strats []entities.Strategy,
		chain entities.ChainType,
	) error
}

type signalService struct {
	db    *sqlx.DB
	sc    *signalapi.Client
	cc    *coinroutesapi.Client
	dl    *tradelogger.DataLogger
	sl    *signallogger.DataLogger
	exdal exchangeDAL.DAL
}

func New(
	s *signalapi.Client,
	c *coinroutesapi.Client,
	d *sqlx.DB,
	t *tradelogger.DataLogger,
	sl *signallogger.DataLogger,
	e exchangeDAL.DAL,
) (SVC, error) {
	return &signalService{
		db:    d,
		sc:    s,
		cc:    c,
		dl:    t,
		sl:    sl,
		exdal: e,
	}, nil
}

func (ss *signalService) checkForSignal(
	ctx context.Context,
	wl wlog.Logger,
	sig *signalapi.SignalResponse,
	strategy *entities.Strategy,
) (*entities.Signal, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %s", err)
	}

	signal := &entities.Signal{
		ID:       id.String(),
		Chain:    sig.Chain,
		Signal:   entities.SignalType(sig.Signal),
		Strategy: strategy.Name,
	}

	delta := false

	// get last chain trade from DB
	lastSignal, err := ss.queryLastSignalForStrategy(ctx, wl, strategy.Name)
	if err != nil {
		// check for DB error
		if errors.Is(err, ErrDBConnection) {
			// return this error so that it can be handled up the stack
			return nil, err
		}
		// if no history found, use the last signal
		wl.Infof(
			"no history found for strategy %s, using last signal %s on %s",
			strategy.Name,
			sig.LastTrade,
			sig.LastTradeTime,
		)
		signal.Signal = entities.SignalType(sig.LastTrade)
		signal.TradeTime = sig.LastTradeTime

		// log signal to BQ
		err = ss.sl.Log(ctx, wl, delta, sig, signal.ID, strategy.Name) // FIXME: Should be passing internal entities
		if err != nil {
			wl.Error(fmt.Errorf("unable to log signal to bigquery: %w", err))
		}
		return signal, nil
	}

	wl.Debugf("lastSignal: %+v", lastSignal)

	// handle no history found
	if lastSignal == nil {
		// this shouldn't happen as sqlx returns an error for 404
		// log signal to BQ
		err = ss.sl.Log(ctx, wl, delta, sig, signal.ID, strategy.Name) // FIXME: Should be passing internal entities
		if err != nil {
			wl.Error(fmt.Errorf("unable to log signal to bigquery: %w", err))
		}
		return nil, fmt.Errorf("no trade history found for strategy %s", strategy.Name)
	}

	// if last trade in DB == last_trade_time in signal, then we no-op
	if lastSignal.Signal == entities.SignalType(sig.LastTrade) &&
		lastSignal.TradeTime.Equal(sig.LastTradeTime) {
		wl.Debugf(
			"no-op: %s == %s && %s == %s",
			lastSignal.Signal,
			entities.SignalType(sig.LastTrade),
			lastSignal.TradeTime,
			sig.LastTradeTime,
		)
		signal.Signal = entities.Null
		signal.TradeTime = lastSignal.TradeTime
	} else { // new signal found!
		wl.Debugf(
			"op: %s == %s && %s == %s",
			lastSignal.Signal,
			entities.SignalType(sig.LastTrade),
			lastSignal.TradeTime,
			sig.LastTradeTime,
		)

		delta = true
		// we need to create a new signal with the last_trade details and return it
		signal.Signal = entities.SignalType(sig.LastTrade)
		signal.TradeTime = sig.LastTradeTime

	}

	// log signal to BQ
	wl.Debug("logging signal to BQ")
	err = ss.sl.Log(ctx, wl, delta, sig, signal.ID, strategy.Name) // FIXME: Should be passing internal entities
	if err != nil {
		wl.Error(fmt.Errorf("unable to log signal to bigquery: %w", err))
	}

	return signal, nil
}

func (ss *signalService) queryLastSignalForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	strategy string,
) (*entities.Signal, error) {
	var lastSignal entities.Signal

	query := "SELECTd id, chain, signal, strategy, trade_time, created_at, updated_at, deleted_at FROM mvp_signal_log WHERE strategy = $1 ORDER BY created_at DESC LIMIT 1"
	err := ss.db.Get(&lastSignal, query, strategy)
	if err != nil {
		err = fmt.Errorf("%w, error querying signal: %w", ErrDBConnection, err)
		return nil, err
	}

	return &lastSignal, nil
}

func (ss *signalService) insertLatestSignalTradedForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	chain entities.ChainType,
	signal *entities.Signal,
	strategy *entities.Strategy,
) error {
	signal.CreatedAt = null.NewTime(time.Now(), true)
	query := `INSERT INTO mvp_signal_log (chain, strategy, signal, trade_time, created_at) 
	VALUES (:chain, :strategy, :signal, :trade_time, :created_at)`
	_, err := ss.db.NamedQuery(query, signal)
	if err != nil {
		return err
	}

	return nil
}

func (ss *signalService) insertNewOrder(
	ctx context.Context,
	wl wlog.Logger,
	o *coinroutesapi.ClientOrderCreateResponse,
) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return fmt.Errorf("error generating id: %s", err)
	}

	order := entities.Order{
		ID:            id.String(),
		ClientOrderId: o.ClientOrderId,
		Strategy:      o.Strategy,
		Status:        entities.OrderStatusType(o.OrderStatus),
		CurrencyPair:  o.CurrencyPair,
		AvgPrice:      o.AvgPrice,
		ExecutedQty:   o.ExecutedQty,
		CreatedAt:     null.NewTime(time.Now(), true),
	}

	var query string

	// check if finished_at is set and only insert it if it is
	if o.FinishedAt != "" {
		finishedAt, err := time.Parse(time.RFC3339Nano, o.FinishedAt)
		if err != nil {
			return err
		}

		order.FinishedAt = null.NewTime(finishedAt, true)
		query = `INSERT INTO mvp_order (
			id,
			client_order_id,
			strategy,
			status,
			currency_pair,
			avg_price,
			executed_qty,
			finished_at,
			created_at)
		VALUES (
			:id,
			:client_order_id,
			:strategy,
			:status,
			:currency_pair,
			:avg_price,
			:executed_qty,
			:finished_at,
			:created_at)`
	} else { // finished_at not set, therefore don't insert it
		query = `INSERT INTO mvp_order (
		id,
		client_order_id,
		strategy,
		status,
		currency_pair,
		avg_price,
		executed_qty,
		created_at)
	VALUES (
		:id,
		:client_order_id,
		:strategy,
		:status,
		:currency_pair,
		:avg_price,
		:executed_qty,
		:created_at)`

	}

	_, err = ss.db.NamedQuery(query, order)
	if err != nil {
		return err
	}

	return nil
}

func (ss *signalService) parseSignalAndCreateTradePayload(
	ctx context.Context,
	wl wlog.Logger,
	sig *entities.Signal,
	cur coinroutesapi.CurrencyType,
	str coinroutesapi.SupportedStrategy,
) (*coinroutesapi.ClientOrderCreateRequest, error) {
	var currencyPair coinroutesapi.CurrencyPairType

	if cur == coinroutesapi.BTC {
		currencyPair = coinroutesapi.BTCInversePerpetual
	} else {
		currencyPair = coinroutesapi.ETHInversePerpetual
	}

	// Parse signal and act on it
	if sig.Signal == entities.Buy || sig.Signal == entities.Long {
		// Handle Buy/Long
		payload := coinroutesapi.ClientOrderCreateRequest{
			OrderType:          coinroutesapi.SmartPost,
			OrderStatus:        coinroutesapi.Open,
			Aggression:         coinroutesapi.Neutral,
			CurrencyPair:       currencyPair,
			Quantity:           "", // this will be hydrated later
			Side:               coinroutesapi.Buy,
			Strategy:           str,
			UseFundingCurrency: true,
			EndOffset:          tradeTTL,
			IntervalLength:     intLength,
			IsTwap:             false,
		}

		// set twap to true for BTC
		if cur == coinroutesapi.BTC {
			payload.IsTwap = true
		}

		return &payload, nil

	} else if sig.Signal == entities.Sell || sig.Signal == entities.Short {
		// Handle Sell/Short
		payload := coinroutesapi.ClientOrderCreateRequest{
			OrderType:          coinroutesapi.SmartPost,
			OrderStatus:        coinroutesapi.Open,
			Aggression:         coinroutesapi.Neutral,
			CurrencyPair:       currencyPair,
			Quantity:           "", // this will be hydrated later
			Side:               coinroutesapi.Sell,
			Strategy:           str,
			UseFundingCurrency: false,
			EndOffset:          tradeTTL,
			IntervalLength:     intLength,
			IsTwap:             false,
		}

		// set twap to true for BTC
		if cur == coinroutesapi.BTC {
			payload.IsTwap = true
		}

		return &payload, nil

	} else if sig.Signal == entities.Null {
		wl.Debug("no-op: null signal received")
		return nil, nil
	} else {
		return nil, fmt.Errorf("invalid signal received: %s", sig.Signal)
	}
}

func (ss *signalService) calculateOrderQuantity(
	ctx context.Context,
	wl wlog.Logger,
	chain entities.ChainType,
	strategy *entities.Strategy,
	payload *coinroutesapi.ClientOrderCreateRequest,
) (float64, error) {
	switch payload.Side {
	case coinroutesapi.Buy:
		// drop contracts (positions)
		position, err := ss.exdal.GetPositionQuantityForStrategy(
			ctx,
			wl,
			strategy,
			payload.CurrencyPair,
		)
		if err != nil {
			return position, err
		}

		return position, nil

	case coinroutesapi.Sell:
		// short available balance (create contracts)
		balance, err := ss.exdal.GetBalanceForStrategy(
			ctx,
			wl,
			coinroutesapi.CurrencyType(chain),
			strategy,
		)
		if err != nil {
			return balance, err
		}

		// add multiple for long strategy
		if strategy.Type == entities.LongShort {
			wl.Debug("multiplying amount for LongShort")
			balance = 2 * balance // 2x leverage for long-short
		}

		return balance, nil

	}

	return 0.0, fmt.Errorf("unsupported side: %s", payload.Side)
}
