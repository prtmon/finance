package signals

import (
	"encoding/json"
	"errors"
	"github.com/prtmon/finance/common"
)

type StarReversalArgs struct {
	SmallBodyRatio float64
	LargeBodyRatio float64
}

func calculateStarReversal(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	var paramStruct StarReversalArgs
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return nil, err
	}
	output := candles.StarReversal(paramStruct.SmallBodyRatio, paramStruct.LargeBodyRatio)
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}
