package domain

import (
	"strings"
)

type SearchQuery struct {
	Page      int
	PerPage   int
	Term      string
	Sort      string
	Direction string
}

type Pagination[T any] struct {
	CurrentPage int
	PerPage     int
	Total       int
	Items       []*T
}

func (sq SearchQuery) SortDirection() string {
	if sq.Sort == "DESC" || sq.Sort == "ASC" {
		return sq.Sort
	}
	return "ASC"
}

func (sq SearchQuery) Limit() int {
	return sq.PerPage
}

func (sq SearchQuery) Offset() int {
	return (sq.Page - 1) * sq.PerPage
}

func (sq SearchQuery) SortColumn() string {
	for _, safeValue := range []string{"name", "description"} {
		if sq.Sort == safeValue {
			return strings.Trim(sq.Sort, " ")
		}
	}
	panic("unsafe sort parameter " + sq.Sort)
}
