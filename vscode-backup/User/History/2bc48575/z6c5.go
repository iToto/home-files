package hellosvc

import (
	"boilerplate-go-api/internal/wlog"
	"context"
)

func (hs *helloService) HelloWorld(ctx context.Context, wl wlog.Logger) error {
	wl.Info("Hello World!")
	return nil
}
