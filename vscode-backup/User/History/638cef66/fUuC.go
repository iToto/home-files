// Package tickr is for handlers that are triggered by a ticker
package tickr

import (
	"context"
	"fmt"
	"time"
	"yield/signal-logger/internal/entities"
	"yield/signal-logger/internal/service/signalsvc"
	"yield/signal-logger/internal/wlog"
)

const (
	defaultBackoff int64 = 15
)

type TIC interface {
	Run(ctx context.Context, wl wlog.Logger)
}

type signalTicker struct {
	res           int64
	backoff       int64
	retryCount    int64
	signals       []entities.SignalSource
	signalService signalsvc.SVC
}

func New(res int64, signals []entities.SignalSource, ss signalsvc.SVC) (TIC, error) {
	return &signalTicker{
		res:           res,
		signals:       signals,
		signalService: ss,
		backoff:       defaultBackoff,
		retryCount:    0,
	}, nil
}

func (st *signalTicker) Run(ctx context.Context, wl wlog.Logger) {
	go st.run(ctx, wl)
}

func (st *signalTicker) run(ctx context.Context, wl wlog.Logger) {
	ticker := time.NewTicker(time.Duration(st.res) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wl := wlog.WithServiceRequest(ctx, wl, "signal-ticker")
			wl.Debug("starting process")
			err := st.signalService.GetAndProcessSignal(ctx, wl, st.signals)
			if err != nil {
				wl.Error(fmt.Errorf("stopping ticker: error when processing signal: %w", err))
				ticker.Stop()
			}
			wl.Debug("finished process")
		}
	}

}
