package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/keremeti/iq-progers/config"
	"github.com/keremeti/iq-progers/internal/entity"
	"github.com/keremeti/iq-progers/internal/service"
	"github.com/keremeti/iq-progers/mocks"
	"github.com/keremeti/iq-progers/pkg/logger"
)

func TestTransferExecute_CorrectData_CorrectTransfer(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	mockInsertRemittance := entity.Remittance{
		SenderId:    1,
		RecipientId: 2,
		Sum:         22,
	}

	mockInsertSenderTra := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.WriteOff,
		UserId:  mockInsertRemittance.SenderId,
		Comment: "Перевод",
		Sum:     mockInsertRemittance.Sum,
	}
	mockInsertRecipientTra := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.Recharge,
		UserId:  mockInsertRemittance.RecipientId,
		Comment: "Перевод",
		Sum:     mockInsertRemittance.Sum,
	}
	var mockSendlerBalance float64 = 1000

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)
	mockTransactionRepo.EXPECT().Transfer(ctx, mockInsertSenderTra, mockInsertRecipientTra).Return(nil).Times(1)
	mockBalanceRepo := mocks.NewMockIBalanceRepo(ctrl)
	mockBalanceRepo.EXPECT().GetByUserId(ctx, mockInsertRemittance.SenderId).Return(mockSendlerBalance, nil).Times(1)

	tms := service.NewTransferMoneyService(l, mockTransactionRepo, mockBalanceRepo)
	err := tms.Execute(ctx, mockInsertRemittance)
	if err != nil {
		t.Fatalf("TestTransferExecute_CorrectData_CorrectTransfer - tms.Execute - error: %s", err)
	}
}

func TestTransferExecute_DbError_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	mockInsertRemittance := entity.Remittance{
		SenderId:    1,
		RecipientId: 2,
		Sum:         22,
	}

	mockInsertSenderTra := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.WriteOff,
		UserId:  mockInsertRemittance.SenderId,
		Comment: "Перевод",
		Sum:     mockInsertRemittance.Sum,
	}
	mockInsertRecipientTra := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.Recharge,
		UserId:  mockInsertRemittance.RecipientId,
		Comment: "Перевод",
		Sum:     mockInsertRemittance.Sum,
	}
	mockErr := fmt.Errorf("Some error")
	var mockSendlerBalance float64 = 1000

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)
	mockTransactionRepo.EXPECT().Transfer(ctx, mockInsertSenderTra, mockInsertRecipientTra).Return(mockErr).Times(1)
	mockBalanceRepo := mocks.NewMockIBalanceRepo(ctrl)
	mockBalanceRepo.EXPECT().GetByUserId(ctx, mockInsertRemittance.SenderId).Return(mockSendlerBalance, nil).Times(1)

	tms := service.NewTransferMoneyService(l, mockTransactionRepo, mockBalanceRepo)
	err := tms.Execute(ctx, mockInsertRemittance)
	if !(err != nil && err == fmt.Errorf("не удалось совершить перевод")) {

	} else {
		t.Fatalf("TestTransferExecute_DbError_Error - == - error: %s", err)
	}
}

func TestTransferExecute_NegativeAmount_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	mockInsertRemittance := entity.Remittance{
		SenderId:    1,
		RecipientId: 2,
		Sum:         -22,
	}

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)
	mockBalanceRepo := mocks.NewMockIBalanceRepo(ctrl)

	tms := service.NewTransferMoneyService(l, mockTransactionRepo, mockBalanceRepo)
	err := tms.Execute(ctx, mockInsertRemittance)
	if !(err != nil && err == fmt.Errorf("нельзя перевести сумму меньше или равную 0 рублей")) {

	} else {
		t.Fatalf("TestTransferExecute_NegativeAmount_Error - == - error: %s", err)
	}
}

func TestTransferExecute_NonExistentSender_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	mockInsertRemittance := entity.Remittance{
		SenderId:    0,
		RecipientId: 2,
		Sum:         22,
	}
	mockErr := fmt.Errorf("Some error")
	var mockSendlerBalance float64 = 0

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)
	mockBalanceRepo := mocks.NewMockIBalanceRepo(ctrl)
	mockBalanceRepo.EXPECT().GetByUserId(ctx, mockInsertRemittance.SenderId).Return(mockSendlerBalance, mockErr).Times(1)

	tms := service.NewTransferMoneyService(l, mockTransactionRepo, mockBalanceRepo)
	err := tms.Execute(ctx, mockInsertRemittance)
	if !(err != nil && err == fmt.Errorf("не удалось найти отправителя")) {

	} else {
		t.Fatalf("TestTransferExecute_NonExistentSender_Error - == - error: %s", err)
	}
}

func TestTransferExecute_NonExistentRecipient_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	mockInsertRemittance := entity.Remittance{
		SenderId:    1,
		RecipientId: 0,
		Sum:         22,
	}
	mockInsertSenderTra := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.WriteOff,
		UserId:  mockInsertRemittance.SenderId,
		Comment: "Перевод",
		Sum:     mockInsertRemittance.Sum,
	}
	mockInsertRecipientTra := entity.Transaction{
		Date:    time.Now().Unix(),
		Type:    entity.Recharge,
		UserId:  mockInsertRemittance.RecipientId,
		Comment: "Перевод",
		Sum:     mockInsertRemittance.Sum,
	}
	var mockSendlerBalance float64 = 1000.00
	mockErr := fmt.Errorf("Some error")

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)
	mockTransactionRepo.EXPECT().Transfer(ctx, mockInsertSenderTra, mockInsertRecipientTra).Return(mockErr).Times(1)
	mockBalanceRepo := mocks.NewMockIBalanceRepo(ctrl)
	mockBalanceRepo.EXPECT().GetByUserId(ctx, mockInsertRemittance.SenderId).Return(mockSendlerBalance, nil).Times(1)

	tms := service.NewTransferMoneyService(l, mockTransactionRepo, mockBalanceRepo)
	err := tms.Execute(ctx, mockInsertRemittance)
	if !(err != nil && err == fmt.Errorf("не удалось совершить перевод")) {

	} else {
		t.Fatalf("TestTransferExecute_NonExistentRecipient_Error - == - error: %s", err)
	}
}

func TestTransferExecute_NotEnoughMoney_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	mockInsertRemittance := entity.Remittance{
		SenderId:    1,
		RecipientId: 2,
		Sum:         22,
	}
	var mockSendlerBalance float64 = 1

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)
	mockBalanceRepo := mocks.NewMockIBalanceRepo(ctrl)
	mockBalanceRepo.EXPECT().GetByUserId(ctx, mockInsertRemittance.SenderId).Return(mockSendlerBalance, nil).Times(1)

	tms := service.NewTransferMoneyService(l, mockTransactionRepo, mockBalanceRepo)
	err := tms.Execute(ctx, mockInsertRemittance)
	if !(err != nil && err == fmt.Errorf("недостаточно средств на счете отправителя")) {

	} else {
		t.Fatalf("TestTransferExecute_NotEnoughMoney_Error - == - error: %s", err)
	}
}
