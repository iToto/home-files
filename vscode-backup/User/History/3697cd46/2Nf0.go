package exchangeclient

import (
	"context"
	"yield-mvp/internal/wlog"
)

type Client interface {
	GetBalance(ctx context.Context, wl wlog.Logger) error
	GetPosition(ctx context.Context, wl wlog.Logger) error
}
