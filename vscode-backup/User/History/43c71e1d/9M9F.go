package orderdal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/coinroutesapi"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid"
)

type orderReportRecord struct {
	Strategy    string `db:"strategy" json:"strategy,omitempty"`
	Coin        string `db:"coin" json:"coin,omitempty"`
	SignalID    string `db:"signal_id" json:"signal_id,omitempty"`
	CreatedAt   string `db:"created_at" json:"created_at,omitempty"`
	Direction   string `db:"direction" json:"direction,omitempty"`
	Signal      string `db:"signal" json:"signal,omitempty"`
	AvgPrice    string `db:"avg_price" json:"avg_price,omitempty"`
	ExecutedQty string `db:"executed_qty" json:"executed_qty,omitempty"`
}

type DAL interface {
	InsertNewOrder(
		ctx context.Context,
		wl wlog.Logger,
		o *coinroutesapi.ClientOrderCreateResponse,
		strat *entities.Strategy,
		sig *entities.Signal,
	) error

	GetOrdersForReport(
		ctx context.Context,
		wl wlog.Logger,
	) ([]*entities.OrderReportRecord, error)
}

type orderDAL struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (DAL, error) {
	return &orderDAL{
		db: db,
	}, nil
}

func (od *orderDAL) InsertNewOrder(
	ctx context.Context,
	wl wlog.Logger,
	o *coinroutesapi.ClientOrderCreateResponse,
	strat *entities.Strategy,
	sig *entities.Signal,
) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return fmt.Errorf("error generating id: %s", err)
	}

	order := entities.Order{
		ID:            id.String(),
		ClientOrderId: o.ClientOrderId,
		Strategy:      strat.Name,
		Status:        entities.OrderStatusType(o.OrderStatus),
		CurrencyPair:  o.CurrencyPair,
		AvgPrice:      o.AvgPrice,
		ExecutedQty:   o.ExecutedQty,
		CreatedAt:     null.NewTime(time.Now(), true),
		Side:          o.Side,
		Coin:          entities.ChainType(sig.Chain),
		SignalID:      strat.SignalSourceID,
		Signal:        sig.Signal,
	}

	if sig.Signal == entities.Sell {
		order.Signal = entities.Short
	}

	if sig.Signal == entities.Buy {
		order.Signal = entities.Long
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
			side,
			coin,
			signal_id,
			signal,
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
			:side,
			:coin,
			:signal_id,
			:signal,
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
		side,
		coin,
		signal_id,
		signal,
		created_at)
	VALUES (
		:id,
		:client_order_id,
		:strategy,
		:status,
		:currency_pair,
		:avg_price,
		:executed_qty,
		:side,
		:coin,
		:signal_id,
		:signal,
		:created_at)`

	}

	rows, err := od.db.NamedQuery(query, order)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (od *orderDAL) GetOrdersForReport(
	ctx context.Context,
	wl wlog.Logger,
) ([]*entities.OrderReportRecord, error) {
	var rows []orderReportRecord

	query := `select o.strategy, o.coin as coin, o.signal_id as signal_id, o.created_at as order_created_at, o.side as direction, o.signal as signal, o.avg_price, o.executed_qty
	from mvp_order o 
	join mvp_strategy st on st."name" = o.strategy
	where st.enabled = true
	order by strategy, o.created_at desc`
	err := od.db.Select(&rows, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no records found for report")
		}

		return nil, err
	}

	var reportData []*entities.OrderReportRecord
	for _, row := range rows {
		record := row
		reportData = append(reportData, &record)
	}

	return nil, nil
}
