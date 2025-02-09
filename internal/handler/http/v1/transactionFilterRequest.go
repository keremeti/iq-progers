package v1

import "github.com/keremeti/iq-progers/internal/entity"

type filterRequest struct {
	UserId int32 `json:"user_id" binding:"required" example:"2"`
}

func (r *filterRequest) ToDomain() entity.Filter {
	return entity.Filter{
		UserId: r.UserId,
	}
}
