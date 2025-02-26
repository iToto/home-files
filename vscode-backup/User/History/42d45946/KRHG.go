package handler

import (
	"net/http"
	"vino-amici/internal/service/hellosvc"
	"vino-amici/internal/wlog"
	"vino-amici/pkg/render"
)

type response struct {
	Hello string `json:"hello,omitempty"`
	Foo   string `json:"foo,omitempty"`
}

func HelloWorld(wl wlog.Logger, helloService hellosvc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := helloService.HelloWorld(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}

		resp := response{
			Hello: "World",
			Foo:   "Bar",
		}

		render.JSON(ctx, wl, w, resp, http.StatusOK)
	}
}
