package repo

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/keremeti/iq-progers/internal/entity"
	"github.com/keremeti/iq-progers/internal/service/repo/model"
	"github.com/keremeti/iq-progers/pkg/postgres"
)

type TransactionRepo struct {
	*postgres.Postgres
}

func NewTransactionRepo(pg *postgres.Postgres) *TransactionRepo {
	return &TransactionRepo{pg}
}

func (r *TransactionRepo) Insert(ctx context.Context, transaction entity.Transaction) (int32, error) {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	tm := model.NewTransactionModel(transaction)
	query := `UPDATE users SET balance = balance + @sum WHERE id = @user_id`
	args := pgx.NamedArgs{
		"user_id": tm.UserId,
		"sum":     tm.Sum,
	}

	tag, err := tx.Exec(ctx, query, args)
	if err != nil {
		return 0, fmt.Errorf("TransactionRepo - Insert - UPDATE: %w", err)
	}
	if tag.String() == "UPDATE 0" {
		return 0, fmt.Errorf("TransactionRepo - Insert - UPDATE: %s", tag.String())
	}

	query = `INSERT INTO transactions (date, type, user_id, comment, sum)
			VALUES (@date, @type, @user_id, @comment, @sum)
			RETURNING id`
	args = pgx.NamedArgs{
		"date":    tm.Date,
		"type":    tm.Type,
		"user_id": tm.UserId,
		"comment": tm.Comment,
		"sum":     tm.Sum,
	}

	row := tx.QueryRow(ctx, query, args)
	var id int32
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("TransactionRepo - Insert - INSERT: %w", err)
	}

	return id, nil
}

func (r *TransactionRepo) Transfer(ctx context.Context, sen, rec entity.Transaction) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	senTm := model.NewTransactionModel(sen)
	query := `UPDATE users SET balance = balance - @sum WHERE id = @user_id`
	args := pgx.NamedArgs{
		"user_id": senTm.UserId,
		"sum":     senTm.Sum,
	}
	tag, err := tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("TransactionRepo - Transfer - UPDATE(senTm): %w", err)
	}
	if tag.String() == "UPDATE 0" {
		return fmt.Errorf("TransactionRepo - Transfer - UPDATE(senTm): %s", tag.String())
	}

	recTm := model.NewTransactionModel(rec)
	query = `UPDATE users SET balance = balance + @sum WHERE id = @user_id`
	args = pgx.NamedArgs{
		"user_id": recTm.UserId,
		"sum":     recTm.Sum,
	}
	tag, err = tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("TransactionRepo - Transfer - UPDATE(recTm): %w", err)
	}
	if tag.String() == "UPDATE 0" {
		return fmt.Errorf("TransactionRepo - Transfer - UPDATE(recTm): %s", tag.String())
	}

	query = `INSERT INTO transactions (date, type, user_id, comment, sum)
			VALUES (@date, @type, @user_id, @comment, @sum)
			RETURNING id`
	args = pgx.NamedArgs{
		"date":    senTm.Date,
		"type":    senTm.Type,
		"user_id": senTm.UserId,
		"comment": senTm.Comment,
		"sum":     senTm.Sum,
	}

	tag, err = tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("TransactionRepo - Transfer - INSERT(senTm): %w", err)
	}
	if tag.String() == "INSERT 0" {
		return fmt.Errorf("TransactionRepo - Transfer - INSERT(senTm): %s", tag.String())
	}

	query = `INSERT INTO transactions (date, type, user_id, comment, sum)
			VALUES (@date, @type, @user_id, @comment, @sum)
			RETURNING id`
	args = pgx.NamedArgs{
		"date":    recTm.Date,
		"type":    recTm.Type,
		"user_id": recTm.UserId,
		"comment": recTm.Comment,
		"sum":     recTm.Sum,
	}

	tag, err = tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("TransactionRepo - Insert - INSERT(recTm): %w", err)
	}
	if tag.String() == "INSERT 0" {
		return fmt.Errorf("TransactionRepo - Insert - INSERT(recTm): %s", tag.String())
	}

	return nil
}
