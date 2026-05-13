package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceParseFloat64(t *testing.T) {
	// Arrange
	assert := assert.New(t)

	tests := []struct {
		name      string
		testValue float64
		wantValue Balance
		wantErr   error
	}{
		{
			name:      "zero value",
			testValue: 0,
			wantValue: 0,
			wantErr:   nil,
		},
		{
			name:      "nothing after dot characters",
			testValue: 100,
			wantValue: 10000,
			wantErr:   nil,
		},
		{
			name:      "one character after dot",
			testValue: 100.2,
			wantValue: 10020,
			wantErr:   nil,
		},
		{
			name:      "two character after dot",
			testValue: 100.22,
			wantValue: 10022,
			wantErr:   nil,
		},
		{
			name:      "three character after dot",
			testValue: 100.222,
			wantValue: 0,
			wantErr:   ErrBalanceOnlyTwoCharAfterDot,
		},
		{
			name:      "non-negative",
			testValue: -1,
			wantValue: 0,
			wantErr:   ErrBalanceNonNegative,
		},
		{
			name:      "Exceeds",
			testValue: 1_000_000_000.2,
			wantValue: 0,
			wantErr:   ErrBalanceExceeds,
		},
	}

	// Act
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotError := BalanceParseFloat64(tt.testValue)
			fl := BalanceParseInt64(gotValue)
			_ = fl
			// Assert
			assert.ErrorIs(gotError, tt.wantErr)
			assert.Equal(tt.wantValue, gotValue)
		})
	}
}
