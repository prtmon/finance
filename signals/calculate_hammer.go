package signals

import (
	"encoding/json"
	"errors"
	"github.com/prtmon/finance/common"
)

type HammerTrendArgs struct {
	TrendConfirmBars int
	smallBodyRatio   float64
	largeShadowRatio float64
	smallShadowRatio float64
}

func calculateHammer(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	var paramStruct HammerTrendArgs
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return nil, err
	}
	output := candles.HammerTrend(paramStruct.TrendConfirmBars, paramStruct.smallBodyRatio, paramStruct.largeShadowRatio, paramStruct.smallShadowRatio)
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}
