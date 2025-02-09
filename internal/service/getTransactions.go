package service

import (
	"context"
	"fmt"

	"github.com/keremeti/iq-progers/internal/entity"
	"github.com/keremeti/iq-progers/pkg/logger"
)

type GetTransactionsService struct {
	transactionsRepo ITransactionsRepo
	l                *logger.Logger
}

func NewGetTransactionsService(l *logger.Logger, tr ITransactionsRepo) *GetTransactionsService {
	return &GetTransactionsService{
		transactionsRepo: tr,
		l:                l,
	}
}

func (s *GetTransactionsService) Execute(ctx context.Context, filt entity.Filter,
	limit, page int) (entity.Pagination[entity.Transaction], error) {
	s.l.Info("GetTransactionsService - Execute")
	s.l.Debug("GetTransactionsService - Execute", "filt", filt, "limit", limit, "page", page)
	songs, err := s.transactionsRepo.GetByFilter(ctx, filt, limit, page)
	if err != nil {
		s.l.Warn("GetTransactionsService - Execute - transactionsRepo.GetByFilter: ", "err", err.Error())
		return entity.Pagination[entity.Transaction]{}, fmt.Errorf("не удалось получить данные")
	}
	return songs, nil
}
