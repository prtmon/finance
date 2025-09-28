package common

import (
	"github.com/markcheno/go-talib"
)

// dataframe K线数据结构
type OHLCV struct {
	Time   []float64
	Open   []float64
	High   []float64
	Low    []float64
	Close  []float64
	Volume []float64
}

func (o OHLCV) ToCandlesticks() Candlesticks {
	var candlesticks Candlesticks
	for i := 0; i < len(o.Close)-1; i++ {
		candlesticks = append(candlesticks, Candlestick{
			Time:   o.Time[i],
			Open:   o.Open[i],
			High:   o.High[i],
			Low:    o.Low[i],
			Close:  o.Close[i],
			Volume: o.Volume[i],
		})
	}
	return candlesticks
}

func (o OHLCV) MA(inTimePeriod int, maType talib.MaType) []float64 {
	return talib.Ma(o.Close, inTimePeriod, maType)
}

func (o OHLCV) EMA(inTimePeriod int) []float64 {
	return talib.Ema(o.Close, inTimePeriod)
}

func (o OHLCV) Engulfing() []int64 {
	c := o.ToCandlesticks()
	return c.Engulfing()
}

// EMACross 均线交叉,短单位均线1下穿长单位均线2,下降趋势,反之为上升趋势 inFastPeriod int, inSlowPeriod int
func (o OHLCV) EMACross(inFastPeriod, inSlowPeriod int) []int64 {
	closes := o.Close
	output := make([]int64, len(closes))
	if inFastPeriod >= inSlowPeriod {
		return output
	}
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
func (o OHLCV) MACDCross(inFastPeriod, inSlowPeriod, inSignalPeriod int) []int64 {
	closes := o.Close
	output := make([]int64, len(closes))
	// MACD指标
	macd, macdSignal, _ := talib.Macd(closes, inFastPeriod, inSlowPeriod, inSignalPeriod)
	for i := 1; i < len(macd); i++ {
		output[i] = 0
		// MACD金叉 = 看多反转
		if macd[i] > macdSignal[i] && macd[i-1] <= macdSignal[i-1] {
			output[i] = 1
		}
		// MACD死叉 = 看跌反转
		if macd[i] < macdSignal[i] &&
			macd[i-1] >= macdSignal[i-1] {
			output[i] = -1
		}
	}
	return output
}

// PeakTrough 滑动窗口极值检测,峰谷峰顶检测
func (o OHLCV) PeakTrough(window int) []int64 {
	c := o.ToCandlesticks()
	return c.PeakTrough(window)
}

// HammerTrend 锤头线趋势（反转信号）
func (o OHLCV) HammerTrend(TrendConfirmBars int, smallBodyRatio, largeShadowRatio, smallShadowRatio float64) []int64 {
	c := o.ToCandlesticks()
	return c.HammerTrend(TrendConfirmBars, smallBodyRatio, largeShadowRatio, smallShadowRatio)
}

// OverTradeRsi Rsi多空力量-超买超卖
// inTimePeriod  一般配置14,周期，短线配置6，中线配置14，长线配置24
// 一般大于highValue70为超买，小于lowValue30为超卖
func (o OHLCV) OverTradeRsi(inTimePeriod int, lowValue, highValue float64) []int64 {
	closes := o.Close
	signals := make([]int64, len(closes))
	//相对强弱指数(RSI)
	rsi := talib.Rsi(closes, inTimePeriod)
	for i := 0; i < len(rsi); i++ {
		k := rsi[i]
		signals[i] = 0
		if k < lowValue {
			signals[i] = -1
		} else if k > highValue {
			signals[i] = 1
		}
	}
	return signals
}

// OverTradeKdj 默认inFastKPeriod=14, inSlowKPeriod=3 ,inSlowDPeriod = 3 ,highValue=80,lowValue=20
func (o OHLCV) OverTradeKdj(inFastKPeriod, inSlowKPeriod, inSlowDPeriod int, lowValue, highValue float64) []int64 {
	k, d, j := KDJ(o, inFastKPeriod, inSlowKPeriod, talib.EMA, inSlowDPeriod, talib.EMA)
	signals := make([]int64, len(k))
	for i := 0; i < len(k); i++ {
		// KDJ超买
		if k[i] > highValue && d[i] > highValue && j[i] > 100 {
			signals[i] = -1
		} else if k[i] < lowValue && d[i] < lowValue && j[i] < 0 {
			signals[i] = 1
		} else {
			signals[i] = 0
		}
	}
	return signals
}

// StarReversal smallBodyRatio=0.3, largeBodyRatio=0.7
func (o OHLCV) StarReversal(smallBodyRatio, largeBodyRatio float64) []int64 {
	c := o.ToCandlesticks()
	return c.StarReversal(smallBodyRatio, largeBodyRatio)
}

func (o OHLCV) ThreeWhiteSoldier() []bool {
	c := o.ToCandlesticks()
	return c.ThreeWhiteSoldier()
}
