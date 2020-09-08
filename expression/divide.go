package expression

import (
	"currency_lib/rate"
	"fmt"
	"github.com/shopspring/decimal"
	"runtime/debug"
)

type DivideExpr struct {
	dividend Expression
	divisor  Expression
}

func NewDivideExpr(dividend Expression, divisor decimal.Decimal) *DivideExpr {
	return &DivideExpr{
		dividend: dividend,
		divisor:  GetExpressionFromDecimal(divisor),
	}
}

func (s *DivideExpr) Plus(expr Expression) Expression {
	return NewPlusExpr(s, expr)
}

func (s *DivideExpr) Minus(expr Expression) Expression {
	return NewMinusExpr(s, expr)
}

func (s *DivideExpr) Multiply(times decimal.Decimal) Expression {
	return NewMultiplyExpr(s, times)
}

func (s *DivideExpr) Divide(times decimal.Decimal) Expression {
	return NewDivideExpr(s, times)
}

func (d *DivideExpr) Reduce(to string, book *rate.Book) (*decimal.Decimal, bool) {
	dividend, valid := d.dividend.Reduce(to, book)
	if !valid {
		return nil, false
	}
	divisor, valid := d.divisor.Reduce(to, book)
	if !valid {
		return nil, false
	}
	result := dividend.Div(*divisor)
	return &result, true
}

// 目前对除法的优化只有两次除法的合并
func (d *DivideExpr) Optimize() (ret Expression) {
	switch firstExpr := d.dividend.(type) {
	case *DivideExpr:
		firstTimes, valid := firstExpr.divisor.(*NumberExpr)
		if !valid {
			fmt.Printf("ERROR! DivideExpr's dividend is DivideExpr is not pointer of NumberExpr, but %T\n", d.divisor)
			debug.PrintStack()
			return d
		}
		secondTimes, valid := d.divisor.(*NumberExpr)
		if !valid {
			fmt.Printf("ERROR! DivideExpr's divisor is not pointer of NumberExpr, but %T\n", d.divisor)
			debug.PrintStack()
			return d
		}
		ret = firstExpr.dividend.Optimize().Divide(firstTimes.num.Mul(secondTimes.num))
	default:
		ret = d
	}
	return
}

func (d *DivideExpr) GetName() string {
	return "divide_expr"
}
