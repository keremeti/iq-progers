package repo

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/keremeti/iq-progers/internal/entity"
	"github.com/keremeti/iq-progers/internal/service/repo/model"
	"github.com/keremeti/iq-progers/pkg/postgres"
)

type TransactionsRepo struct {
	*postgres.Postgres
}

func NewTransactionsRepo(pg *postgres.Postgres) *TransactionsRepo {
	return &TransactionsRepo{pg}
}
func (r *TransactionsRepo) GetByFilter(ctx context.Context, filt entity.Filter,
	limit, page int) (entity.Pagination[entity.Transaction], error) {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return entity.Pagination[entity.Transaction]{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	offset := limit * (page - 1)
	query := `SELECT id, date, type, user_id, comment, sum FROM transactions WHERE user_id = @user_id
		ORDER BY date DESC LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"user_id": filt.UserId,
		"limit":   limit,
		"offset":  offset,
	}
	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return entity.Pagination[entity.Transaction]{}, fmt.Errorf("SongsRepo - GetByFilter - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := []entity.Transaction{}

	for rows.Next() {
		m := model.TransactionModel{}

		err = rows.Scan(&m.Id, &m.Date, &m.Type, &m.UserId, &m.Comment, &m.Sum)
		if err != nil {
			return entity.Pagination[entity.Transaction]{}, fmt.Errorf("SongsRepo - GetByFilter - model.SongModel: %w", err)
		}
		entities = append(entities, m.ToDomain())
	}

	query = `SELECT COUNT(*) FROM transactions WHERE user_id = @user_id`
	row := tx.QueryRow(ctx, query, args)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return entity.Pagination[entity.Transaction]{}, fmt.Errorf("SongsRepo - GetByFilter - count: %w", err)
	}
	p := entity.NewPagination(entities, limit, page, count)
	return p, nil
}
