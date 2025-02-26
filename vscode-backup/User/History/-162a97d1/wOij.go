package socialsvc

import (
	"context"
	"social-links-api/internal/entities"
	"social-links-api/internal/wlog"
)

func (ss *socialService) CreateSocialURL(
	ctx context.Context,
	wl wlog.Logger,
	socialNetworks []entities.SocialURL,
) (*entities.SocialMediaLink, error) {
	return nil, nil
}
