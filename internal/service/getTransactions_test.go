package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/keremeti/iq-progers/config"
	"github.com/keremeti/iq-progers/internal/entity"
	"github.com/keremeti/iq-progers/internal/service"
	"github.com/keremeti/iq-progers/mocks"
	"github.com/keremeti/iq-progers/pkg/logger"
)

func TestGetTsExecute_CorrectData_Success(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	testFilt := entity.Filter{UserId: 1}
	limit := 2
	page := 1
	total := 3
	mockTransactions := []entity.Transaction{
		{Id: 1,
			Date:    1738947039,
			Type:    entity.Recharge,
			UserId:  1,
			Comment: "Пополнение",
			Sum:     1},
		{Id: 2,
			Date:    1738947039,
			Type:    entity.WriteOff,
			UserId:  1,
			Comment: "Перевод",
			Sum:     1},
	}
	mockList := entity.NewPagination(mockTransactions, limit, page, total)

	mockTransactionsRepo := mocks.NewMockITransactionsRepo(ctrl)
	mockTransactionsRepo.EXPECT().GetByFilter(ctx, testFilt, limit, page).Return(mockList, nil).Times(1)

	gts := service.NewGetTransactionsService(l, mockTransactionsRepo)
	list, err := gts.Execute(ctx, testFilt, limit, page)
	if err != nil {
		t.Fatalf("TestGetTsExecute_CorrectData_Success - gts.Execute - error: %s", err)
	}
	if !PaginationEquals(mockList, list) {
		t.Fatalf("TestGetTsExecute_CorrectData_Success - PaginationEquals - error: %s",
			"список, полученый из репозитория, был изменен сервисом")
	}
}

func TestGetTsExecute_DbError_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	testFilt := entity.Filter{UserId: 1}
	limit := 2
	page := 1
	mockErr := fmt.Errorf("Some error")

	mockTransactionsRepo := mocks.NewMockITransactionsRepo(ctrl)
	mockTransactionsRepo.EXPECT().GetByFilter(ctx,
		testFilt, limit, page).Return(entity.Pagination[entity.Transaction]{}, mockErr).Times(1)

	gts := service.NewGetTransactionsService(l, mockTransactionsRepo)
	_, err := gts.Execute(ctx, testFilt, limit, page)
	if !(err != nil && err == fmt.Errorf("не удалось получить данные")) {

	} else {
		t.Fatalf("TestGetTsExecute_DbError_Error - == - error: %s", err)
	}
}

func TestGetTsExecute_NonExistentUser_Error(t *testing.T) {
	l := logger.New(config.Dev)

	ctrl, ctx := gomock.WithContext(context.Background(), t)
	defer ctrl.Finish()
	testFilt := entity.Filter{UserId: 0}
	limit := 2
	page := 1
	mockErr := fmt.Errorf("Some error")

	mockTransactionsRepo := mocks.NewMockITransactionsRepo(ctrl)
	mockTransactionsRepo.EXPECT().GetByFilter(ctx,
		testFilt, limit, page).Return(entity.Pagination[entity.Transaction]{}, mockErr).Times(1)

	gts := service.NewGetTransactionsService(l, mockTransactionsRepo)
	_, err := gts.Execute(ctx, testFilt, limit, page)
	if !(err != nil && err == fmt.Errorf("не удалось получить данные")) {

	} else {
		t.Fatalf("TestGetTsExecute_NonExistentUser_Error - == - error: %s", err)
	}
}

func PaginationEquals(first, next entity.Pagination[entity.Transaction]) bool {
	current_eq := first.CurrentPage == next.CurrentPage
	next_eq := first.HasNext == next.HasNext
	previous_eq := first.HasPrevious == next.HasPrevious
	per_eq := first.RecordPerPage == next.RecordPerPage
	tcount_eq := first.TotalCount == next.TotalCount
	tpage_eq := first.TotalPage == next.TotalPage

	slise_eq := len(first.Slice) == len(next.Slice)
	if slise_eq {
		for index, item := range first.Slice {
			if item != next.Slice[index] {
				slise_eq = false
			}
		}
	}
	return (current_eq && next_eq && previous_eq && per_eq && tcount_eq && tpage_eq && slise_eq)
}
