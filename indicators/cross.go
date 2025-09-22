package indicators

import (
	"github.com/markcheno/go-talib"
	"github.com/prtmon/finance/common"
)

// MACDCross 金叉死叉
// inFastPeriod=12,inSlowPeriod=26,inSlowPeriod=9
func MACDCross(ohlcv common.OHLCV, inFastPeriod, inSlowPeriod, inSignalPeriod int) []int64 {
	closes := ohlcv.Close
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

// MACross 均线交叉,短单位均线1下穿长单位均线2,下降趋势,反之为上升趋势 inFastPeriod int, inSlowPeriod int, inSignalPeriod int
func MACross(closes []float64, inFastPeriod, inSlowPeriod int, maType int) []int64 {
	output := make([]int64, len(closes))
	if inFastPeriod >= inSlowPeriod {
		return output
	}
	fastMA := talib.Ma(closes, inFastPeriod, talib.MaType(maType))
	slowMA := talib.Ma(closes, inSlowPeriod, talib.MaType(maType))
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
