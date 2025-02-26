package signalsvc

import "errors"

var (
	ErrNoOpCurrentPosition = errors.New("no-op: current position is already desired position")
	ErrDBConnection        = errors.New("error with DB connection")
	ErrNoSignalHistory     = errors.New("no signal history found")
	ErrSignalClient        = errors.New("error received from Signal Client")
	ErrNoOpSignal          = errors.New("no-op: signal received cannot be processed")
	ErrClientBadVersion    = errors.New("unsupported signal version found")
	ErrClientBadSignalType = errors.New("invalid signal type found")
)
