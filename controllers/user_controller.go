package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/atifali-pm/go-blog-posts/helpers"
	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

		var reviews []models.Review
		db.Model(&users[i]).Association("Reviews").Find(&reviews)
		users[i].Reviews = reviews

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

	var reviews []models.Review
	db.Model(&user).Association("Reviews").Find(&reviews)

	// Create a separate struct for the response body
	type UserPostsBody struct {
		User    models.User     `json:"user"`
		Posts   []models.Post   `json:"posts"`
		Reviews []models.Review `json:"reviews"`
	}

	// Create the response object
	response := gin.H{
		"status": gin.H{
			"code":  http.StatusOK,
			"error": false,
			"text":  "success",
		},
		"body": UserPostsBody{
			User:    user,
			Posts:   posts,
			Reviews: reviews,
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

type SignupRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone"`
}

func Signup(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Avatar:    req.Avatar,
		Phone:     req.Phone,
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := gin.H{
		"status": gin.H{
			"code":  http.StatusCreated,
			"error": false,
			"text":  "success",
		},
		"body": user,
	}

	c.JSON(http.StatusCreated, response)

}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate the access token
	accessToken := generateToken(user.ID, time.Hour*1)

	// Generate the refresh token
	refreshToken := generateToken(user.ID, time.Hour*24*7)

	response := gin.H{
		"token_type":    "Bearer",
		"expires_in":    time.Hour * 1, // Set the expiration time for the access token
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	c.JSON(http.StatusOK, response)

}

func generateToken(userID uint, expiration time.Duration) string {
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(expiration).Unix(),
	})

	// Sign the token with the JWT secret
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString
}
