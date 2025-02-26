package entities

import "time"

type SocialMediaLink struct {
	ID         string      `json:"id"`
	ShortCode  string      `json:"short_code"`
	SocialURLs []SocialURL `json:"social_urls"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  time.Time   `json:"deleted_at"`
}
