package expression

import (
	"currency_lib/rate"
	"fmt"
	"github.com/shopspring/decimal"
)

type Money struct {
	Amount       decimal.Decimal
	CurrencyName string
}

func NewMoney(amount decimal.Decimal, currencyName string) *Money {
	return &Money{
		Amount:       amount,
		CurrencyName: currencyName,
	}
}

func (m *Money) GetCurrencyName() string {
	return m.CurrencyName
}

func (m *Money) Plus(money Expression) Expression {
	if realMoney, valid := money.(*Money); valid {
		if realMoney.CurrencyName == m.CurrencyName {
			return NewMoney(m.Amount.Add(realMoney.Amount), m.CurrencyName)
		}
	}
	return NewPlusExpr(m, money)
}

func (m *Money) Minus(money Expression) Expression {
	if realMoney, valid := money.(*Money); valid {
		if realMoney.CurrencyName == m.CurrencyName {
			return NewMoney(m.Amount.Add(realMoney.Amount.Neg()), m.CurrencyName)
		}
	}
	return NewMinusExpr(m, money)
}

func (m *Money) Multiply(times decimal.Decimal) Expression {
	return NewMoney(m.Amount.Mul(times), m.CurrencyName)
}

func (m *Money) Divide(times decimal.Decimal) Expression {
	// divide计算需要延迟
	return NewDivideExpr(m, times)
}

func (m *Money) Reduce(to string, book *rate.Book) (*decimal.Decimal, bool) {
	r, err := book.GetRate(m.CurrencyName, to)
	if err != nil || r == nil {
		fmt.Printf("get rate error, from %s to %s\n", m.CurrencyName, to)
		return nil, false
	}
	result := m.Amount.Div(*r)
	return &result, true
}

func (m *Money) Optimize() Expression {
	return m
}

func (m *Money) GetName() string {
	return "money"
}
