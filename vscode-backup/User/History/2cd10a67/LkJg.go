// Package signalsvc is the service that handles getting and processing trade signals
package signalsvc

import (
	"context"
	"yield/signal-logger/internal/signallogger"
	"yield/signal-logger/internal/signalloggerv2"
	"yield/signal-logger/internal/wlog"
	"yield/signal-logger/pkg/signalapi"

	"github.com/jmoiron/sqlx"
)

type SVC interface {
	HelloWorld(ctx context.Context, wl wlog.Logger) error
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
