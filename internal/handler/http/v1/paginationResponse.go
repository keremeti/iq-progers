package v1

import (
	"github.com/keremeti/iq-progers/internal/entity"
)

type paginationResponse struct {
	HasNext       bool `json:"HasNext"`
	HasPrevious   bool `json:"HasPrevious"`
	RecordPerPage int  `json:"PageSize"`
	CurrentPage   int  `json:"CurrentPage"`
	TotalPage     int  `json:"TotalPages"`
	TotalCount    int  `json:"TotalCount"`
}

func NewPaginationResponse[T any](p *entity.Pagination[T]) paginationResponse {
	return paginationResponse{
		HasNext:       p.HasNext,
		HasPrevious:   p.HasPrevious,
		RecordPerPage: p.RecordPerPage,
		CurrentPage:   p.CurrentPage,
		TotalPage:     p.TotalPage,
		TotalCount:    p.TotalCount,
	}
}
