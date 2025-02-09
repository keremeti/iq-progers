package service

import (
	"context"
	"fmt"
	"time"

	"github.com/keremeti/iq-progers/internal/entity"
	"github.com/keremeti/iq-progers/pkg/logger"
)

type TopUpBalanceService struct {
	transactionRepo ITransactionRepo
	l               *logger.Logger
}

func NewTopUpBalanceService(l *logger.Logger, tr ITransactionRepo) *TopUpBalanceService {
	return &TopUpBalanceService{
		transactionRepo: tr,
		l:               l,
	}
}

func (s *TopUpBalanceService) Execute(ctx context.Context, recharge entity.RechargeTransaction) (entity.Transaction, error) {
	s.l.Info("TopUpBalanceService - Execute")
	s.l.Debug("TopUpBalanceService - Execute", "recharge", recharge)
	transaction := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.Recharge,
		UserId:  recharge.UserId,
		Comment: "Пополнение",
		Sum:     recharge.Sum,
	}
	if transaction.Sum < 0 {
		s.l.Debug("TopUpBalanceService - Execute - transaction.Sum >= 0")
		return entity.Transaction{}, fmt.Errorf("нельзя пополнить баланс на сумму меньше или равную 0 рублей")
	}
	id, err := s.transactionRepo.Insert(ctx, transaction)
	if err != nil {
		s.l.Warn("TopUpBalanceService - Execute - transactionRepo.Insert: ", "err", err.Error())
		return entity.Transaction{}, fmt.Errorf("не удалось пополнить баланс")
	}

	transaction.Id = id
	return transaction, nil
}
