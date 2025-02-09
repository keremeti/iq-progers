package entity

import (
	"math"
)

type Pagination[T any] struct {
	Slice         []T
	HasNext       bool
	HasPrevious   bool
	RecordPerPage int
	CurrentPage   int
	TotalPage     int
	TotalCount    int
}

func NewPagination[T any](s []T, recordPerPage, page, total int) Pagination[T] {
	tp := int(math.Ceil(float64(total) / float64(recordPerPage)))
	return Pagination[T]{
		Slice:         s,
		HasNext:       total > page,
		HasPrevious:   page > 1 && tp != 0,
		RecordPerPage: recordPerPage,
		CurrentPage:   page,
		TotalPage:     tp,
		TotalCount:    total,
	}
}
