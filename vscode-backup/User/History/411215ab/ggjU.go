package signalsvc

import "errors"

var (
	ErrNoOpCurrentPosition = errors.New("current position is already desired position: no-op")
	ErrDBConnection        = errors.New("error with DB connection")
	ErrSignalClient        = errors.New("error received from Signal Client")
)
