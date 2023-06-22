package controllers

import (
	"net/http"
	"strconv"

	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset := (page - 1) * limit

	var users []models.User

	db.Offset(offset).Limit(limit).Order("first_name").Find(&users)

	c.JSON(http.StatusOK, users)

}
