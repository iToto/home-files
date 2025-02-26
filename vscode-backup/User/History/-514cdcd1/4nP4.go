package handler

import (
	"boilerplate-go-api/internal/service/hellosvc"
	"boilerplate-go-api/internal/wlog"
	"boilerplate-go-api/pkg/render"
	"net/http"
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

		resp :=

			render.JSON(ctx, wl, w, nil, http.StatusOK)
	}
}
