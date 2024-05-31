package domain

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
	Items       []T
}
