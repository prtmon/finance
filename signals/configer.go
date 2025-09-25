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
	Value int64
	Score float64
}

// 指标函数映射表
var indicatorHandlers = map[string]func(candles common.Candlesticks, jsonParam json.RawMessage) (int64, error){
	"MACross":         calculateMACross,
	"MACDCross":       calculateMACDCross,
	"StarReversal":    calculateStarReversal,
	"RsiKdjOverTrade": calculateRsiKdjOverTrade,
	"Hammer":          calculateHammer,
	"Engulfing":       calculateEngulfing,
}

func (ic *IndicatorConfig) Calculate(candles common.Candlesticks) (IndicatorResult, error) {
	handler, exists := indicatorHandlers[ic.Handler]
	if !exists {
		return IndicatorResult{}, fmt.Errorf("unknown handler: %s", ic.Handler)
	}

	value, err := handler(candles, ic.Params)
	if err != nil {
		return IndicatorResult{}, err
	}

	return IndicatorResult{
		Name:  ic.Name,
		Value: value,
		Score: tools.ToFloat64(value) * ic.Weight,
	}, nil
}

// TradeDecision buyLow=1.5,sellHigh=-1.5
func TradeDecision(candles common.Candlesticks, configs []IndicatorConfig, buyLow, sellHigh float64) (TradeSignal, float64, []IndicatorResult) {
	var totalScore float64
	var results []IndicatorResult

	for _, config := range configs {
		result, err := config.Calculate(candles)
		if err != nil {
			return Error, 0, results
		}
		totalScore += result.Score
		results = append(results, result)
	}

	if totalScore > buyLow {
		return Buy, totalScore, results
	} else if totalScore < sellHigh {
		return Sell, totalScore, results
	}
	return Hold, totalScore, results
}
