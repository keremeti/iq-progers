package v1

import (
	"strconv"

	"github.com/keremeti/iq-progers/internal/entity"
)

type remittanceRequest struct {
	SenderId    int32  `json:"sender_id" binding:"required"  example:"1"`
	RecipientId int32  `json:"recipient_id" binding:"required"  example:"2"`
	Sum         string `json:"sum" binding:"required"  example:"100"`
}

func (r *remittanceRequest) ToDomain() (entity.Remittance, error) {
	sum, err := strconv.ParseFloat(r.Sum, 64)
	if err != nil {
		return entity.Remittance{}, err
	}
	return entity.Remittance{
		SenderId:    r.SenderId,
		RecipientId: r.RecipientId,
		Sum:         sum,
	}, nil
}
