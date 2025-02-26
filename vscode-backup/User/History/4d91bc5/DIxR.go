package handler

import (
	"net/http"
	"social-links-api/internal/service/socialsvc"
	"social-links-api/internal/wlog"
)

type socialResponse struct {
	URL string `json:"url"`
}

func CreateSocialURL(wl wlog.Logger, socialService socialsvc.SVC) http.HandlerFunc {

}
