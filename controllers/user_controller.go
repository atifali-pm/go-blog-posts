package controllers

import (
	"net/http"
	"strconv"

	"github.com/atifali-pm/go-blog-posts/helpers"
	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginationMeta struct {
	CurrentPage int    `json:"current_page"`
	From        int    `json:"from"`
	LastPage    int    `json:"last_page"`
	Path        string `json:"path"`
	PerPage     int    `json:"per_page"`
	To          int    `json:"to"`
	Total       int    `json:"total"`
}

type PaginationLinks struct {
	First string      `json:"first"`
	Last  string      `json:"last"`
	Prev  interface{} `json:"prev"`
	Next  interface{} `json:"next"`
}

type UserListResponse struct {
	Data  []models.User   `json:"data"`
	Meta  PaginationMeta  `json:"meta"`
	Links PaginationLinks `json:"links"`
}

func ListUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	db.Offset(offset).Limit(limit).Order("first_name").Find(&users)
	db.Model(&models.User{}).Count(&total)

	response := UserListResponse{
		Data: users,
		Meta: PaginationMeta{
			CurrentPage: page,
			From:        offset + 1,
			LastPage:    int(total/int64(limit)) + 1,
			Path:        c.Request.URL.Path,
			PerPage:     limit,
			To:          offset + len(users),
			Total:       int(total),
		},
		Links: PaginationLinks{
			First: helpers.GenerateURL(c.Request, 1),
			Last:  helpers.GenerateURL(c.Request, int(total/int64(limit))+1),
			Prev:  helpers.GeneratePrevURL(c.Request, page),
			Next:  helpers.GenerateNextURL(c.Request, int(total/int64(limit))+1, page),
		},
	}

	c.JSON(http.StatusOK, response)

}
