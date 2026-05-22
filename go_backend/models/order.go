package models

import "time"

type Order struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	Status              string    `json:"status"`
	TotalAmount         float64   `json:"total_amount"`
	PaypalTransactionID string    `json:"paypal_transaction_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID              string  `json:"id"`
	OrderID         string  `json:"order_id"`
	ProductID       string  `json:"product_id"`
	Quantity        int     `json:"quantity"`
	PriceAtPurchase float64 `json:"price_at_purchase"`
}
