package main

import (
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

	// defer db.Close()

	if err := models.MigrateUser(db); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router := routes.SetupRoutes(db)
	router.Run(":8000")
}
