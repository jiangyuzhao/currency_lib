package rate

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

var (
	One = decimal.NewFromInt(1)
)

// 非并发安全, 看使用场景是否需要并发使用
type Book struct {
	innerTable map[string]*decimal.Decimal
}

// Add添加一条记录, 可以生成相反的两条记录, 但目前没有这么做, 看需求
func (b *Book) Add(from, to string, rate decimal.Decimal) {
	if b.innerTable == nil {
		b.innerTable = make(map[string]*decimal.Decimal)
	}
	key := fmt.Sprintf("%s:%s", strings.ToLower(from), strings.ToLower(to))
	b.innerTable[key] = &rate
}

func (b *Book) GetRate(from, to string) (*decimal.Decimal, error) {
	from = strings.ToLower(from)
	to = strings.ToLower(to)
	if from == to {
		return &One, nil
	}
	if b.innerTable == nil {
		// error
		return nil, errors.New("no rate book")
	}
	key := from + ":" + to
	return b.innerTable[key], nil
}
