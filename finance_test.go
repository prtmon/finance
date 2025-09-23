package finance

import (
	"fmt"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/signals"
)

func test(candles common.Candlesticks) {
	configs := []signals.IndicatorConfig{
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

	signal, results := signals.TradeDecision(candles, configs, 1.5, 1.5)
	fmt.Printf("Trade Signal: %s\n", signal)
	fmt.Println("Indicator Results:")
	for _, r := range results {
		fmt.Printf("- %s: Value=%.2f, Score=%.2f\n", r.Name, r.Value, r.Score)
	}
}
