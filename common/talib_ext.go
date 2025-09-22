package common

import "github.com/markcheno/go-talib"

// 6. 随机指标(KDJ),默认：inFastKPeriod = 14,inSlowKPeriod = 3,inSlowDMatype=3,talib.EMA
func KDJ(ohlcv OHLCV, inFastKPeriod, inSlowKPeriod int, inSlowKMatype talib.MaType, inSlowDPeriod int, inSlowDMatype talib.MaType) ([]float64, []float64, []float64) {
	slowK, slowD := talib.Stoch(ohlcv.High, ohlcv.Low, ohlcv.Close, inFastKPeriod, inSlowKPeriod, inSlowKMatype, inSlowDPeriod, inSlowDMatype)
	slowJ := calculateJ(slowK, slowD)
	return slowK, slowD, slowJ
}

func calculateJ(k, d []float64) []float64 {
	j := make([]float64, len(k))
	for i := range k {
		j[i] = 3*k[i] - 2*d[i]
	}
	return j
}
