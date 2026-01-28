package model

// Tour represents a travel product
type Tour struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PriceCents  int    `json:"price_cents"`
}
