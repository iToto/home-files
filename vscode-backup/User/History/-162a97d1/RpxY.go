package socialsvc

import (
	"context"
	"fmt"
	"social-links-api/internal/entities"
	"social-links-api/internal/wlog"

	"github.com/oklog/ulid/v2"
)

func (ss *socialService) CreateSocialURL(
	ctx context.Context,
	wl wlog.Logger,
	socialNetworks []entities.SocialURL,
) (*entities.SocialMediaLink, error) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint:gosec
	id, err := ulid.New(uilid.Now(), rng)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %s, err")
	}
	links := entities.SocialMediaLink{
		SocialURLs: socialNetworks,

	return nil, nil
}
