package signals

type TradeSignal string

const (
	Buy   TradeSignal = "BUY"
	Sell  TradeSignal = "SELL"
	Hold  TradeSignal = "HOLD"
	Error TradeSignal = "ERROR"
)
