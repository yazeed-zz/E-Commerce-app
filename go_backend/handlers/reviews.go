package handlers

import (
	"ecommerce/backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	productID := c.Param("id")

	var input struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate rating
	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	// Check if user purchased this product
	var count int
	config.DB.QueryRow(
		`SELECT COUNT(*) FROM order_items oi
		 JOIN orders o ON oi.order_id = o.id
		 WHERE o.user_id = $1 AND oi.product_id = $2`,
		userID, productID,
	).Scan(&count)

	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only review products you have purchased"})
		return
	}

	// Check if already reviewed
	var existing int
	config.DB.QueryRow(
		`SELECT COUNT(*) FROM reviews 
		 WHERE user_id = $1 AND product_id = $2`,
		userID, productID,
	).Scan(&existing)

	if existing > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reviewed this product"})
		return
	}

	// Insert review
	reviewID := uuid.New()
	_, err := config.DB.Exec(
		`INSERT INTO reviews (id, user_id, product_id, rating, comment)
		 VALUES ($1, $2, $3, $4, $5)`,
		reviewID, userID, productID, input.Rating, input.Comment,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit review"})
		return
	}

	// Update product average rating
	config.DB.Exec(
		`UPDATE products SET rating_avg = (
			SELECT AVG(rating) FROM reviews WHERE product_id = $1
		) WHERE id = $1`, productID,
	)

	c.JSON(http.StatusCreated, gin.H{"message": "Review submitted successfully"})
}

func GetReviews(c *gin.Context) {
	productID := c.Param("id")

	rows, err := config.DB.Query(
		`SELECT r.id, u.name, r.rating, r.comment, r.created_at
		 FROM reviews r
		 JOIN users u ON r.user_id = u.id
		 WHERE r.product_id = $1
		 ORDER BY r.created_at DESC`, productID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer rows.Close()

	var reviews []gin.H
	for rows.Next() {
		var id, name, comment, createdAt string
		var rating int

		rows.Scan(&id, &name, &rating, &comment, &createdAt)
		reviews = append(reviews, gin.H{
			"id":         id,
			"user":       name,
			"rating":     rating,
			"comment":    comment,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
