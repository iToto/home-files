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

		lug := entities.Luggage{
			ID:           "01GTAWH1GYGMT0B56H3SJQ921W",
			Manufacturer: "Samsonite",
			Color:        "Black",
			Image:        "https://www.samsonite.com/on/demandware.static/-/Sites-samsonite-master-catalog/default/dw8b3b3b3a/images/hi-res/68302-1041.jpg",
			SerialNumber: "123456789",
			Owner: entities.Owner{
				ID:              "01GTAWHXMXBW0EMB0GAGTBGXTT",
				FirstName:       "John",
				LastName:        "Doe",
				Email:           "jdoe@email.com",
				Phone:           "123-456-7890",
				TwitterHandle:   "@jdoe",
				LinkedInProfile: "https://www.linkedin.com/in/jdoe",
			},
		}

		resp := luggageResponse{
			Luggage: lug,
		}

		render.JSON(ctx, wl, w, resp, http.StatusOK)
	}
}
