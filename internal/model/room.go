package model

// Room represents a hotel room
type Room struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Location           string  `json:"location"`
	Destination        string  `json:"destination"`
	Rating             float64 `json:"rating"`
	Reviews            int     `json:"reviews"`
	PriceCents         int     `json:"price_cents"`
	OriginalPriceCents *int    `json:"original_price_cents,omitempty"`
	Amenities          string  `json:"amenities"`
	Featured           bool    `json:"featured"`
	MaxAdults          int     `json:"max_adults"`
	MaxChildren        int     `json:"max_children"`
	RoomsTotal         int     `json:"rooms_total"`
	RoomsAvailable     int     `json:"rooms_available"`
	Status             string  `json:"status"`
}
