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

	// Fetch posts for each user
	for i := range users {
		var posts []models.Post
		db.Model(&users[i]).Association("Posts").Find(&posts)
		users[i].Posts = posts
	}

	meta := helpers.GeneratePaginationMeta(page, limit, offset, int(total))
	links := helpers.GeneratePaginationLinks(c.Request, meta.LastPage, page)

	response := gin.H{
		"body": gin.H{
			"users": users,
		},
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"meta":  meta,
		"links": links,
	}

	c.JSON(http.StatusOK, response)

}

func GetUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id, _ := strconv.Atoi(c.Param("user_id"))

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "user_not_found",
				"text":       "User not found",
			},
			"body": nil,
		})
		return
	}

	var posts []models.Post
	db.Model(&user).Association("Posts").Find(&posts)

	// Create a separate struct for the response body
	type UserPostsBody struct {
		User  models.User   `json:"user"`
		Posts []models.Post `json:"posts"`
	}

	// Create the response object
	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": UserPostsBody{
			User:  user,
			Posts: posts,
		},
	}

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
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusBadRequest,
				"error":      true,
				"error_type": "user_not_created",
				"text":       "User not created",
			},
			"body": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, user)

}

func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("user_id"))

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "user_not_found",
				"text":       "User not found",
			},
			"body": nil,
		})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&user)

	// Create the response object
	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": user,
	}

	c.JSON(http.StatusOK, response)

}

func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("user_id"))

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "user_not_found",
				"text":       "User not found",
			},
			"body": nil,
		})
		return
	}

	db.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
