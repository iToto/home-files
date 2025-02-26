package hellosvc

import (
	"context"
	"social-links-api/internal/wlog"
)

func (hs *helloService) HelloWorld(ctx context.Context, wl wlog.Logger) error {
	wl.Info("Hello World!")
	return nil
}
