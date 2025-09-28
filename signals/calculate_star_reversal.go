package signals

import (
	"encoding/json"
	"errors"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/indicators"
)

func calculateStarReversal(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	output := indicators.DetectStarReversal(candles)
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}
