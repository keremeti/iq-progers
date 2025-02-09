package service

import (
	"context"

	"github.com/keremeti/iq-progers/internal/entity"
)

type (
	GenericTypePaginationTransaction = entity.Pagination[entity.Transaction]

	ITopUpBalanceService interface {
		Execute(ctx context.Context, recharge entity.RechargeTransaction) (entity.Transaction, error)
	}
	ITransferMoneyService interface {
		Execute(ctx context.Context, remittance entity.Remittance) error
	}
	IGetTransactionsService interface {
		Execute(ctx context.Context, filt entity.Filter,
			limit, page int) (GenericTypePaginationTransaction, error)
	}

	ITransactionsRepo interface {
		GetByFilter(ctx context.Context, filt entity.Filter,
			limit, page int) (GenericTypePaginationTransaction, error)
	}

	ITransactionRepo interface {
		Insert(ctx context.Context, transaction entity.Transaction) (int32, error)
		Transfer(ctx context.Context, s, r entity.Transaction) error
	}

	IBalanceRepo interface {
		GetByUserId(ctx context.Context, userId int32) (float64, error)
	}
)
