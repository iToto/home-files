package hellosvc

import (
	"context"
	"vino-amici/internal/wlog"
)

func (hs *helloService) HelloWorld(ctx context.Context, wl wlog.Logger) error {
	wl.Info("Hello World!")
	return nil
}
