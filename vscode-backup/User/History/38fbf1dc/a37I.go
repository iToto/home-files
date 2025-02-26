package strategyDAL

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/wlog"

	"github.com/jmoiron/sqlx"
)

type DAL interface {
	GetActiveStrategiesForSignalSourceAndUser(
		ctx context.Context,
		wl wlog.Logger,
		signalSourceID string,
		userID string,
	) ([]*entities.Strategy, error)

	GetActiveStrategies(
		ctx context.Context,
		wl wlog.Logger,
	) ([]*entities.Strategy, error)

	DisableStrategyByName(
		ctx context.Context,
		wl wlog.Logger,
		name string,
	) error

	EnableStrategyByName(
		ctx context.Context,
		wl wlog.Logger,
		name string,
	) error

	GetStrategyByName(
		ctx context.Context,
		wl wlog.Logger,
		name string,
	) (*entities.Strategy, error)

	CreateStrategy(
		ctx context.Context,
		wl wlog.Logger,
		strategy *entities.Strategy,
	) error
}

type strategyDAL struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (DAL, error) {
	return &strategyDAL{
		db: db,
	}, nil
}

func (sd *strategyDAL) GetActiveStrategiesForSignalSourceAndUser(
	ctx context.Context,
	wl wlog.Logger,
	signalSourceID string,
	userID string,
) ([]*entities.Strategy, error) {
	var rows []entities.Strategy

	query := "SELECT id, enabled, user_id, signal_source_id, type, name, exchange, margin, leverage, trade_strategy, fixed_trade_amount, currency_pair FROM mvp_strategy WHERE enabled = TRUE AND signal_source_id = $1 AND user_id = $2"
	err := sd.db.Select(&rows, query, signalSourceID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no active strategies found for user and signal")
		}

		return nil, err
	}

	var strategies []*entities.Strategy
	for _, row := range rows {
		strategy := row
		strategies = append(strategies, &strategy)
	}

	return strategies, nil
}

func (sd *strategyDAL) GetActiveStrategies(
	ctx context.Context,
	wl wlog.Logger,
) ([]*entities.Strategy, error) {
	var rows []entities.Strategy

	query := "SELECT id, enabled, user_id, signal_source_id, type, name, exchange, margin, leverage, trade_strategy, fixed_trade_amount, currency_pair FROM mvp_strategy WHERE enabled = TRUE"
	err := sd.db.Select(&rows, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no active strategies found")
		}

		return nil, err
	}

	var strategies []*entities.Strategy
	for _, row := range rows {
		strategy := row
		strategies = append(strategies, &strategy)
	}

	return strategies, nil
}

func (sd *strategyDAL) DisableStrategyByName(
	ctx context.Context,
	wl wlog.Logger,
	name string,
) error {
	query := "UPDATE mvp_strategy SET enabled = FALSE WHERE name = $1;"
	_, err := sd.db.Query(query, name)
	if err != nil {
		return err
	}
	return nil
}

func (sd *strategyDAL) EnableStrategyByName(
	ctx context.Context,
	wl wlog.Logger,
	name string,
) error {
	query := "UPDATE mvp_strategy SET enabled = TRUE WHERE name = $1;"
	_, err := sd.db.Query(query, name)
	if err != nil {
		return err
	}
	return nil
}

func (sd *strategyDAL) GetStrategyByName(
	ctx context.Context,
	wl wlog.Logger,
	name string,
) (*entities.Strategy, error) {
	var strategy entities.Strategy

	query := "SELECT id, enabled, user_id, signal_source_id, type, name, exchange, margin, leverage, trade_strategy, fixed_trade_amount, currency_pair FROM mvp_strategy WHERE name = $1"
	err := sd.db.Get(&strategy, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no signal found")
		}

		return nil, err
	}

	return &strategy, nil
}

func (sd *strategyDAL) CreateStrategy(
	ctx context.Context,
	wl wlog.Logger,
	strategy *entities.Strategy,
) error {
	query :=
		`INSERT INTO "public"."mvp_strategy"
("id", "enabled", "user_id", "signal_source_id", "type", "name", "exchange", "margin", "leverage", "trade_strategy", "fixed_trade_amount", "currency_pair")
VALUES
(:id, :enabled, :user_id, :signal_source_id, :type, :name, :exchange, :margin, :leverage, :trade_strategy, :fixed_trade_amount, :currency_pair)`

	_, err := sd.db.NamedExec(query, strategy)
	if err != nil {
		return err
	}

	return nil
}
