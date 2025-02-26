package socialsvc

import (
	"boilerplate-go-api/internal/entities"
	"boilerplate-go-api/internal/wlog"
	"context"

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
	// DB *sqlx.DB
}

func New(db *sqlx.DB) (SVC, error) {
	return &socialService{}, nil
}
