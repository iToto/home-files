package handler

import (
	"boilerplate-go-api/internal/entities"
	"boilerplate-go-api/internal/service/hellosvc"
	"boilerplate-go-api/internal/wlog"
	"boilerplate-go-api/pkg/render"
	"net/http"
)

type luggageResponse struct {
	Luggage entities.Luggage `json:"luggage,omitempty"`
}

func GetLuggageData(wl wlog.Logger, helloService hellosvc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := helloService.HelloWorld(ctx, wl)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}

		resp := luggageResponse{
			Hello: "World",
			Foo:   "Bar",
		}

		render.JSON(ctx, wl, w, resp, http.StatusOK)
	}
}
