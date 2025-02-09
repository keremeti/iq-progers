package service

import (
	"context"
	"fmt"
	"time"

	"github.com/keremeti/iq-progers/internal/entity"
	"github.com/keremeti/iq-progers/pkg/logger"
)

type TransferMoneyService struct {
	transactionRepo ITransactionRepo
	balanceRepo     IBalanceRepo
	l               *logger.Logger
}

func NewTransferMoneyService(l *logger.Logger,
	tr ITransactionRepo, br IBalanceRepo) *TransferMoneyService {
	return &TransferMoneyService{
		transactionRepo: tr,
		balanceRepo:     br,
		l:               l,
	}
}

func (s *TransferMoneyService) Execute(ctx context.Context, remittance entity.Remittance) error {
	s.l.Info("TransferMoneyService - Execute")
	s.l.Debug("TransferMoneyService - Execute", "remittance", remittance)

	if remittance.Sum < 0 {
		s.l.Debug("TopUpBalanceService - Execute - transaction.Sum >= 0")
		return fmt.Errorf("нельзя перевести сумму меньше или равную 0 рублей")
	}

	senderBal, err := s.balanceRepo.GetByUserId(ctx, remittance.SenderId)
	if err != nil {
		s.l.Warn("TransferMoneyService - Execute - balanceRepo.GetByUserId: ", "err", err.Error())
		return fmt.Errorf("не удалось найти отправителя")
	}
	if senderBal-remittance.Sum < 0 {
		s.l.Info("TransferMoneyService - Execute - senderBal - remittance.Sum: ", "sender_balance", remittance.Sum)
		return fmt.Errorf("недостаточно средств на счете отправителя")
	}

	senderTra := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.WriteOff,
		UserId:  remittance.SenderId,
		Comment: "Перевод",
		Sum:     remittance.Sum,
	}
	recipientTra := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.Recharge,
		UserId:  remittance.RecipientId,
		Comment: "Перевод",
		Sum:     remittance.Sum,
	}

	err = s.transactionRepo.Transfer(ctx, senderTra, recipientTra)
	if err != nil {
		s.l.Warn("TransferMoneyService - Execute - transactionRepo.Transfer", "err", err.Error())
		return fmt.Errorf("не удалось совершить перевод")
	}
	return nil
}
