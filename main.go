package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/atifali-pm/go-blog-posts/config"
	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/atifali-pm/go-blog-posts/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Invalid port number")
	}

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

	router.Run(":" + strconv.Itoa(port))
}
