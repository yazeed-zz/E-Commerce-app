package handlers

import (
	"ecommerce/backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Get cart items
	rows, err := config.DB.Query(
		`SELECT ci.product_id, ci.quantity, p.price, p.stock
		 FROM cart_items ci
		 JOIN products p ON ci.product_id = p.id
		 WHERE ci.user_id = $1`, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}
	defer rows.Close()

	type CartItem struct {
		ProductID string
		Quantity  int
		Price     float64
		Stock     int
	}

	var items []CartItem
	var total float64

	for rows.Next() {
		var item CartItem
		rows.Scan(&item.ProductID, &item.Quantity, &item.Price, &item.Stock)

		if item.Quantity > item.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product " + item.ProductID})
			return
		}

		total += item.Price * float64(item.Quantity)
		items = append(items, item)
	}

	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Create order
	orderID := uuid.New()
	_, err = config.DB.Exec(
		`INSERT INTO orders (id, user_id, status, total_amount) 
		 VALUES ($1, $2, $3, $4)`,
		orderID, userID, "pending", total,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order items & update stock
	for _, item := range items {
		itemID := uuid.New()
		config.DB.Exec(
			`INSERT INTO order_items (id, order_id, product_id, quantity, price_at_purchase)
			 VALUES ($1, $2, $3, $4, $5)`,
			itemID, orderID, item.ProductID, item.Quantity, item.Price,
		)

		// Reduce stock
		config.DB.Exec(
			`UPDATE products SET stock = stock - $1 WHERE id = $2`,
			item.Quantity, item.ProductID,
		)
	}

	// Clear cart
	config.DB.Exec(
		`DELETE FROM cart_items WHERE user_id = $1`, userID,
	)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Order placed successfully",
		"order_id": orderID,
		"total":    total,
	})
}

func GetOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := config.DB.Query(
		`SELECT id, status, total_amount, created_at 
		 FROM orders WHERE user_id = $1
		 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []gin.H
	for rows.Next() {
		var id, status, createdAt string
		var total float64

		rows.Scan(&id, &status, &total, &createdAt)
		orders = append(orders, gin.H{
			"order_id":   id,
			"status":     status,
			"total":      total,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
