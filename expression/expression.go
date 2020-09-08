package expression

import (
	"currency_lib/rate"
	"github.com/shopspring/decimal"
)

type Expression interface {
	Plus(expr Expression) Expression
	Minus(expr Expression) Expression
	Multiply(times decimal.Decimal) Expression
	Divide(times decimal.Decimal) Expression
	Reduce(to string, book *rate.Book) (*decimal.Decimal, bool)
	Optimize() Expression // get compute plan and optimize, just empty now, should return a new Expression
	GetName() string      // for debug
}

type NumberExpr struct {
	num decimal.Decimal
}

func GetExpressionFromDecimal(num decimal.Decimal) *NumberExpr {
	return &NumberExpr{
		num: num,
	}
}

func (n *NumberExpr) Plus(expr Expression) Expression {
	return NewPlusExpr(n, expr)
}

func (n *NumberExpr) Minus(expr Expression) Expression {
	return NewMinusExpr(n, expr)
}

func (n *NumberExpr) Multiply(times decimal.Decimal) Expression {
	return NewMultiplyExpr(n, times)
}

func (n *NumberExpr) Divide(times decimal.Decimal) Expression {
	return NewDivideExpr(n, times)
}

func (n *NumberExpr) Reduce(to string, book *rate.Book) (*decimal.Decimal, bool) {
	return &n.num, true
}

func (n *NumberExpr) Optimize() Expression {
	return n
}

func (n *NumberExpr) GetName() string {
	return "number_expr"
}
