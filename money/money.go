package money

import (
	"github.com/oarkflow/pkg/decimal"
)

type Money struct {
	Amount   decimal.Decimal
	Currency string
}

func (m *Money) ToCurrency(currency string) decimal.Decimal {
	return Convert(m.Amount, m.Currency, currency)
}
