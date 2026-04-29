package store

import (
	"net/http"
	"strconv"
)

type PaginationFeedsQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (q *PaginationFeedsQuery) Parse(r *http.Request) (*PaginationFeedsQuery, error) {
	qs := r.URL.Query()
	limit := qs.Get("limit")

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
		q.Limit = limitInt
	}

	offset := qs.Get("offset")
	if offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return nil, err
		}
		q.Offset = offsetInt
	}

	sort := qs.Get("sort")
	if sort != "" {
		q.Sort = sort
	}

	return q, nil
}
