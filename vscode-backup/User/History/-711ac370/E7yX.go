package entities

type Luggage struct {
	ID             string   `json:"id"`
	Manufacturer   string   `json:"manufacturer"`
	Color          string   `json:"color"`
	Image          string   `json:"image"`
	SerialNumber   string   `json:"serial_number"`
	Owner          Owner    `json:"owner"`
	CurrentFlights []Flight `json:"current_flights"`
}

type Owner struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	TwitterHandle   string `json:"twitter_handle"`
	LinkedInProfile string `json:"linkedin_profile"`
}

type Flight struct {
	ID           string    `json:"id"`
	FlightNumber string    `json:"flight_number"`
	Origin       string    `json:"origin"`
	Destination  string    `json:"destination"`
	Departure    Time.Time `json:"departure"`
	Arrival      string    `json:"arrival"`
}
