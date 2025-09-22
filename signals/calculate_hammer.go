package signals

import (
	"encoding/json"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/indicators"
)

func calculateHammer(candles common.Candlesticks, params json.RawMessage) (int64, error) {
	output := indicators.DetectHammerSignals(candles)
	if len(output) > 0 {
		return output[len(output)-1], nil
	}
	return 0, nil
}
