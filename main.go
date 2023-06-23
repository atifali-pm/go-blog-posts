package main

import (
	"net/http"

	"github.com/atifali-pm/go-blog-posts/config"
	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/atifali-pm/go-blog-posts/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.SetupDB()
	if err != nil {
		panic(err)
	}

	if err := models.Migrate(db); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router := routes.SetupRoutes(db)

	// Custom 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "product_not_found",
				"text":       "Product not found",
			},
			"body": nil,
		})
	})

	router.Run(":8000")
}
