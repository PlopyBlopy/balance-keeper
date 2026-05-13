package domain

import (
	"errors"
	"math"
)

const MaxBalance float64 = 1_000_000_000.00

const (
	Exceeds             = ""
	NonNegative         = ""
	OnlyTwoCharAfterDot = ""
)

var (
	ErrBalanceExceeds             = errors.New(Exceeds)
	ErrBalanceNonNegative         = errors.New(NonNegative)
	ErrBalanceOnlyTwoCharAfterDot = errors.New(OnlyTwoCharAfterDot)
)

type Balance = int64

func NewBalance() Balance {
	return 0
}

// Balance => amount
func BalanceParseInt64(v Balance) float64 {
	return float64(v) / 100.00
}

// amount => Balance
func BalanceParseFloat64(v float64) (Balance, error) {
	if err := ValidateBalanceAmount(v); err != nil {
		return 0, err
	}

	return Balance(math.Round(v * 100)), nil
}

func ValidateBalanceAmount(v float64) error {
	if v > MaxBalance {
		return ErrBalanceExceeds
	}

	if v < 0 {
		return ErrBalanceNonNegative
	}

	scaled := v * 100

	rounded := math.Round(scaled)

	if math.Abs(scaled-rounded) > 1e-9 {
		return ErrBalanceOnlyTwoCharAfterDot
	}

	return nil
}

///
///
///

// func ParseDecimal(s string) (Balance, error) {
// 	//...
// }
// func ParseDecimalToInt(s string) (int64, error) {
// 	//...
// }

// func FromInt64(v int64) Balance {
// 	//...
// }

// // driver.Value implementation?
// func (b *Balance) Value() (driver.Value, error) {
// 	//...
// }

// // sql.Scanner implementation
// func (b *Balance) Scan(src any) error {
// 	//...
// }

// свой валидарор в playground validation
