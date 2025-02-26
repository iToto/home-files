package socialsvc

import (
	"boilerplate-go-api/internal/entities"
	"boilerplate-go-api/internal/wlog"
	"context"
)

type SVC interface {
	CreateSocialURL(ctx context.Context, wl wlog.Logger, socialNetworks []entities.SocialURL) (*entities.SocialMediaLink, error)
}
