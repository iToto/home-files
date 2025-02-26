// Package hellosvc is the service that handles getting and processing trade signals
package hellosvc

import (
	"context"
	"yield/signal-logger/internal/wlog"
)

type SVC interface {
	// declare any methods that this service will expose
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
