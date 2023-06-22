package controllers

import (
	"net/http"
	"strconv"

	"github.com/atifali-pm/go-blog-posts/helpers"
	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	db.Offset(offset).Limit(limit).Order("first_name").Find(&users)
	db.Model(&models.User{}).Count(&total)

	meta := helpers.GeneratePaginationMeta(page, limit, offset, int(total))
	links := helpers.GeneratePaginationLinks(c.Request, meta.LastPage, page)

	response := helpers.GenerateListResponse(users, meta, links)

	c.JSON(http.StatusOK, response)

}

func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)

}
