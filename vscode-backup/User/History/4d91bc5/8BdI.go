package handler

import (
	"net/http"
	"social-links-api/internal/service/socialsvc"
	"social-links-api/internal/wlog"
	"social-links-api/pkg/render"

	"github.com/gorilla/mux"
)

type socialResponse struct {
	URL string `json:"url"`
}

func CreateSocialURL(wl wlog.Logger, socialService socialsvc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		shortCode := vars["short_code"]

		social, err := socialService.CreateSocialURL(ctx, wl, shortCode)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
		}

		resp := socialResponse{
			URL: social.URL,
		}

		render.JSON(ctx, wl, w, resp, http.StatusOK)
	}
}
