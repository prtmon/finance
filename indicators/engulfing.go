package indicators

import (
	"github.com/prtmon/finance/common"
	"github.com/prtmon/utils"
)

/*
DetectEngulfing 判断是否为吞没形态,0=不是吞没形态,1=看涨吞没，-1=看跌吞没
形态判断标准：
1，由两根K线组成
2，两根K线实体颜色相反
3，第二根K线实体完全覆盖第一根K线实体
通常出现在明确的趋势之后
*/
func DetectEngulfing(ohlcv common.OHLCV) []int64 {
	output := make([]int64, len(ohlcv.Close))
	if len(output) < 2 {
		return output
	}

	for i := 1; i < len(output); i++ {

		prevBody := utils.Abs(ohlcv.Close[i-1] - ohlcv.Open[i-1])
		currentBody := utils.Abs(ohlcv.Open[i] - ohlcv.Open[i])

		isEngulfing := false
		output[i] = 0
		// 第二根K线实体必须大于第一根
		if currentBody <= prevBody {
			output[i] = 0
		}

		// 完全覆盖前一根K线实体
		if (ohlcv.Close[i] > ohlcv.Open[i-1] && ohlcv.Close[i] > ohlcv.Close[i-1]) &&
			(ohlcv.Open[i] < ohlcv.Open[i-1] && ohlcv.Open[i] < ohlcv.Close[i-1]) {
			isEngulfing = true
		}

		if (ohlcv.Close[i] < ohlcv.Open[i-1] && ohlcv.Close[i] < ohlcv.Close[i-1]) &&
			(ohlcv.Open[i] > ohlcv.Open[i-1] && ohlcv.Open[i] > ohlcv.Close[i-1]) {
			isEngulfing = true
		}
		if isEngulfing {
			//看涨吞没形态
			if ohlcv.Close[i] > ohlcv.Open[i] && ohlcv.Close[i-1] < ohlcv.Open[i-1] {
				output[i] = 1
			}
			//看跌吞没形态
			if ohlcv.Close[i] < ohlcv.Open[i] && ohlcv.Close[i-1] > ohlcv.Open[i-1] {
				output[i] = -1
			}
		}
	}

	return output
}
