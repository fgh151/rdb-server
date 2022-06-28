package utils

import (
	"net/http"
	"strconv"
)

func GetPagination(r *http.Request) (int, int, string, string) {

	v := r.URL.Query()

	limit, err := strconv.Atoi(v.Get("_end"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(v.Get("_start"))
	if err != nil {
		offset = 0
	}

	order := v.Get("_order")
	if order == "" {
		order = "ASC"
	}
	sort := v.Get("_sort")
	if sort == "" {
		sort = "id"
	}

	return limit - offset, offset, order, sort
}
