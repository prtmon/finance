package signals

import (
	"encoding/json"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/indicators"
)

type MACrossArgs struct {
	InReal       []float64
	InFastPeriod int
	InSlowPeriod int
	MaType       int
}

func calculateMACross(candles common.Candlesticks, params json.RawMessage) (int64, error) {
	var paramStruct MACrossArgs
	ohlcv := candles.ToOhlcv()
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return 0, err
	}
	output := indicators.MACross(ohlcv.Close, paramStruct.InFastPeriod, paramStruct.InSlowPeriod, paramStruct.MaType)
	if len(output) > 0 {
		return output[len(output)-1], nil
	}
	return 0, nil
}
