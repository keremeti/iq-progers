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

func TestTopUpExecute_CorrectData_CorrectTransaction(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()

	recharge := entity.RechargeTransaction{
		UserId: 1,
		Sum:    22,
	}
	var mockId int32 = 1
	date := time.Now().Unix()
	mockInsertTransaction := entity.Transaction{
		Id:      0,
		Type:    entity.Recharge,
		Date:    date,
		UserId:  1,
		Comment: "Пополнение",
		Sum:     22,
	}

	mockTransaction := entity.Transaction{
		Id:      1,
		Type:    entity.Recharge,
		Date:    date,
		UserId:  1,
		Comment: "Пополнение",
		Sum:     22,
	}

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)
	mockTransactionRepo.EXPECT().Insert(ctx, mockInsertTransaction).Return(mockId, nil).Times(1)

	tubs := service.NewTopUpBalanceService(l, mockTransactionRepo)
	transaction, err := tubs.Execute(ctx, recharge)
	if err != nil {
		t.Fatalf("TestTopUpExecute_CorrectData_SuccessRecharge - tubs.Execute - error: %s", err)
	}
	if !TransactionEquals(mockTransaction, transaction) {
		t.Fatalf("TestTopUpExecute_CorrectData_SuccessRecharge - !TransactionEquals - error: %s", "транзакция, созданная сервисом, не верна")
	}
}

func TestTopUpExecute_NegativeAmount_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()

	recharge := entity.RechargeTransaction{
		UserId: 1,
		Sum:    -22,
	}

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)

	tubs := service.NewTopUpBalanceService(l, mockTransactionRepo)
	_, err := tubs.Execute(ctx, recharge)
	if !(err != nil && err == fmt.Errorf("нельзя пополнить баланс на сумму меньше или равную 0 рублей")) {

	} else {
		t.Fatalf("TestTopUpExecute_NegativeAmount_Error - == - error: %s", err)
	}
}

func TestTopUpExecute_NonExistentUser_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()

	recharge := entity.RechargeTransaction{
		UserId: 1,
		Sum:    -22,
	}

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)

	tubs := service.NewTopUpBalanceService(l, mockTransactionRepo)
	_, err := tubs.Execute(ctx, recharge)
	if !(err != nil && err == fmt.Errorf("не удалось пополнить баланс")) {

	} else {
		t.Fatalf("TestTopUpExecute_NonExistentUser_Error - == - error: %s", err)
	}
}

func TestTopUpExecute_DbError_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()

	recharge := entity.RechargeTransaction{
		UserId: 1,
		Sum:    22,
	}
	var mockId int32 = 0
	date := time.Now().Unix()
	mockInsertTransaction := entity.Transaction{
		Id:      0,
		Type:    entity.Recharge,
		Date:    date,
		UserId:  1,
		Comment: "Пополнение",
		Sum:     22,
	}
	mockErr := fmt.Errorf("Some error")

	mockTransactionRepo := mocks.NewMockITransactionRepo(ctrl)
	mockTransactionRepo.EXPECT().Insert(ctx, mockInsertTransaction).Return(mockId, mockErr).Times(1)

	tubs := service.NewTopUpBalanceService(l, mockTransactionRepo)
	_, err := tubs.Execute(ctx, recharge)
	if !(err != nil && err == fmt.Errorf("не удалось пополнить баланс")) {

	} else {
		t.Fatalf("TestTopUpExecute_DbError_Error - == - error: %s", err)
	}
}

func TransactionEquals(first, next entity.Transaction) bool {
	id_eq := first.Id == next.Id
	type_eq := first.Type == next.Type
	user_eq := first.UserId == next.UserId
	comment_eq := first.Comment == next.Comment
	sum_eq := first.Sum == next.Sum
	return (id_eq && type_eq && user_eq && comment_eq && sum_eq)
}
