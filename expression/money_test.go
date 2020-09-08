package expression

import (
	"currency_lib/rate"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
)

func TestMoneyNormal(t *testing.T) {
	book := new(rate.Book)
	book.Add("rmb", "usd", decimal.NewFromFloat(2.1))
	five := NewMoney(decimal.NewFromInt(5), "rmb")
	result := five.Plus(NewMoney(decimal.NewFromInt(12), "rmb"))
	d, v := result.Reduce("USD", book)
	expect := decimal.NewFromInt(17).Div(decimal.NewFromFloat(2.1))
	if !v || !reflect.DeepEqual(*d, expect) {
		t.Errorf("expect: %+v, got: %+v", expect, d)
	}
	result = result.Minus(NewMoney(decimal.NewFromInt(6), "usd"))
	d, v = result.Reduce("USD", book)
	expect = decimal.NewFromInt(17).Div(decimal.NewFromFloat(2.1)).Add(decimal.NewFromInt(6).Neg())
	if !v || !reflect.DeepEqual(*d, expect) {
		t.Errorf("expect: %+v, got: %+v", expect, d)
	}
	result = result.Multiply(decimal.NewFromFloat(2.2))
	d, v = result.Reduce("USD", book)
	expect = decimal.NewFromFloat(17).Div(decimal.NewFromFloat(2.1)).Add(decimal.NewFromFloat(6).Neg()).Mul(decimal.NewFromFloat(2.2))
	if !v || !reflect.DeepEqual(*d, expect) {
		t.Errorf("expect: %+v, got: %+v", expect, d)
	}
}

func TestMultiplyOptimize(t *testing.T) {
	book := new(rate.Book)
	book.Add("rmb", "usd", decimal.NewFromFloat(2.1))
	result := NewMoney(decimal.NewFromInt(5), "rmb").Plus(NewMoney(decimal.NewFromInt(12), "rmb")).Minus(NewMoney(decimal.NewFromInt(6), "usd")).Multiply(decimal.NewFromFloat(2.2))
	d, v := result.Optimize().Reduce("USD", book)
	expect := decimal.NewFromFloat(17 * 2.2).Div(decimal.NewFromFloat(2.1)).Add(decimal.NewFromFloat(6 * 2.2).Neg())
	if !v || !reflect.DeepEqual(*d, expect) {
		t.Errorf("expect: %+v, got: %+v", expect, d)
	}
}

func TestDivideExprOptimize(t *testing.T) {
	book := new(rate.Book)
	result := NewMoney(decimal.NewFromInt(5), "rmb").Divide(decimal.NewFromFloat(7)).Divide(decimal.NewFromFloat(1.0 / 7.0))
	d, v := result.Optimize().Reduce("rmb", book)
	expect := decimal.NewFromFloat(5).Div(decimal.NewFromFloat(7).Mul(decimal.NewFromFloat(1.0 / 7.0)))
	if !reflect.DeepEqual(expect, decimal.NewFromFloat(5).Div(decimal.NewFromFloat(1))) {
		t.Logf("precision problem for decimal, expect: 5, got: %+v\n", expect)
	}
	if !v || !reflect.DeepEqual(*d, expect) {
		t.Errorf("expect: %+v, got: %+v", expect, d)
	}
}

func TestDivideMultiNotOptimize(t *testing.T) {
	book := new(rate.Book)
	book.Add("rmb", "usd", decimal.NewFromFloat(2.3))
	result := NewMoney(decimal.NewFromInt(5), "rmb").Plus(NewMoney(decimal.NewFromInt(12), "rmb")).Minus(NewMoney(decimal.NewFromInt(6), "usd")).Divide(decimal.NewFromFloat(7)).Multiply(decimal.NewFromFloat(7))
	d, v := result.Reduce("usd", book)
	// 由于精度原因, 优化版和非优化版的decimal的值不同
	expect := decimal.NewFromInt(17).Div(decimal.NewFromFloat(2.3)).Add(decimal.NewFromFloat(6).Neg()).Div(decimal.NewFromFloat(7)).Mul(decimal.NewFromFloat(7))
	if !v || !reflect.DeepEqual(*d, expect) {
		t.Errorf("expect: %+v, got: %+v", expect, d)
	}
}

func TestDivideMultiOptimize(t *testing.T) {
	book := new(rate.Book)
	book.Add("rmb", "usd", decimal.NewFromFloat(2.3))
	result := NewMoney(decimal.NewFromInt(5), "rmb").Plus(NewMoney(decimal.NewFromInt(12), "rmb")).Minus(NewMoney(decimal.NewFromInt(6), "usd")).Divide(decimal.NewFromFloat(7)).Multiply(decimal.NewFromFloat(7))
	d, v := result.Optimize().Reduce("usd", book)
	expect := decimal.NewFromInt(17).Div(decimal.NewFromFloat(2.3)).Add(decimal.NewFromFloat(6).Neg())
	if !v || !reflect.DeepEqual(*d, expect) {
		t.Errorf("expect: %+v, got: %+v", expect, d)
	}
}
