package signals

import (
	"encoding/json"
	"fmt"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/utils"
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
		Score: utils.ToFloat64(value) * ic.Weight,
	}, nil
}

// TradeDecision buyLow=1.5,sellHigh=-1.5
func TradeDecision(candles common.Candlesticks, configs []IndicatorConfig, buyLow, sellHigh float64) (TradeSignal, []IndicatorResult) {
	var totalScore float64
	var results []IndicatorResult

	for _, config := range configs {
		result, err := config.Calculate(candles)
		if err != nil {
			return Error, results
		}
		totalScore += result.Score
		results = append(results, result)
	}

	if totalScore > buyLow {
		return Buy, results
	} else if totalScore < sellHigh {
		return Sell, results
	}
	return Hold, results
}

/*
func test() {
	configs := []IndicatorConfig{
		{
			Name:        "Relative Strength Index",
			Handler:     "RSI",
			Params:      []byte(`{"period":14}`),
			Weight:      0.4,
			Description: "超买超卖指标",
		},
		{
			Name:        "Moving Average Convergence Divergence",
			Handler:     "MACD",
			Params:      []byte(`{"fast":12,"slow":26}`),
			Weight:      0.3,
			Description: "趋势指标",
		},
		{
			Name:        "Bollinger Bands",
			Handler:     "Bollinger",
			Params:      []byte(`{"period":20}`),
			Weight:      0.3,
			Description: "波动率指标",
		},
	}

	signal, results := TradeDecision(configs)
	fmt.Printf("Trade Signal: %s\n", signal)
	fmt.Println("Indicator Results:")
	for _, r := range results {
		fmt.Printf("- %s: Value=%.2f, Score=%.2f\n", r.Name, r.Value, r.Score)
	}
}*/
