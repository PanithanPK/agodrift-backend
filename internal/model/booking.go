package model

import "time"

type Booking struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	HotelID        int       `json:"hotel_id"`
	CheckIn        time.Time `json:"check_in"`
	CheckOut       time.Time `json:"check_out"`
	Adults         int       `json:"adults"`
	Children       int       `json:"children"`
	Rooms          int       `json:"rooms"`
	TotalPriceCents int      `json:"total_price_cents"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
