package common

// Candlestick 定义K线数据结构
type Candlestick struct {
	Time   float64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	//预计算字段
	//RealBody    float64 //abs(收盘价-开盘价)实体
	//UpperShadow float64 //上影线
	//LowerShadow float64 //下影线
	//IsBullish   bool    //涨或跌
}

type Candlesticks []Candlestick

func (c Candlesticks) Closes() []float64 {
	closes := make([]float64, len(c))
	for i, v := range c {
		closes[i] = v.Close
	}
	return closes
}

func (c Candlesticks) Opens() []float64 {
	opens := make([]float64, len(c))
	for i, v := range c {
		opens[i] = v.Open
	}
	return opens
}

func (c Candlesticks) Highs() []float64 {
	highs := make([]float64, len(c))
	for i, v := range c {
		highs[i] = v.High
	}
	return highs
}

func (c Candlesticks) Lows() []float64 {
	lows := make([]float64, len(c))
	for i, v := range c {
		lows[i] = v.Low
	}
	return lows
}

func (c Candlesticks) ToOhlcv() OHLCV {
	ohlcv := OHLCV{
		Time:   make([]float64, len(c)),
		Open:   make([]float64, len(c)),
		High:   make([]float64, len(c)),
		Low:    make([]float64, len(c)),
		Close:  make([]float64, len(c)),
		Volume: make([]float64, len(c)),
	}

	for i, k := range c {
		ohlcv.Time[i] = k.Time
		ohlcv.Open[i] = k.Open
		ohlcv.High[i] = k.High
		ohlcv.Low[i] = k.Low
		ohlcv.Close[i] = k.Close
		ohlcv.Volume[i] = k.Volume
	}

	return ohlcv
}
