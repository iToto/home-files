package hellosvc

import "errors"

var (
	ErrNoOpCurrentPosition = errors.New("current position is already desired position: no-op")
	ErrDBConnection        = errors.New("error with DB connection")
)
