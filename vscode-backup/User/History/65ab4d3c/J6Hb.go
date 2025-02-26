package socialsvc

import (
	"context"
	"social-links-api/internal/entities"
	"social-links-api/internal/wlog"

	"github.com/jmoiron/sqlx"
)

type SVC interface {
	CreateSocialURL(
		ctx context.Context,
		wl wlog.Logger,
		socialNetworks []entities.SocialURL,
	) (*entities.SocialMediaLink, error)
}

type socialService struct {
	db       *sqlx.DB
	base_url string
}

func New(d *sqlx.DB) (SVC, error) {
	return &socialService{
		db: d,
	}, nil
}
