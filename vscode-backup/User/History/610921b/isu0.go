// Package tickr is for handlers that are triggered by a ticker
package tickr

import (
	"context"
	"errors"
	"fmt"
	"time"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/service/signalsvc"
	"yield-mvp/internal/wlog"
)

const (
	defaultBackoff int64 = 60
)

type TIC interface {
	Run(ctx context.Context, wl wlog.Logger)
}

type signalTicker struct {
	res           int64
	backoff       int64
	retryCount    int
	chain         entities.ChainType
	strats        []entities.Strategy
	signalService signalsvc.SVC
}

func New(res int64, chain entities.ChainType, strats []entities.Strategy, ss signalsvc.SVC) (TIC, error) {
	return &signalTicker{
		res:           res,
		chain:         chain,
		strats:        strats,
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
	done := make(chan bool)

	wl := wlog.WithServiceRequest(ctx, wl, "signal-ticker")
	wl = wlog.WithChain(wl, string(st.chain))
	for {
		select {
		case <-done:

		case <-ticker.C:
			// check BTC Signal
			wl.Debug("starting process")
			err := st.signalService.GetAndProcessSignal(ctx, wl, st.strats, st.chain)
			if err != nil {
				if errors.Is(err, signalsvc.ErrDBConnection) {
					wl.Error(fmt.Errorf("stopping ticker: error when processing signal: %w", err))
					ticker.Stop()
				}
				if errors.Is(err, signalsvc.ErrSignalClient) {
					st.retryCount += 1
					wl.Error(fmt.Errorf("error when calling signal client, pausing ticker"))
					ticker.Reset(time.Duration(st.backoff) * time.Second * time.Duration(st.retryCount))
				}
			}
			wl.Debug("finished process")
		}
	}

}
