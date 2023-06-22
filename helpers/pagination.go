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
