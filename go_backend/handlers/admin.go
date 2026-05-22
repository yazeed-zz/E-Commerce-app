package handlers

import (
	"ecommerce/backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetStats(c *gin.Context) {
	var totalUsers, totalProducts, totalOrders int
	var totalRevenue float64

	config.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	config.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalProducts)
	config.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&totalOrders)
	config.DB.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status != 'cancelled'").Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"total_users":    totalUsers,
		"total_products": totalProducts,
		"total_orders":   totalOrders,
		"total_revenue":  totalRevenue,
	})
}

func GetUsers(c *gin.Context) {
	rows, err := config.DB.Query(
		`SELECT id, name, email, role, created_at 
		 FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var id, name, email, role, createdAt string
		rows.Scan(&id, &name, &email, &role, &createdAt)
		users = append(users, gin.H{
			"id":         id,
			"name":       name,
			"email":      email,
			"role":       role,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func CreateProduct(c *gin.Context) {
	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		CategoryID  string  `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New()
	_, err := config.DB.Exec(
		`INSERT INTO products (id, name, description, price, stock, category_id)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		id, input.Name, input.Description, input.Price, input.Stock, input.CategoryID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"id":      id,
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	_, err := config.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[string]bool{
		"pending": true, "paid": true,
		"shipped": true, "delivered": true, "cancelled": true,
	}

	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	_, err := config.DB.Exec(
		"UPDATE orders SET status = $1 WHERE id = $2",
		input.Status, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
