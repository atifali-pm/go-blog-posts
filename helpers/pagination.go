package helpers

import (
	"net/http"
	"strconv"
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

type ListResponse struct {
	Data  interface{}     `json:"data"`
	Meta  PaginationMeta  `json:"meta"`
	Links PaginationLinks `json:"links"`
}

func GenerateURL(r *http.Request, page int) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + r.URL.Path + "?page=" + strconv.Itoa(page)
}

func GeneratePrevURL(r *http.Request, page int) interface{} {
	if page > 1 {
		return GenerateURL(r, page-1)
	}
	return nil
}

func GenerateNextURL(r *http.Request, lastPage int, page int) interface{} {
	if page < lastPage {
		return GenerateURL(r, page+1)
	}
	return nil
}

func GeneratePaginationMeta(page, limit, offset, total int) PaginationMeta {
	return PaginationMeta{
		CurrentPage: page,
		From:        offset + 1,
		LastPage:    int(total/limit) + 1,
		Path:        "",
		PerPage:     limit,
		To:          offset + limit,
		Total:       total,
	}
}

func GenerateListResponse(data interface{}, meta PaginationMeta, links PaginationLinks) ListResponse {
	return ListResponse{
		Data:  data,
		Meta:  meta,
		Links: links,
	}
}

func GeneratePaginationLinks(r *http.Request, lastPage, page int) PaginationLinks {
	return PaginationLinks{
		First: GenerateURL(r, 1),
		Last:  GenerateURL(r, lastPage),
		Prev:  GeneratePrevURL(r, page),
		Next:  GenerateNextURL(r, lastPage, page),
	}
}
