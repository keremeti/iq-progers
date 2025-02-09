package model

import "github.com/keremeti/iq-progers/internal/entity"

type TransactionModel struct {
	Id      int32                  `db:"id"`
	Date    int64                  `db:"date"`
	Type    entity.TransactionType `db:"type"`
	UserId  int32                  `db:"user_id"`
	Comment string                 `db:"comment"`
	Sum     float64                `db:"sum"`
}

func (m *TransactionModel) ToDomain() entity.Transaction {
	return entity.Transaction{
		Id:      m.Id,
		Date:    m.Date,
		Type:    m.Type,
		UserId:  m.UserId,
		Comment: m.Comment,
		Sum:     m.Sum,
	}
}

func NewTransactionModel(e entity.Transaction) TransactionModel {
	return TransactionModel{
		Id:      e.Id,
		Date:    e.Date,
		Type:    e.Type,
		UserId:  e.UserId,
		Comment: e.Comment,
		Sum:     e.Sum,
	}
}
