package routes

import (
	"github.com/atifali-pm/go-blog-posts/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", controllers.ListUsers)
		v1.POST("/users", controllers.CreateUser)
		v1.GET("/posts", controllers.ListPosts)
	}

	return r
}
