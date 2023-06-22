package controllers

import (
	"net/http"
	"strconv"

	"github.com/atifali-pm/go-blog-posts/helpers"
	"github.com/atifali-pm/go-blog-posts/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListPosts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var users []models.Post
	var total int64

	db.Offset(offset).Limit(limit).Order("title").Find(&users)
	db.Model(&models.Post{}).Count(&total)

	meta := helpers.GeneratePaginationMeta(page, limit, offset, int(total))
	links := helpers.GeneratePaginationLinks(c.Request, meta.LastPage, page)

	response := helpers.GenerateListResponse(users, meta, links)

	c.JSON(http.StatusOK, response)

}
