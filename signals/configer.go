package signals

import (
	"encoding/json"
	"fmt"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/tools"
)

type IndicatorConfig struct {
	Name        string          `json:"name"`
	Handler     string          `json:"handler"`
	Params      json.RawMessage `json:"params"`
	Weight      float64         `json:"weight"`
	Description string          `json:"description"`
}

type IndicatorResult struct {
	Name  string
	Value []int64
	Score []float64
}

// 指标函数映射表
var indicatorHandlers = map[string]func(candles common.Candlesticks, jsonParam json.RawMessage) ([]int64, error){
	"MACross":      calculateMACross,
	"MACDCross":    calculateMACDCross,
	"StarReversal": calculateStarReversal,
	"OverTradeRsi": calculateOverTradeRsi,
	"OverTradeKdj": calculateOverTradeKdj,
	"Hammer":       calculateHammer,
	"Engulfing":    calculateEngulfing,
}

func (ic *IndicatorConfig) Calculate(candles common.Candlesticks) (IndicatorResult, error) {
	handler, exists := indicatorHandlers[ic.Handler]
	if !exists {
		return IndicatorResult{}, fmt.Errorf("unknown handler: %s", ic.Handler)
	}

	values, err := handler(candles, ic.Params)
	if err != nil {
		return IndicatorResult{}, err
	}
	scores := make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		scores[i] = tools.ToFloat64(values[i]) * ic.Weight
	}
	return IndicatorResult{
		Name:  ic.Name,
		Value: values,
		Score: scores,
	}, nil
}

// TradeDecision buyLow=1.5,sellHigh=-1.5
func TradeDecision(candles common.Candlesticks, configs []IndicatorConfig, buyLow, sellHigh float64) ([]TradeSignal, []float64, []IndicatorResult) {

	var results []IndicatorResult
	tradeSignals := make([]TradeSignal, len(candles))
	totalScores := make([]float64, len(candles))

	for _, config := range configs {
		result, err := config.Calculate(candles)
		if err != nil {
			return nil, nil, results
		}
		results = append(results, result)
	}

	for i := 0; i < len(candles); i++ {
		for _, result := range results {
			totalScores[i] += result.Score[i]
		}
		if totalScores[i] > buyLow {
			tradeSignals[i] = Buy
		} else if totalScores[i] < sellHigh {
			tradeSignals[i] = Sell
		} else {
			tradeSignals[i] = Hold
		}
	}
	return tradeSignals, totalScores, results
}
