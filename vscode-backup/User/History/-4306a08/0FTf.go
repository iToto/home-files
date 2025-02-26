// Package signalsvc is the service that handles getting and processing trade signals
package signalsvc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"yield/signal-logger/internal/entities"
	"yield/signal-logger/internal/wlog"

	"github.com/oklog/ulid"
)

// GetAndProcessSignal will ping for a specified signal, parse it and process it
func (ss *signalService) GetAndProcessSignal(
	ctx context.Context,
	wl wlog.Logger,
	signals []entities.SignalSource,
) error {

	for _, signal := range signals {
		rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint: gosec
		id, err := ulid.New(ulid.Now(), rng)
		if err != nil {
			return fmt.Errorf("error generating id: %s", err)
		}

		wl := wlog.WithSignalIP(wl, signal.IP)

		// check for history for signal
		prevSignal, err := ss.queryLastSignal(ctx, wl, signal.IP)
		if err != nil {
			if errors.Is(err, ErrDBConnection) {
				wl.Error(err)
				return err
			}

			if errors.Is(err, ErrNoSignalHistory) {
				wl.Infof("no history found")
			}
		}

		if signal.SignalVersion == entities.V1 {
			signalResp, err := ss.sc.GetSignalFromIPV1(ctx, wl, signal)
			if err != nil {
				wl.Error(fmt.Errorf("error getting signal: %w", err))
				continue
			}

			// convert to entity
			signalLog := entities.SignalLogV1{
				ID:            id.String(),
				Chain:         string(signal.Type),
				IP:            signal.IP,
				Signal:        string(signalResp.Signal),
				CurrentTime:   signalResp.CurrentTime,
				LastTrade:     string(signalResp.LastTrade),
				LastData:      signalResp.LastData,
				LastTradeTime: signalResp.LastTradeTime,
				CreatedAt:     time.Now(),
			}

			// check for delta
			if prevSignal != nil {
				if prevSignal.Signal.IsEquivalent(entities.SignalType(signalResp.LastTrade)) &&
					prevSignal.TradeTime.Equal(signalLog.LastTradeTime) {
					wl.Info("no signal change, skipping")
					continue
				}
			}

			// new signal found, let's record it for next time
			s := signalHistory{
				ID:        signalLog.ID,
				Chain:     signalLog.Chain,
				IP:        signalLog.IP,
				Signal:    entities.SignalType(strings.ToLower(signalLog.LastTrade)),
				TradeTime: signalLog.LastTradeTime,
			}

			err = ss.insertLatestSignalTradedForStrategy(ctx, wl, s)
			if err != nil {
				wl.Error(fmt.Errorf("unable to insert signal into log table: %w", err))
				return err
			}

			// log signal to BQ
			err = ss.sl.Log(ctx, wl, signalLog)
			if err != nil {
				wl.Error(fmt.Errorf("unable to log signal with error: %w", err))
				continue
			}
		} else {
			signalResp, err := ss.sc.GetSignalFromIPV2(ctx, wl, signal)
			if err != nil {
				wl.Error(fmt.Errorf("error getting signal: %w", err))
				continue
			}

			wl.Debugf("signal received: %+v", signalResp)

			// convert to entity
			signalLog := entities.SignalLogV2{
				ID:                  id.String(),
				Chain:               string(signal.Type),
				IP:                  signal.IP,
				FetchResultStatus:   signalResp.FetchResultStatus,
				FetchType:           signalResp.FetchType,
				StrategyState:       string(signalResp.StrategyState),
				StrategyVersion:     signalResp.StrategyVersion,
				LastChecked:         signalResp.LastChecked,
				LastTradeSignal:     string(signalResp.LastTradeSignal),
				LastTradeSignalTime: signalResp.LastTradeSignalTime,
				CreatedAt:           time.Now(),
			}

			// check for delta
			if prevSignal != nil {
				if prevSignal.Signal.IsEquivalent(entities.SignalType(signalResp.LastTradeSignal)) &&
					prevSignal.TradeTime.Equal(signalLog.LastTradeSignalTime) {
					wl.Info("no signal change, skipping")
					continue
				}
			}

			// new signal found, let's record it for next time
			s := signalHistory{
				ID:        signalLog.ID,
				Chain:     signalLog.Chain,
				IP:        signalLog.IP,
				Signal:    entities.SignalType(signalResp.LastTradeSignal),
				TradeTime: signalLog.LastTradeSignalTime,
			}
			err = ss.insertLatestSignalTradedForStrategy(ctx, wl, s)
			if err != nil {
				wl.Error(fmt.Errorf("unable to insert signal into log table: %w", err))
				return err
			}
			wl.Debug("signal history saved")

			// log to BQ
			err = ss.slv2.Log(ctx, wl, signalLog)
			if err != nil {
				wl.Error(fmt.Errorf("unable to log signal v2 with error: %w", err))
				continue
			}
		}

		wl.Info("signal successfully logged")
	}
	return nil
}
