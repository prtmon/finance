package signals

import (
	"encoding/json"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/indicators"
)

type MACDCrossArgs struct {
	InFastPeriod   int
	InSlowPeriod   int
	InSignalPeriod int
}

func calculateMACDCross(candles common.Candlesticks, params json.RawMessage) (int64, error) {
	var paramStruct MACDCrossArgs
	ohlcv := candles.ToOhlcv()
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return 0, err
	}
	output := indicators.MACDCross(ohlcv, paramStruct.InFastPeriod, paramStruct.InSlowPeriod, paramStruct.InSignalPeriod)
	if len(output) > 0 {
		return output[len(output)-1], nil
	}
	return 0, nil
}
