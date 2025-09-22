package indicators

import (
	"github.com/prtmon/finance/common"
	"github.com/prtmon/tools"
)

// 实体类型判断（新增阈值常量）
const (
	LongBodyThreshold  = 0.7
	ShortBodyThreshold = 0.3
)

// StarReversal 检测星线反转形态（结构体封装+条件复用）
type StarReversal struct {
	MorningConditions bool
	EveningConditions bool
}

// 检测单根K线是否为长实体
func isLongBody(c common.Candlestick) bool {
	return tools.Abs(c.Close-c.Open)/(c.High-c.Low) > LongBodyThreshold
}

// 检测单根K线是否为小实体
func isShortBody(c common.Candlestick) bool {
	return tools.Abs(c.Close-c.Open)/(c.High-c.Low) < ShortBodyThreshold
}

func (sr *StarReversal) Validate(candles common.Candlesticks, i int) {
	prev2, prev1, curr := candles[i-2], candles[i-1], candles[i]

	// 早晨之星条件
	sr.MorningConditions = isLongBody(prev2) && prev2.Close < prev2.Open &&
		isShortBody(prev1) && prev1.High < prev2.Low && prev1.Low < prev2.Close &&
		isLongBody(curr) && curr.Close > prev2.Open && curr.Close > prev1.Close

	// 黄昏之星条件
	sr.EveningConditions = isLongBody(prev2) && prev2.Close > prev2.Open &&
		isShortBody(prev1) && prev1.High > prev2.High && prev1.Low > prev2.Close &&
		isLongBody(curr) && curr.Close < prev2.Open && curr.Close < prev1.Close
}

func DetectStarReversal(candles common.Candlesticks) []int64 {
	output := make([]int64, len(candles))
	for i := 2; i < len(candles); i++ {
		sr := StarReversal{}
		sr.Validate(candles, i)

		if sr.MorningConditions {
			output[i] = 1
		} else if sr.EveningConditions {
			output[i] = -1
		} else {
			output[i] = 0
		}
	}
	return output
}
