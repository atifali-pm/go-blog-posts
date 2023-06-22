package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func SetupRoutes(db *gorm.DB) *gin.Engine{

	r := gin.Default()

	r.Use(func (c *gin.Context)  {
		c.Set("db", db)
		c.Next()
	})

	v1 := r.Group("/api/v1"){
		v1.GET("/users", controllers.ListUsers)
	}

	return r
}