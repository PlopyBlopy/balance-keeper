package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Balance   Balance   `json:"balance" db:"balance"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func NewAccount(id uuid.UUID) Account {
	curTime := time.Now().UTC().Truncate(time.Microsecond)
	return Account{
		Id:        id,
		Balance:   NewBalance(),
		UpdatedAt: curTime,
		CreatedAt: curTime,
	}
}

func NewAccountWithBalance(id uuid.UUID, amount float64) (*Account, error) {

	balance, err := BalanceParseFloat64(amount)
	if err != nil {
		return nil, err
	}

	curTime := time.Now().UTC().Truncate(time.Microsecond)
	return &Account{
		Id:        id,
		Balance:   balance,
		UpdatedAt: curTime,
		CreatedAt: curTime,
	}, nil
}
