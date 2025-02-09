package repo

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/keremeti/iq-progers/pkg/postgres"
)

type BalanceRepo struct {
	*postgres.Postgres
}

func NewBalanceRepo(pg *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{pg}
}

func (r *BalanceRepo) GetByUserId(ctx context.Context, userId int32) (float64, error) {
	query := `SELECT balance FROM users WHERE id=@user_id`
	args := pgx.NamedArgs{
		"user_id": userId,
	}
	row := r.Pool.QueryRow(ctx, query, args)

	var balance float64
	err := row.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("BalanceRepo - GetByUserId - row.Scan: %w", err)
	}

	return balance, nil
}
