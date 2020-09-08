package expression

import (
	"currency_lib/rate"
	"github.com/shopspring/decimal"
)

type PlusExpr struct {
	augend Expression
	addend Expression
}

func NewPlusExpr(augend, addend Expression) *PlusExpr {
	return &PlusExpr{
		augend: augend,
		addend: addend,
	}
}

func (s *PlusExpr) Plus(expr Expression) Expression {
	return NewPlusExpr(s, expr)
}

func (s *PlusExpr) Minus(expr Expression) Expression {
	return NewMinusExpr(s, expr)
}

func (s *PlusExpr) Multiply(times decimal.Decimal) Expression {
	return NewMultiplyExpr(s, times)
}

func (s *PlusExpr) Divide(times decimal.Decimal) Expression {
	return NewDivideExpr(s, times)
}

func (s *PlusExpr) Reduce(to string, book *rate.Book) (*decimal.Decimal, bool) {
	augend, valid := s.augend.Reduce(to, book)
	if !valid {
		return nil, false
	}
	addend, valid := s.addend.Reduce(to, book)
	if !valid {
		return nil, false
	}
	result := augend.Add(*addend)
	return &result, true
}

func (s *PlusExpr) Optimize() Expression {
	return s
}

func (s *PlusExpr) GetName() string {
	return "plus_expr"
}
