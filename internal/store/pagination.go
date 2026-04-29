package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginationFeedsQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Until  string   `json:"until" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
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

	tags := qs.Get("tags")
	if tags != "" {
		q.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		q.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		since, err := parseTime(since)
		if err != nil {
			return nil, err
		}
		q.Since = since
	}

	until := qs.Get("until")
	if until != "" {
		until, err := parseTime(until)
		if err != nil {
			return nil, err
		}
		q.Until = until
	}

	return q, nil
}

func parseTime(timeStr string) (string, error) {
	it, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return "", err
	}
	return it.Format(time.RFC3339), nil
}
