package v1

import (
	"strconv"

	"github.com/keremeti/iq-progers/internal/entity"
)

type transactionResponse struct {
	Id      int32  `json:"id" binding:"required"  example:"12"`
	Date    int64  `json:"date" binding:"required" example:"21212"`
	Type    string `json:"type" binding:"required" example:"RECHARGE"`
	UserId  int32  `json:"user_id" binding:"required" example:"2"`
	Comment string `json:"comment" binding:"required" example:"Пополнение баланса"`
	Sum     string `json:"sum" binding:"required" example:"100"`
}

func newTransactionResponse(e entity.Transaction) transactionResponse {
	return transactionResponse{
		Id:      e.Id,
		Date:    e.Date,
		Type:    e.Type.ToString(),
		UserId:  e.UserId,
		Comment: e.Comment,
		Sum:     strconv.FormatFloat(e.Sum, 'f', 2, 64),
	}
}
