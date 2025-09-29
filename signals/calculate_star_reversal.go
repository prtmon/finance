package signals

import (
	"encoding/json"
	"errors"
	"github.com/prtmon/finance/common"
)

type StarReversalArgs struct {
	smallBodyRatio float64
	largeBodyRatio float64
}

func calculateStarReversal(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	var paramStruct StarReversalArgs
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return nil, err
	}
	output := candles.StarReversal(paramStruct.smallBodyRatio, paramStruct.largeBodyRatio)
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}
