package v1

import (
	"strconv"

	"github.com/keremeti/iq-progers/internal/entity"
)

type rechargeRequest struct {
	UserId int32  `json:"user_id" binding:"required"  example:"1"`
	Sum    string `json:"sum" binding:"required"  example:"100"`
}

func (r *rechargeRequest) ToDomain() (entity.RechargeTransaction, error) {
	sum, err := strconv.ParseFloat(r.Sum, 64)
	if err != nil {
		return entity.RechargeTransaction{}, err
	}
	return entity.RechargeTransaction{
		UserId: r.UserId,
		Sum:    sum,
	}, nil
}
