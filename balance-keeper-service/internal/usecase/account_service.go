package usecase

import (
	"context"

	"github.com/PlopyBlopy/balance-keeper-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AddAccountFunc func(uuid.UUID, context.Context) (uuid.UUID, error)

func AddAccount(tm domain.TxManager, accountRep domain.AccountTxAdder, outboxRep domain.OutboxTxInsert) AddAccountFunc {
	return func(id uuid.UUID, ctx context.Context) (uuid.UUID, error) {
		var newId uuid.UUID

		err := tm.WithinTransaction(ctx, func(tx pgx.Tx) error {
			account := domain.NewAccount(id)

			err := accountRep.AddAccountTx(tx, account, ctx)
			if err != nil {
				return err
			}

			msg, err := domain.NewAccountCreatedEvent(
				domain.AccountCreatedEvent{
					Id:             account.Id,
					InitialBalance: domain.BalanceParseInt64(account.Balance),
					CreatedAt:      account.CreatedAt,
				},
			)
			if err != nil {
				return err
			}

			err = outboxRep.InsertTx(tx, msg, ctx)
			if err != nil {
				return err
			}

			err = tx.Commit(ctx)
			if err != nil {
				return err
			}

			newId = id

			return nil
		})

		return newId, err
	}
}

// func AddNewAccount() func(context.Context) (uuid.UUID, error) {
// 	return func(ctx context.Context) (uuid.UUID, error) {

// 	}
// }

// func GetAccount() func(uuid.UUID, context.Context) (domain.Account, error) {
// 	return func(id uuid.UUID, ctx context.Context) (domain.Account, error) {

// 	}
// }

// func GetAccounts() func(int, context.Context) ([]domain.Account, error) {
// 	return func(limit int, ctx context.Context) ([]domain.Account, error) {

// 	}
// }

// func GetBalance() func(uuid.UUID, context.Context) (domain.Balance, error) {
// 	return func(id uuid.UUID, ctx context.Context) (domain.Balance, error) {

// 	}
// }

// func Deposit() func(float64, context.Context) error {
// 	return func(amount float64, ctx context.Context) error {

// 	}
// }

// func Withdraw() func(float64, context.Context) error {
// 	return func(amount float64, ctx context.Context) error {

// 	}
// }

// func Transfer() func(uuid.UUID, uuid.UUID, float64, context.Context) error {
// 	return func(fromId, toId uuid.UUID, amount float64, ctx context.Context) error {

// 	}
// }
