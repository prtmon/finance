package signals

import (
	"encoding/json"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/indicators"
)

func calculateEngulfing(candles common.Candlesticks, params json.RawMessage) (int64, error) {
	ohlcv := candles.ToOhlcv()
	output := indicators.DetectEngulfing(ohlcv)
	if len(output) > 0 {
		return output[len(output)-1], nil
	}
	return 0, nil
}
