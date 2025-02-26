package entities

type SocialMediaLink struct {
	ID         string      `json:"id"`
	ShortCode  string      `json:"short_code"`
	SocialURLs []SocialURL `json:"social_urls"`
	CreatedAt  DateTime    `json:"created_at"`
	UpdatedAt  DateTime    `json:"updated_at"`
	DeletedAt  DateTime    `json:"deleted_at"`
}
