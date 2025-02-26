package handler

import (
	"boilerplate-go-api/internal/service/hellosvc"
	"boilerplate-go-api/internal/wlog"
	"net/http"
)

func HelloWorld(wl wlog.Logger, helloService hellosvc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := helloService.HelloWorld(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}
	}
}
