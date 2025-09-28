package signals

import (
	"encoding/json"
	"errors"
	"github.com/prtmon/finance/common"
)

func calculateEngulfing(candles common.Candlesticks, params json.RawMessage) ([]int64, error) {
	output := candles.Engulfing()
	if len(output) > 0 {
		return output, nil
	}
	return output, errors.New("the length of the K-line sample data is insufficient")
}
