package domain

import (
	"errors"
	"slices"
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
	CurrentPage int  `json:"currentPage"`
	PerPage     int  `json:"perPage"`
	Total       int  `json:"total"`
	IsLast      bool `json:"isLast"`
	Items       []*T `json:"items"`
}

var safeValues = []string{"name", "description"}

func (sq SearchQuery) SortDirection() string {
	if sq.Direction == "" {
		return "ASC"
	}
	return sq.Direction
}

func (sq SearchQuery) Limit() int {
	return sq.PerPage
}

func (sq SearchQuery) Offset() int {
	return (sq.Page - 1) * sq.PerPage
}

func (sq SearchQuery) SortColumn() string {
	if sq.Sort == "" {
		return "name"
	}
	return sq.Sort
}

func (sq *SearchQuery) Validate() error {
	if sq.Page < 1 {
		return errors.New("invalid page")
	}
	if sq.PerPage < 1 {
		return errors.New("perPage should be greater than zero")
	}

	if sq.Sort != "" && !slices.Contains(safeValues, sq.Sort) {
		return errors.New("can only sort by 'name' and 'description'")
	}
	if sq.Direction != "" && !slices.Contains([]string{"ASC", "DESC"}, strings.ToUpper(sq.Direction)) {
		return errors.New("invalid direction")
	}

	return nil
}
