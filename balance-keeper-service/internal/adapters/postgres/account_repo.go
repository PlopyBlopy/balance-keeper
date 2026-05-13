package postgres

import (
	"context"
	"errors"

	"github.com/PlopyBlopy/balance-keeper-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type AccountRepository interface {
// 	AddAccountTx(tx pgx.Tx, account domain.Account, ctx context.Context)
// }

type AccountRepository struct {
	db
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		db: db{pool: pool},
	}
}

func (r *AccountRepository) AddAccountTx(tx pgx.Tx, account domain.Account, ctx context.Context) error {
	ct, err := tx.Exec(ctx, "INSERT INTO accounts (id, balance, updated_at, created_at) VALUES ($1, $2, $3, $4)", account.Id, account.Balance, account.UpdatedAt, account.CreatedAt)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotAdded
	}

	return nil
}

func (r *AccountRepository) GetAccount(id uuid.UUID, ctx context.Context) (domain.Account, error) {
	account := domain.Account{}

	err := r.pool.QueryRow(ctx, "SELECT * FROM accounts WHERE id = $1", id).Scan(&account.Id, &account.Balance, &account.UpdatedAt, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return account, domain.ErrNotFound
		}
		return account, err
	}

	return account, nil
}
