package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type AccountTxAdder interface {
	AddAccountTx(tx pgx.Tx, account Account, ctx context.Context) error
}

type OutboxTxInsert interface {
	InsertTx(tx pgx.Tx, msg OutboxMessage, ctx context.Context) error
}
