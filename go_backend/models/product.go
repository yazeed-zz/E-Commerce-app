package models

import "time"

type Product struct {
	ID          string    `json:"id"`
	SellerID    string    `json:"seller_id"`
	CategoryID  string    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Images      []string  `json:"images"`
	RatingAvg   float64   `json:"rating_avg"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
