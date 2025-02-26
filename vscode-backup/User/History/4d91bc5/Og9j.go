package handler

import (
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
}
