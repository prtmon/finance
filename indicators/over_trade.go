package indicators

import (
	"github.com/markcheno/go-talib"
	"github.com/prtmon/finance/common"
)

// RsiOverTrade Rsi多空力量-超买超卖
// inTimePeriod  一般配置14,周期，短线配置6，中线配置14，长线配置24
// 一般大于highValue70为超买，小于lowValue30为超卖
func RsiOverTrade(closes []float64, inTimePeriod int, lowValue, highValue float64) []int64 {
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

// KdjOverTrade 默认highValue=80,lowValue=20
func KdjOverTrade(ohlcv common.OHLCV, inTimePeriod int, lowValue, highValue float64) []int64 {
	k, d := talib.Stoch(ohlcv.High, ohlcv.Low, ohlcv.Close, 14, 3, talib.EMA, 3, talib.EMA)
	signals := make([]int64, len(k))
	for i := 0; i < len(k); i++ {
		// RSI超买且KDJ超买
		if k[i] > highValue && d[i] > highValue {
			signals[i] = -1
		} else if k[i] < lowValue && d[i] < lowValue {
			signals[i] = 1
		} else {
			signals[i] = 0
		}
	}
	return signals
}

// RsiKdjOverTrade 双确认策略结合RSI超买超卖区(30/70)和KDJ超买超卖区(20/80)生成交易信号
func RsiKdjOverTrade(ohlcv common.OHLCV, inTimePeriod int, rsiLowValue, rsiHighValue, kdjLowValue, kdjHighValue float64) []int64 {
	rsi := talib.Rsi(ohlcv.Close, inTimePeriod)
	k, d := talib.Stoch(ohlcv.High, ohlcv.Low, ohlcv.Close, 14, 3, talib.EMA, 3, talib.EMA)
	signals := make([]int64, len(rsi))
	for i := 0; i < len(rsi); i++ {
		// RSI超买且KDJ超买
		if rsi[i] > rsiHighValue && k[i] > kdjHighValue && d[i] > kdjHighValue {
			signals[i] = -1
		} else if rsi[i] < rsiLowValue && k[i] < kdjLowValue && d[i] < kdjLowValue {
			signals[i] = 1
		} else {
			signals[i] = 0
		}
	}
	return signals
}
