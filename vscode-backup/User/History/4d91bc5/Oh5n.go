package handler

import (
	"encoding/json"
	"net/http"
	"social-links-api/internal/entities"
	"social-links-api/internal/service/socialsvc"
	"social-links-api/internal/wlog"
	"social-links-api/pkg/render"
)

type socialResponse struct {
	URL string `json:"url"`
}

type socialRequest struct {
	URLS []entities.SocialURL `json:"urls"`
}

func CreateSocialURL(wl wlog.Logger, socialService socialsvc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wl.Info("in social handler")
		ctx := r.Context()

		req := socialRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.BadRequest(ctx, wl, w, render.ErrJSONDecode)
			return
		}

		data, err := socialService.CreateSocialURL(ctx, wl, req.URLS)
		if err != nil {
			render.InternalError(ctx, wl, w, render.ErrInternal)
			return
		}

		resp := socialResponse{
			URL: data.URL,
		}

		render.JSON(ctx, wl, w, resp, http.StatusOK)
	}
}
