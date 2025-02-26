package handler

import (
	"encoding/json"
	"net/http"
	"social-links-api/internal/entities"
	"social-links-api/internal/service/socialsvc"
	"social-links-api/internal/wlog"
)

type socialResponse struct {
	URL string `json:"url"`
}

type socialRequest struct {
	URLS []entities.SocialURL `json:"urls"`
}

func CreateSocialURL(wl wlog.Logger, socialService socialsvc.SVC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := socaiclRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.BadRequest(ctx, wl, w, render.ErrJSONDecode)
			return
		}

		data, err := socialService.CreateSocialURL(ctx, wl, req.URLS)
		if err != nil {
			render.InternalError(ctx, wl, w, render.InternalError)
			return
		}
}
