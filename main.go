package main

import (
	"github.com/atifali-pm/go-blog-posts/config"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.SetupDB()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router := router.SetupRoutes(db)
	router.Run(":8000")
}
