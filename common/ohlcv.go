package common

import "github.com/markcheno/go-talib"

// dataframe K线数据结构
type OHLCV struct {
	Time   []float64
	Open   []float64
	High   []float64
	Low    []float64
	Close  []float64
	Volume []float64
}

func (o OHLCV) MA(inTimePeriod int, maType talib.MaType) []float64 {
	return talib.Ma(o.Close, inTimePeriod, maType)
}

func (o OHLCV) EMA(inTimePeriod int) []float64 {
	return talib.Ema(o.Close, inTimePeriod)
}
