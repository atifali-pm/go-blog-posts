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

	user := r.Group("/api/v1")
	{
		user.GET("/users", controllers.ListUsers)
		user.GET("/users/:user_id", controllers.GetUser)
		user.POST("/users", controllers.CreateUser)
		user.PUT("/users/:user_id", controllers.UpdateUser)
		user.DELETE("/users/:user_id", controllers.DeleteUser)
	}

	post := r.Group("/api/v1")
	{
		post.GET("/posts", controllers.ListPosts)
	}

	return r
}
