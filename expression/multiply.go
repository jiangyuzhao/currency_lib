package expression

import (
	"currency_lib/rate"
	"fmt"
	"github.com/shopspring/decimal"
	"runtime/debug"
)

type MultiplyExpr struct {
	multiplicand Expression
	multiplier   Expression
}

func NewMultiplyExpr(multiplicand Expression, multiplier decimal.Decimal) *MultiplyExpr {
	return &MultiplyExpr{
		multiplicand: multiplicand,
		multiplier:   GetExpressionFromDecimal(multiplier),
	}
}

func (s *MultiplyExpr) Plus(expr Expression) Expression {
	return NewPlusExpr(s, expr)
}

func (s *MultiplyExpr) Minus(expr Expression) Expression {
	return NewMinusExpr(s, expr)
}

func (s *MultiplyExpr) Multiply(times decimal.Decimal) Expression {
	return NewMultiplyExpr(s, times)
}

func (s *MultiplyExpr) Divide(times decimal.Decimal) Expression {
	return NewDivideExpr(s, times)
}

func (m *MultiplyExpr) Reduce(to string, book *rate.Book) (*decimal.Decimal, bool) {
	multiplicand, valid := m.multiplicand.Reduce(to, book)
	if !valid {
		return nil, false
	}
	multiplier, valid := m.multiplier.Reduce(to, book)
	if !valid {
		return nil, false
	}
	result := multiplicand.Mul(*multiplier)
	return &result, true
}

// 比如在目前的实现中, multiplier一定是常数, 如果被乘数是常数, 或者Money就可以直接修改数值
// (缺点是Money本可以自己一个包, 纯粹管理和金钱相关的数据, 后期可以和Bank等金钱流转系统类放在一起(也不算很严重, 因为Money是基础类, 金钱流转系统类依赖Money即可), 由于不能循环依赖, 因此Money会被放入表达式这个包)
// 若第一个操作数是加法表达式或者减法表达式, 就可以用分配律, 构造: 第一个操作数' = 第一个操作数 * NumberExpr, 第二个操作数' = 第二个操作数 * NumberExpr
// 递归调用第一个操作数'和第二个操作数'的Optimize方法, 之后返回它们的和或者差
// 若第一个操作数是乘法表达式, 则可以跟前一个乘法表达式的操作数合并
// 若第一个操作数是除法表达式, 则可以进行乘法前置, 把乘法前置到除法之前, 总之原则是把除法的计算尽可能滞后(汇率转换涉及到除法, 因此汇率转换也应该尽可能后置)
func (m *MultiplyExpr) Optimize() (ret Expression) {
	times, valid := m.multiplier.(*NumberExpr)
	if !valid {
		fmt.Printf("ERROR! multiply's times is not pointer of NumberExpr, but %T\n", m.multiplier)
		debug.PrintStack()
		return m
	}
	switch firstExpr := m.multiplicand.(type) {
	case *Money:
		ret = NewMoney(firstExpr.Amount.Mul(times.num), firstExpr.CurrencyName)
	case *NumberExpr:
		ret = GetExpressionFromDecimal(firstExpr.num.Mul(times.num))
	case *PlusExpr:
		firstExprTmp := firstExpr.augend.Multiply(times.num)
		secondExprTmp := firstExpr.addend.Multiply(times.num)
		ret = firstExprTmp.Optimize().Plus(secondExprTmp.Optimize())
	case *MinusExpr:
		firstExprTmp := firstExpr.minuend.Multiply(times.num)
		secondExprTmp := firstExpr.subtrahend.Multiply(times.num)
		ret = firstExprTmp.Optimize().Minus(secondExprTmp.Optimize())
	case *MultiplyExpr:
		firstTimes, valid := firstExpr.multiplier.(*NumberExpr)
		if !valid {
			fmt.Printf("ERROR! multiply's first expr is multiply but its times is not pointer of NumberExpr, but %T\n", firstExpr.multiplier)
			debug.PrintStack()
			return m
		}
		ret = firstExpr.multiplicand.Multiply(firstTimes.num.Mul(times.num)).Optimize()
	case *DivideExpr:
		divisorNum, valid := firstExpr.divisor.(*NumberExpr)
		if !valid {
			fmt.Printf("ERROR! multiply's first expr is DivideExpr but its times is not pointer of NumberExpr, but %T\n", firstExpr.divisor)
			debug.PrintStack()
			return m
		}
		ret = firstExpr.dividend.Multiply(times.num).Optimize().Divide(divisorNum.num)
	default:
		ret = m
	}
	return ret
}

func (m *MultiplyExpr) GetName() string {
	return "multiply_expr"
}
