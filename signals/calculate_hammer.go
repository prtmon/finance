package signals

import (
	"encoding/json"
	"errors"
	"github.com/prtmon/finance/common"
)

type HammerTrendArgs struct {
	TrendConfirmBars int
	SmallBodyRatio   float64
	LargeShadowRatio float64
	SmallShadowRatio float64
}

func calculateHammer(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	var paramStruct HammerTrendArgs
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return nil, err
	}
	output := candles.HammerTrend(paramStruct.TrendConfirmBars, paramStruct.SmallBodyRatio, paramStruct.LargeShadowRatio, paramStruct.SmallShadowRatio)
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}
