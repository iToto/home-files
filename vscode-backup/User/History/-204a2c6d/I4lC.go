// Package signalsvc is the service that handles getting and processing trade signals
package signalsvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"yield/signal-logger/internal/entities"
	"yield/signal-logger/internal/signallogger"
	"yield/signal-logger/internal/signalloggerv2"
	"yield/signal-logger/internal/wlog"
	"yield/signal-logger/pkg/signalapi"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

type signalStrategyLog struct {
	ID        string              `db:"id" json:"id"`
	Chain     string              `db:"chain" json:"chain"`
	Signal    entities.SignalType `db:"signal" json:"signal"`
	Strategy  string              `db:"strategy" json:"strategy"`
	TradeTime time.Time           `db:"trade_time" json:"trade_time"`
	CreatedAt null.Time           `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt null.Time           `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt null.Time           `db:"deleted_at" json:"deleted_at,omitempty"`
}

type SVC interface {
	// GetAndProcessSignal will ping for a specified signal, parse it and process it
	GetAndProcessSignal(ctx context.Context,
		wl wlog.Logger,
		signals []entities.SignalSource,
	) error
}

type signalService struct {
	db   *sqlx.DB
	sc   *signalapi.Client
	sl   *signallogger.DataLogger
	slv2 *signalloggerv2.DataLogger
}

func New(
	db *sqlx.DB,
	s *signalapi.Client,
	sl *signallogger.DataLogger,
	slv2 *signalloggerv2.DataLogger,
) (SVC, error) {
	return &signalService{
		db:   db,
		sc:   s,
		sl:   sl,
		slv2: slv2,
	}, nil
}

func (ss *signalService) queryLastSignalForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	strategy string,
) (*signalStrategyLog, error) {
	var lastSignal signalStrategyLog

	query := "SELECT id, chain, signal, strategy, trade_time, created_at, updated_at, deleted_at FROM mvp_signal_log WHERE strategy = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1"
	err := ss.db.Get(&lastSignal, query, strategy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no rows found")
		}

		wl.Error(fmt.Errorf("error querying signal: %w", err))
		return nil, fmt.Errorf("error with DB connection")
	}

	return &lastSignal, nil
}

func (ss *signalService) insertLatestSignalTradedForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	signal signalStrategyLog,
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
