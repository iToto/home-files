package socialsvc

import (
	"context"
	"fmt"
	"math/rand"
	"social-links-api/internal/entities"
	"social-links-api/internal/wlog"
	"time"

	"github.com/oklog/ulid/v2"
)

func (ss *socialService) CreateSocialURL(
	ctx context.Context,
	wl wlog.Logger,
	socialNetworks []entities.SocialURL,
) (*entities.SocialMediaLink, error) {

	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // nolint:gosec
	id, err := ulid.New(ulid.Now(), rng)
	if err != nil {
		return nil, fmt.Errorf("error generating id: %s, err")
	}

	link := entities.SocialMediaLink{
		ID:         id.String(),
		ShortCode:  id.String(),
		SocialURLs: socialNetworks,
		URL:        "https://" + ss.base_url + "/" + id.String(),
	}

	return &link, nil
}
