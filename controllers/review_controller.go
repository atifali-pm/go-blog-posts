package controllers

import (
	"net/http"
	"strconv"

	"github.com/atifali-pm/go-blog-posts/helpers"
	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetReviews(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var reviews []models.Review
	var total int64

	db.Offset(offset).Limit(limit).Order("title").Find(&reviews)
	db.Model(&models.Review{}).Count(&total)

	meta := helpers.GeneratePaginationMeta(page, limit, offset, int(total))
	links := helpers.GeneratePaginationLinks(c.Request, meta.LastPage, page)

	response := gin.H{
		"body": gin.H{
			"reviews": reviews,
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

func GetReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id, _ := strconv.Atoi(c.Param("review_id"))

	var review models.Review

	if err := db.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "review_not_found",
				"text":       "Review not found",
			},
			"body": nil,
		})
		return
	}

	// Create the review response object
	reviewResponse := models.ReviewResponse{
		ID:          review.ID,
		Title:       review.Title,
		Description: review.Description,
		UserID:      review.UserID,
		PostID:      review.PostID,
	}

	var user models.User
	db.Model(&review).Association("User").Find(&user)
	// Create the user response object
	userResponse := models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Phone:     user.Phone,
	}

	var post models.Post
	db.Model(&review).Association("Post").Find(&post)
	// Create the post response object
	postResponse := models.PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		UserID:      post.UserID,
	}

	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": ReviewBodyResponse{
			Review: reviewResponse,
			Post:   postResponse,
			User:   userResponse,
		},
	}

	c.JSON(http.StatusOK, response)
}

func CreateReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var review models.Review

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func UpdateReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("review_id"))

	var review models.Review

	if err := db.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "review_not_found",
				"text":       "Review not found",
			},
			"body": nil,
		})
		return
	}

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&review)

	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": review,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(c.Param("review_id"))

	var review models.Review

	if err := db.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": gin.H{
				"code":       http.StatusNotFound,
				"error":      true,
				"error_type": "review_not_found",
				"text":       "Review not found",
			},
			"body": nil,
		})
		return
	}

	db.Delete(&review)

	// Create the response object
	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": review,
	}

	c.JSON(http.StatusOK, response)
}

type ReviewBodyResponse struct {
	Review models.ReviewResponse `json:"review"`
	Post   models.PostResponse   `json:"post"`
	User   models.UserResponse   `json:"user"`
}
