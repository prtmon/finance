package signals

import (
	"encoding/json"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/indicators"
)

type RsiKdjOverTradeArgs struct {
	InTimePeriod int
	RsiLowValue  float64
	RsiHighValue float64
	KdjLowValue  float64
	KdjHighValue float64
}

func calculateRsiKdjOverTrade(candles common.Candlesticks, params json.RawMessage) (int64, error) {
	var paramStruct RsiKdjOverTradeArgs
	ohlcv := candles.ToOhlcv()
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return 0, err
	}
	output := indicators.RsiKdjOverTrade(ohlcv, paramStruct.InTimePeriod, paramStruct.RsiLowValue, paramStruct.RsiHighValue, paramStruct.KdjLowValue, paramStruct.KdjHighValue)
	if len(output) > 0 {
		return output[len(output)-1], nil
	}
	return 0, nil
}
