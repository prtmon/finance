package common

import (
	"github.com/prtmon/tools"
	"math"
)

// Candlestick 定义K线数据结构
type Candlestick struct {
	Time   float64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// RealBody 实体大小
func (c Candlestick) RealBody() float64 {
	return tools.Abs(c.Close - c.Open)
}

// TotalBody 实体大小
func (c Candlestick) TotalBody() float64 {
	return tools.Abs(c.High - c.Close)
}

// IsSmallBody 是否小实体
func (c Candlestick) IsSmallBody(bodyRatio float64) bool {
	rate := math.Abs(bodyRatio)
	if rate == 0 {
		rate = 0.3
	}
	return c.RealBody()/c.TotalBody() < rate
}

// IsLargeBody 是否大实体
func (c Candlestick) IsLargeBody(bodyRatio float64) bool {
	rate := math.Abs(bodyRatio)
	if rate == 0 {
		rate = 0.7
	}
	return c.RealBody()/c.TotalBody() > rate
}

// UpperShadow 上影线
func (c Candlestick) UpperShadow() float64 {
	return c.High - math.Max(c.Open, c.Close)
}

// LowerShadow 下影线
func (c Candlestick) LowerShadow() float64 {
	return math.Min(c.Open, c.Close) - c.Low
}

// IsBullish 是否上涨
func (c Candlestick) IsBullish() bool {
	return c.Close > c.Open
}
