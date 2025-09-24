package finance

import (
	"fmt"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/signals"
)

func test(candles common.Candlesticks) {
	configs := []signals.IndicatorConfig{
		{
			Name:        "MACDCross",
			Handler:     "MACDCross",
			Params:      []byte(`{"InFastPeriod":12,"InSlowPeriod":26,"InSignalPeriod":9}`),
			Weight:      1,
			Description: "MACD金叉死叉信号",
		},
		{
			Name:        "MACross",
			Handler:     "MACross",
			Params:      []byte(`{"InFastPeriod":5,"InSlowPeriod":30,"MaType":1}`),
			Weight:      1,
			Description: "均线交叉信号",
		},
		{
			Name:        "RsiKdjOverTrade",
			Handler:     "RsiKdjOverTrade",
			Params:      []byte(`{}`),
			Weight:      1,
			Description: "rsi结合kdj超买超卖信号",
		},
		{
			Name:        "StarReversal",
			Handler:     "StarReversal",
			Params:      []byte(`{}`),
			Weight:      1,
			Description: "早晨之星/黄昏之星反转信号",
		},
		{
			Name:        "Engulfing",
			Handler:     "Engulfing",
			Params:      []byte(`{}`),
			Weight:      1,
			Description: "吞没型态信号",
		},
		{
			Name:        "Hammer",
			Handler:     "Hammer",
			Params:      []byte(`{}`),
			Weight:      1,
			Description: "锤头线买卖信号",
		},
	}

	signal, results := signals.TradeDecision(candles, configs, 3.5, 2.5)
	fmt.Printf("Trade Signal: %s\n", signal)
	fmt.Println("Indicator Results:")
	for _, r := range results {
		fmt.Printf("- %s: Value=%.2f, Score=%.2f\n", r.Name, r.Value, r.Score)
	}
}
