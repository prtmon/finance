package indicators

import (
	"github.com/prtmon/finance/common"
	"math"
)

const (
	BodyRatio        = 0.3
	ShadowRatio      = 2.0
	UpperShadowRatio = 0.1
	TrendConfirmBars = 3
	//RequireNextConfirm = true // 是否需要下一根K线确认
)

// 锤头线检测（底部反转信号）
func isHammer(c common.Candlestick) bool {
	bodySize := math.Abs(c.Close - c.Open)
	totalRange := c.High - c.Low
	if totalRange == 0 {
		return false
	}
	lowerShadow := math.Min(c.Open, c.Close) - c.Low
	upperShadow := c.High - math.Max(c.Open, c.Close)

	// 条件1：实体较小（占比<30%）
	isSmallBody := bodySize/totalRange < BodyRatio
	// 条件2：下影线至少是实体的2倍
	hasLongLowerShadow := lowerShadow >= ShadowRatio*bodySize
	// 条件3：上影线极短（<10%总范围）
	hasTinyUpperShadow := upperShadow/totalRange < UpperShadowRatio

	return isSmallBody && hasLongLowerShadow && hasTinyUpperShadow
}

// 射击之星检测（顶部反转信号）
func isShootingStar(c common.Candlestick) bool {
	bodySize := math.Abs(c.Close - c.Open)
	totalRange := c.High - c.Low
	if totalRange == 0 {
		return false
	}

	upperShadow := c.High - math.Max(c.Open, c.Close)
	lowerShadow := math.Min(c.Open, c.Close) - c.Low

	return bodySize/totalRange < BodyRatio &&
		upperShadow >= ShadowRatio*bodySize &&
		lowerShadow/totalRange < UpperShadowRatio
}

// 趋势判断
func isTrend(candles common.Candlesticks, idx int, upTrend bool) bool {
	if idx < TrendConfirmBars {
		return false
	}

	count := 0
	for i := idx - TrendConfirmBars; i < idx; i++ {
		if upTrend && candles[i].Close > candles[i].Open {
			count++
		} else if !upTrend && candles[i].Close < candles[i].Open {
			count++
		}
	}
	return count >= TrendConfirmBars/2
}

// DetectHammerSignals 检测锤头线买入信号（需结合趋势）
func DetectHammerSignals(candles common.Candlesticks) []int64 {
	signals := make([]int64, len(candles))
	for i := 1; i < len(candles); i++ {
		current := candles[i]

		// 锤头线买入信号
		if isHammer(current) && isTrend(candles, i, false) {
			//可加入判断下一根K线是否是上涨来确认此形态
			signals[i] = 1
		}

		// 射击之星卖出信号
		if isShootingStar(current) && isTrend(candles, i, true) {
			//可加入判断下一根K线是否是下跌来确认此形态
			signals[i] = -1
		}
	}
	return signals
}
