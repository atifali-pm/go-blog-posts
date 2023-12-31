package routes

import (
	"github.com/atifali-pm/go-blog-posts/controllers"
	"github.com/atifali-pm/go-blog-posts/middlewares"
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
		user.POST("/users/signup", controllers.Signup)
		user.POST("/users/login", controllers.Login)
	}

	post := r.Group("/api/v1")
	post.Use(middlewares.JWTMiddleware())
	{
		post.GET("/posts", controllers.GetPosts)
		post.GET("/posts/:post_id", controllers.GetPost)
		post.POST("/posts", controllers.CreatePost)
		post.PUT("/posts/:post_id", controllers.UpdatePost)
		post.DELETE("/posts/:post_id", controllers.DeletePost)
	}

	review := r.Group("/api/v1")
	review.Use(middlewares.JWTMiddleware())
	{
		review.GET("/reviews", controllers.GetReviews)
		review.POST("/reviews", controllers.CreateReview)
		review.GET("/reviews/:review_id", controllers.GetReview)
		review.PUT("/reviews/:review_id", controllers.UpdateReview)
		review.DELETE("/reviews/:review_id", controllers.DeleteReview)
	}

	return r
}
