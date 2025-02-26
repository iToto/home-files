// Package signalsvc is the service that handles getting and processing trade signals
package signalsvc

import (
	"context"
	"time"
	"yield/signal-logger/internal/entities"
	"yield/signal-logger/internal/signallogger"
	"yield/signal-logger/internal/signalloggerv2"
	"yield/signal-logger/internal/wlog"
	"yield/signal-logger/pkg/signalapi"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

type signalLog struct {
	ID        string     `db:"id" json:"id"`
	Chain     string     `db:"chain" json:"chain"`
	Signal    SignalType `db:"signal" json:"signal"`
	Strategy  string     `db:"strategy" json:"strategy"`
	TradeTime time.Time  `db:"trade_time" json:"trade_time"`
	CreatedAt null.Time  `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt null.Time  `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt null.Time  `db:"deleted_at" json:"deleted_at,omitempty"`
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
