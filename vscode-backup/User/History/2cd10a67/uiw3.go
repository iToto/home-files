// Package signalsvc is the service that handles getting and processing trade signals
package signalsvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"yield/signal-logger/internal/wlog"

	"github.com/guregu/null"
)

type SVC interface {
	HelloWorld(ctx context.Context, wl wlog.Logger) error
}

type helloService struct {
	// add any dependencies here (DB, Client, etc.)
}

func New(
// pass in any dependencies
) (SVC, error) {
	return &helloService{}, nil
}

func (ss *signalService) queryLastSignal(
	ctx context.Context,
	wl wlog.Logger,
	ip string,
) (*signalHistory, error) {
	var lastSignal signalHistory

	query := "SELECT id, ip, chain, signal, trade_time, created_at, updated_at, deleted_at FROM logger_signal_log WHERE ip = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1"
	err := ss.db.Get(&lastSignal, query, ip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoSignalHistory
		}

		wl.Error(fmt.Errorf("error querying signal: %w", err))
		return nil, ErrDBConnection
	}

	return &lastSignal, nil
}

func (ss *signalService) insertLatestSignalTradedForStrategy(
	ctx context.Context,
	wl wlog.Logger,
	signal signalHistory,
) error {
	signal.CreatedAt = null.NewTime(time.Now(), true)
	query := `INSERT INTO logger_signal_log (ip, chain, signal, trade_time, created_at) 
	VALUES (:ip, :chain, :signal, :trade_time, :created_at)`
	rows, err := ss.db.NamedQuery(query, signal)
	if err != nil {
		return err
	}

	rows.Close()

	return nil
}
