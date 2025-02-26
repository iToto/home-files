package signalsvc

import "errors"

var (
	ErrNoOpCurrentPosition = errors.New("current position is already desired position: no-op")
	ErrDBConnection        = errors.New("error with DB connection")
	ErrNoSignalHistory     = errors.New("no signal history found")
)
