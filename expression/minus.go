package expression

import (
	"currency_lib/rate"
	"github.com/shopspring/decimal"
)

type MinusExpr struct {
	minuend    Expression
	subtrahend Expression
}

func NewMinusExpr(minuend, subtrahend Expression) *MinusExpr {
	return &MinusExpr{
		minuend:    minuend,
		subtrahend: subtrahend,
	}
}

func (s *MinusExpr) Plus(expr Expression) Expression {
	return NewPlusExpr(s, expr)
}

func (s *MinusExpr) Minus(expr Expression) Expression {
	return NewMinusExpr(s, expr)
}

func (s *MinusExpr) Multiply(times decimal.Decimal) Expression {
	return NewMultiplyExpr(s, times)
}

func (s *MinusExpr) Divide(times decimal.Decimal) Expression {
	return NewDivideExpr(s, times)
}

func (m *MinusExpr) Reduce(to string, book *rate.Book) (*decimal.Decimal, bool) {
	minuend, valid := m.minuend.Reduce(to, book)
	if !valid {
		return nil, false
	}
	subtrahend, valid := m.subtrahend.Reduce(to, book)
	if !valid {
		return nil, false
	}
	result := minuend.Add((*subtrahend).Neg())
	return &result, true
}

func (m *MinusExpr) Optimize() Expression {
	return m
}

func (m *MinusExpr) GetName() string {
	return "minus_expr"
}
