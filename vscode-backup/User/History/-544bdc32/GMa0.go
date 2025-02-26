package handler

import (
	"boilerplate-go-api/internal/entities"
	"boilerplate-go-api/internal/service/hellosvc"
	"boilerplate-go-api/internal/wlog"
	"boilerplate-go-api/pkg/render"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type luggageResponse struct {
	Luggage entities.Luggage `json:"luggage,omitempty"`
}

func GetLuggageData(wl wlog.Logger, helloService hellosvc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			render.BadRequest(ctx, wl, w, fmt.Errorf("missing required id in URL"))
			return
		}

		lug := entities.Luggage{
			ID:           id,
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
			CurrentFlights: []entities.Flight{
				{
					ID:               "01GTAWHXMXBW0EMB0GAGTBGXTT",
					FlightNumber:     "UA123",
					DepartureAirport: "SFO",
					ArrivalAirport:   "LAX",
					DepartureTime:    "2021-01-01T12:00:00Z",
					ArrivalTime:      "2021-01-01T14:00:00Z",
				},
			},
		}

		resp := luggageResponse{
			Luggage: lug,
		}

		render.JSON(ctx, wl, w, resp, http.StatusOK)
	}
}
