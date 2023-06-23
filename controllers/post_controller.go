package controllers

import (
	"net/http"
	"strconv"

	"github.com/atifali-pm/go-blog-posts/helpers"
	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPosts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var posts []models.Post
	var total int64

	db.Offset(offset).Limit(limit).Order("title").Find(&posts)
	db.Model(&models.Post{}).Count(&total)

	// Fetch posts for each user
	for i := range posts {
		db.Model(&posts[i]).Association("Users").Find(&posts[i].User)
	}

	meta := helpers.GeneratePaginationMeta(page, limit, offset, int(total))
	links := helpers.GeneratePaginationLinks(c.Request, meta.LastPage, page)

	response := gin.H{
		"body": gin.H{
			"posts": posts,
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

func GetPost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id, _ := strconv.Atoi(c.Param("post_id"))

	var post models.Post

	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "post_not_found",
				"text":       "Post not found",
			},
			"body": nil,
		})
		return
	}

	var users []models.User
	db.Model(&post).Association("User").Find(&users)

	// Create a separate struct for the response body
	type UserPostsBody struct {
		Post models.Post `json:"post"`
	}

	// Create the response object
	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": UserPostsBody{
			Post: post,
		},
	}

	c.JSON(http.StatusOK, response)

}

func CreatePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func UpdatePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("post_id"))

	var post models.Post

	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "post_not_found",
				"text":       "Post not found",
			},
			"body": nil,
		})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&post)

	// Create the response object
	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": post,
	}

	c.JSON(http.StatusOK, response)

}

func DeletePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("post_id"))

	var post models.Post

	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "post_not_found",
				"text":       "Post not found",
			},
			"body": nil,
		})
		return
	}

	db.Delete(&post)

	// Create the response object
	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": post,
	}

	c.JSON(http.StatusOK, response)

}
