package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error
}
