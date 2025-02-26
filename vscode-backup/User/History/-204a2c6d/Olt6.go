// Package signalsvc is the service that handles getting and processing trade signals
package signalsvc

import (
	"context"
	"yield/signal-logger/internal/entities"
	"yield/signal-logger/internal/signallogger"
	"yield/signal-logger/internal/signalloggerv2"
	"yield/signal-logger/internal/wlog"
	"yield/signal-logger/pkg/signalapi"

	"github.com/jmoiron/sqlx"
)

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
