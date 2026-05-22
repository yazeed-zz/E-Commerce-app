package handlers

import (
	"ecommerce/backend/config"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetProducts(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")

	query := `SELECT id, name, description, price, stock, images, rating_avg 
			  FROM products WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if category != "" {
		query += ` AND category_id = $` + string(rune('0'+argCount))
		args = append(args, category)
		argCount++
	}

	if search != "" {
		query += ` AND name ILIKE $` + string(rune('0'+argCount))
		args = append(args, "%"+search+"%")
		argCount++
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []gin.H
	for rows.Next() {
		var id, name, description string
		var price float64
		var stock int
		var images []string
		var ratingAvg float64

		rows.Scan(&id, &name, &description, &price, &stock, &images, &ratingAvg)
		products = append(products, gin.H{
			"id":          id,
			"name":        name,
			"description": description,
			"price":       price,
			"stock":       stock,
			"images":      images,
			"rating_avg":  ratingAvg,
		})
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var productID, name, description string
	var price float64
	var stock int
	var ratingAvg float64

	err := config.DB.QueryRow(
		`SELECT id, name, description, price, stock, rating_avg 
		 FROM products WHERE id = $1`, id,
	).Scan(&productID, &name, &description, &price, &stock, &ratingAvg)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          productID,
		"name":        name,
		"description": description,
		"price":       price,
		"stock":       stock,
		"rating_avg":  ratingAvg,
	})
}

func AddToCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New()
	_, err := config.DB.Exec(
		`INSERT INTO cart_items (id, user_id, product_id, quantity) 
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id, product_id) 
		 DO UPDATE SET quantity = cart_items.quantity + $4`,
		id, userID, input.ProductID, input.Quantity,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to cart successfully"})
}

func GetCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := config.DB.Query(
		`SELECT p.id, p.name, p.price, ci.quantity 
		 FROM cart_items ci
		 JOIN products p ON ci.product_id = p.id
		 WHERE ci.user_id = $1`, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}
	defer rows.Close()

	var items []gin.H
	for rows.Next() {
		var id, name string
		var price float64
		var quantity int

		rows.Scan(&id, &name, &price, &quantity)
		items = append(items, gin.H{
			"product_id": id,
			"name":       name,
			"price":      price,
			"quantity":   quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{"cart": items})
}

func GetRecommendations(c *gin.Context) {
	recType := c.Param("type")
	id := c.Query("id")

	url := "http://127.0.0.1:8001/recommendations/" + recType
	if id != "" {
		url += "/" + id
	}

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Recommendation service unavailable"})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	c.JSON(http.StatusOK, result)
}
