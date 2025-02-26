package handler

type response struct {
	IsOnAir     bool   `json:"is_on_air,omitempty"`
	LastUpdated string `json:"last_updated,omitempty"`
	LastOnAir   string `json:"last_on_air,omitempty"`
}
