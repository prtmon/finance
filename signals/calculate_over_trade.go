package signals

import (
	"encoding/json"
	"errors"
	"github.com/prtmon/finance/common"
)

type OverTradeRsiArgs struct {
	InTimePeriod int
	RsiLowValue  float64
	RsiHighValue float64
}

func calculateOverTradeRsi(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	var paramStruct OverTradeRsiArgs
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return nil, err
	}
	output := candles.OverTradeRsi(paramStruct.InTimePeriod, paramStruct.RsiLowValue, paramStruct.RsiHighValue)
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}

type OverTradeKdjArgs struct {
	InFastKPeriod int
	InSlowKPeriod int
	InSlowDPeriod int
	KdjLowValue   float64
	KdjHighValue  float64
}

func calculateOverTradeKdj(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	var paramStruct OverTradeKdjArgs
	err := json.Unmarshal(params, &paramStruct) // 反序列化为RawMessage
	if err != nil {
		return nil, err
	}
	output := candles.OverTradeKdj(paramStruct.InFastKPeriod, paramStruct.InSlowKPeriod, paramStruct.InSlowDPeriod, paramStruct.KdjLowValue, paramStruct.KdjHighValue)
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}
