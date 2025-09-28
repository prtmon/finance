package common

import (
	"github.com/markcheno/go-talib"
	"github.com/prtmon/tools"
)

type Candlesticks []Candlestick

func (c Candlesticks) Closes() []float64 {
	closes := make([]float64, len(c))
	for i, v := range c {
		closes[i] = v.Close
	}
	return closes
}

func (c Candlesticks) Opens() []float64 {
	opens := make([]float64, len(c))
	for i, v := range c {
		opens[i] = v.Open
	}
	return opens
}

func (c Candlesticks) Highs() []float64 {
	highs := make([]float64, len(c))
	for i, v := range c {
		highs[i] = v.High
	}
	return highs
}

func (c Candlesticks) Lows() []float64 {
	lows := make([]float64, len(c))
	for i, v := range c {
		lows[i] = v.Low
	}
	return lows
}

func (c Candlesticks) ToOhlcv() OHLCV {
	ohlcv := OHLCV{
		Time:   make([]float64, len(c)),
		Open:   make([]float64, len(c)),
		High:   make([]float64, len(c)),
		Low:    make([]float64, len(c)),
		Close:  make([]float64, len(c)),
		Volume: make([]float64, len(c)),
	}

	for i, k := range c {
		ohlcv.Time[i] = k.Time
		ohlcv.Open[i] = k.Open
		ohlcv.High[i] = k.High
		ohlcv.Low[i] = k.Low
		ohlcv.Close[i] = k.Close
		ohlcv.Volume[i] = k.Volume
	}

	return ohlcv
}

/*
Engulfing 判断是否为吞没形态,0=不是吞没形态,1=看涨吞没，-1=看跌吞没
形态判断标准：
1，由两根K线组成
2，两根K线实体颜色相反
3，第二根K线实体完全覆盖第一根K线实体
通常出现在明确的趋势之后
*/
func (c Candlesticks) Engulfing() []int64 {
	output := make([]int64, len(c))
	if len(output) < 2 {
		return output
	}

	for i := 1; i < len(output); i++ {
		prevCandle := c[i-1]
		currentCandle := c[i]

		prevBody := tools.Abs(prevCandle.Close - prevCandle.Open)
		currentBody := tools.Abs(currentCandle.Close - currentCandle.Open)

		isEngulfing := false
		output[i] = 0
		// 第二根K线实体必须大于第一根
		if currentBody <= prevBody {
			output[i] = 0
		}

		// 完全覆盖前一根K线实体
		if (currentCandle.Close > prevCandle.Open && currentCandle.Close > prevCandle.Close) &&
			(currentCandle.Open < prevCandle.Open && currentCandle.Open < prevCandle.Close) {
			isEngulfing = true
		}

		if (currentCandle.Close < prevCandle.Open && currentCandle.Close < prevCandle.Close) &&
			(currentCandle.Open > prevCandle.Open && currentCandle.Open > prevCandle.Close) {
			isEngulfing = true
		}
		if isEngulfing {
			//看涨吞没形态
			if currentCandle.Close > currentCandle.Open && prevCandle.Close < prevCandle.Open {
				output[i] = 1
			}
			//看跌吞没形态
			if currentCandle.Close < currentCandle.Open && prevCandle.Close > prevCandle.Open {
				output[i] = -1
			}
		}
	}

	return output
}

// EMACross 均线交叉,短单位均线1下穿长单位均线2,下降趋势,反之为上升趋势 inFastPeriod int, inSlowPeriod int
func (c Candlesticks) EMACross(inFastPeriod, inSlowPeriod int) []int64 {
	output := make([]int64, len(c))
	if inFastPeriod >= inSlowPeriod {
		return output
	}
	closes := c.Closes()
	fastMA := talib.Ma(closes, inFastPeriod, talib.EMA)
	slowMA := talib.Ma(closes, inSlowPeriod, talib.EMA)
	for i := range closes {
		if i > inSlowPeriod {
			if fastMA[i] > slowMA[i] && fastMA[i-1] < slowMA[i-1] {
				//上穿
				output[i] = 1
			} else if fastMA[i] < slowMA[i] && fastMA[i-1] > slowMA[i-1] {
				//下穿
				output[i] = -1
			}
		} else {
			output[i] = 0
		}
	}
	return output
}

// MACDCross 金叉死叉 inFastPeriod=12,inSlowPeriod=26,inSlowPeriod=9
func (c Candlesticks) MACDCross(inFastPeriod, inSlowPeriod, inSignalPeriod int) []int64 {
	o := c.ToOhlcv()
	return o.MACDCross(inFastPeriod, inSlowPeriod, inSignalPeriod)
}

// PeakTrough 滑动窗口极值检测,峰谷峰顶检测
func (c Candlesticks) PeakTrough(window int) []int64 {
	n := len(c)
	output := make([]int64, n)
	for i := window; i < n-window; i++ {
		current := c[i]
		isPeak, isTrough := true, true

		// 检查窗口范围内极值
		for j := i - window; j <= i+window; j++ {
			if j == i {
				continue
			}
			if current.High <= c[j].High {
				isPeak = false
			}
			if current.Low >= c[j].Low {
				isTrough = false
			}
		}
		if isPeak {
			//峰顶
			output[i] = 1
		}
		if isTrough {
			//谷底
			output[i] = -1
		}
	}
	return output
}

// HammerTrend 锤头线趋势（反转信号）smallBodyRatio = 0.3, largeShadowRatio = 2.0 smallShadowRatio = 0.1 TrendConfirmBars = 5
// smallBodyRatio 实体body小于整体body 30%
// largeShadowRatio 长影线至少是实体的2倍
// smallShadowRatio 短影线极短（<10%总范围）
// TrendConfirmBars 参与计算趋势bar的数量
func (c Candlesticks) HammerTrend(TrendConfirmBars int, smallBodyRatio, largeShadowRatio, smallShadowRatio float64) []int64 {
	output := make([]int64, len(c))
	for i := 1; i < len(c); i++ {
		// 判断趋势
		if i < TrendConfirmBars {
			output[i] = 0
			continue
		}
		current := c[i]

		if current.TotalBody() == 0 {
			output[i] = 0
			continue
		}

		// 判断是否是锤子线形态
		// 条件2：下影线至少是实体的2倍
		// 条件3：上影线极短（<10%总范围）
		if current.IsSmallBody(smallBodyRatio) && current.LowerShadow() >= largeShadowRatio*current.RealBody() && current.UpperShadow()/current.TotalBody() < smallShadowRatio {
			//判断过去TrendConfirmBars个bar是否是上涨趋势
			count := 0
			for j := i - TrendConfirmBars; j < i; j++ {
				if c[j].Close < c[j].Open {
					count++
				}
			}
			if count >= TrendConfirmBars/2 {
				output[i] = 1
			}
		}

		// 判断射击之星形态
		// 条件2：上影线至少是实体的2倍
		// 条件3：下影线极短（<10%总范围）
		if current.IsSmallBody(smallBodyRatio) && current.UpperShadow() >= largeShadowRatio*current.RealBody() && current.LowerShadow()/current.TotalBody() < smallShadowRatio {
			//判断过去TrendConfirmBars个bar是否是上涨趋势
			count := 0
			for j := i - TrendConfirmBars; j < i; j++ {
				if c[j].Close > c[j].Open {
					count++
				}
			}
			if count >= TrendConfirmBars/2 {
				output[i] = -1
			}
		}
	}
	return output
}

// OverTradeRsi Rsi多空力量-超买超卖
func (c Candlesticks) OverTradeRsi(inTimePeriod int, lowValue, highValue float64) []int64 {
	ohlcv := c.ToOhlcv()
	return ohlcv.OverTradeRsi(inTimePeriod, lowValue, highValue)
}

// OverTradeKdj 默认inFastKPeriod=14, inSlowKPeriod=3 ,inSlowDPeriod = 3 ,highValue=80,lowValue=20
func (c Candlesticks) OverTradeKdj(inFastKPeriod, inSlowKPeriod, inSlowDPeriod int, lowValue, highValue float64) []int64 {
	ohlcv := c.ToOhlcv()
	return ohlcv.OverTradeKdj(inFastKPeriod, inSlowKPeriod, inSlowDPeriod, lowValue, highValue)
}

// StarReversal smallBodyRatio=0.3, largeBodyRatio=0.7
func (c Candlesticks) StarReversal(smallBodyRatio, largeBodyRatio float64) []int64 {
	output := make([]int64, len(c))
	for i := 2; i < len(c); i++ {
		prev2, prev1, curr := c[i-2], c[i-1], c[i]

		// 早晨之星条件
		morningConditions := prev2.IsLargeBody(largeBodyRatio) && prev2.Close < prev2.Open &&
			prev1.IsSmallBody(smallBodyRatio) && prev1.High < prev2.Low && prev1.Low < prev2.Close &&
			curr.IsLargeBody(largeBodyRatio) && curr.Close > prev2.Open && curr.Close > prev1.Close

		// 黄昏之星条件
		eveningConditions := prev2.IsLargeBody(largeBodyRatio) && prev2.Close > prev2.Open &&
			prev1.IsSmallBody(smallBodyRatio) && prev1.High > prev2.High && prev1.Low > prev2.Close &&
			curr.IsLargeBody(largeBodyRatio) && curr.Close < prev2.Open && curr.Close < prev1.Close

		if morningConditions {
			output[i] = 1
		} else if eveningConditions {
			output[i] = -1
		}
	}
	return output
}

// ThreeWhiteSoldier 红三兵识别 (Three White Soldiers)
func (c Candlesticks) ThreeWhiteSoldier() []bool {
	output := make([]bool, len(c))
	if len(c) < 3 {
		return output
	}

	for i := 2; i < len(c); i++ {
		c1 := c[i-2]
		c2 := c[i-1]
		c3 := c[i]

		// 1. 连续三根阳线 (收盘 > 开盘)
		condition1 := (c1.Close > c1.Open) && (c2.Close > c2.Open) && (c3.Close > c3.Open)

		// 2. 每根K线收盘价高于前一根收盘价
		condition2 := c2.Close > c1.Close && c3.Close > c2.Close

		// 3. 每根K线开盘价在前一根实体范围内
		condition3 := c2.Open > c1.Open && c2.Open < c1.Close &&
			c3.Open > c2.Open && c3.Open < c2.Close

		// 4. 实体长度递增 (可选)
		condition4 := (c2.Close-c2.Open) > (c1.Close-c1.Open) &&
			(c3.Close-c3.Open) > (c2.Close-c2.Open)

		output[i] = condition1 && condition2 && condition3 && condition4
	}
	return output
}
