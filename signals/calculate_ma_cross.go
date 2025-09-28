package signals

import (
	"encoding/json"
	"errors"
	"github.com/prtmon/finance/common"
)

type MACrossArgs struct {
	InFastPeriod int
	InSlowPeriod int
	MaType       int
}

func calculateMACross(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	var paramStruct MACrossArgs
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return nil, err
	}
	output := candles.EMACross(paramStruct.InFastPeriod, paramStruct.InSlowPeriod)
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}
