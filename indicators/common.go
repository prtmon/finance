/*
以下是判断趋势反转最有效的技术指标及其综合应用方法，结合K线形态与量化指标形成多维度验证体系：
一、核心反转指标

	MACD指标‌
		零轴上方死叉（DIF下穿DEA）为顶部反转信号
		零轴下方金叉（DIF上穿DEA）为底部反转信号
		柱状线与价格背离（顶背离/底背离）增强信号可靠性‌

	KDJ+RSI组合‌
		KDJ超买区（J值>90）+RSI>70时顶部反转概率增大
		KDJ超卖区（J值<10）+RSI<30时底部反转概率增大
		两者同步背离时信号强度提升3倍‌

	‌移动平均线系统‌
		5日线上穿20日线形成黄金交叉（上涨反转）
		5日线下穿20日线形成死亡交叉（下跌反转）
		均线斜率变化比交叉点更早预示趋势衰竭‌
二、K线形态验证
	‌经典反转形态‌
		黄昏之星（上升趋势末端三根K线组合）
		早晨之星（下跌趋势末端三根K线组合）
		吞没形态（阳包阴/阴包阳）‌

	‌复合形态‌
		头肩顶/底（需突破颈线确认）
		双重顶/底（需成交量配合）
		岛形反转（缺口+跳空确认）‌

三、辅助验证指标
	‌成交量分析‌
		顶部反转时成交量萎缩
		底部反转时成交量放大
		量价背离（价格新高但成交量递减）‌

	‌市场宽度指标‌
		连涨连跌占比（超过80%个股同步上涨时警惕顶部）
		涨跌停家数比（极端值预示反转）‌

	‌资金流向指标‌
		融资买入比（异常放大预示顶部）
		主力资金净流入（持续背离价格走势）‌

四、实战应用策略
	多周期验证‌
		周线定方向（趋势级别）
		日线找买点（具体入场位）
		60分钟线确认信号‌

	‌指标组合优先级‌
		形态信号 > MACD交叉 > 均线系统 > KDJ/RSI
		当三个及以上指标同步发出信号时，反转成功率可达78%

	‌风险控制要点‌
		设置3%止损位（针对假突破）
		避免在消息面真空期操作
		极端行情中指标可能失效（如2020年原油负值事件）‌
	注：实际交易中建议结合基本面分析，技术指标仅提供概率性参考。历史数据显示，综合使用3种以上指标可将误判率降低42%
	对于加密货币的指标计算不能与股票的一样，周期更适合短一些，如：原来14天周期的最好改成7天周期，因为行情变换更快
*/

package indicators

import (
	"github.com/markcheno/go-talib"
	"github.com/prtmon/finance/common"
)

// 扩展指标计算函数
func CalculateAllIndicators(series common.Candlesticks) map[string]interface{} {
	results := make(map[string]interface{})
	ohlcv := series.ToOhlcv()

	//时间轴
	results["TIME"] = ohlcv.Time
	// 1. 移动平均线(MA)
	results["MA5"] = talib.Ma(ohlcv.Close, 5, talib.EMA)
	results["MA10"] = talib.Ma(ohlcv.Close, 10, talib.EMA)
	results["MA20"] = talib.Ma(ohlcv.Close, 20, talib.EMA)
	results["MA30"] = talib.Ma(ohlcv.Close, 30, talib.EMA)
	results["MA60"] = talib.Ma(ohlcv.Close, 60, talib.EMA)
	// 2. 指数移动平均线(EMA)
	results["EMA12"] = talib.Ema(ohlcv.Close, 12)
	results["EMA26"] = talib.Ema(ohlcv.Close, 26)

	// 3. 相对强弱指数(RSI)
	results["RSI14"] = talib.Rsi(ohlcv.Close, 14)

	// 4. MACD指标
	macd, macdSignal, _ := talib.Macd(ohlcv.Close, 12, 26, 9)
	results["MACD"] = macd
	results["MACD_Signal"] = macdSignal

	// 5. 布林带
	upper, middle, lower := talib.BBands(ohlcv.Close, 20, 2, 2, talib.EMA)
	results["BB_Upper"] = upper
	results["BB_Middle"] = middle
	results["BB_Lower"] = lower

	// 6. 随机指标(KDJ)
	slowK, slowD, slowJ := common.KDJ(ohlcv, 14, 3, talib.EMA, 3, talib.EMA)
	results["KDJ_K"] = slowK
	results["KDJ_D"] = slowD
	results["KDJ_J"] = slowJ

	// 7. 平均真实波幅(ATR)
	results["ATR14"] = talib.Atr(ohlcv.High, ohlcv.Low, ohlcv.Close, 14)

	// 8. 商品通道指数(CCI)
	results["CCI20"] = talib.Cci(ohlcv.High, ohlcv.Low, ohlcv.Close, 20)

	// 9. 能量潮指标(OBV)
	results["OBV"] = talib.Obv(ohlcv.Close, ohlcv.Volume)

	// 10. 平均趋向指数(ADX)
	results["ADX14"] = talib.Adx(ohlcv.High, ohlcv.Low, ohlcv.Close, 14)

	// 11. 威廉指标(WILLR)
	results["WILLR14"] = talib.WillR(ohlcv.High, ohlcv.Low, ohlcv.Close, 14)

	// 12. 资金流量指标(MFI)
	results["MFI14"] = talib.Mfi(ohlcv.High, ohlcv.Low, ohlcv.Close, ohlcv.Volume, 14)

	// 13. 抛物线转向指标(SAR)
	results["SAR"] = talib.Sar(ohlcv.High, ohlcv.Low, 0.02, 0.2)

	// 14. 三重指数平滑平均线(TRIX)
	results["TRIX18"] = talib.Trix(ohlcv.Close, 18)

	// 15. 变动率指标(ROC)
	results["ROC10"] = talib.Roc(ohlcv.Close, 10)

	//16.吞没形态
	results["ENGULFING"] = DetectEngulfing(ohlcv)

	//17.锤头线识别
	results["HAMMER"] = DetectHammerSignals(series)

	//18.早晨之星形态
	results["MORNING_STAR"] = DetectStarReversal(series)

	//19.RSI强弱指数-超买超卖
	results["OverTrade"] = RsiKdjOverTrade(ohlcv, 14, 30, 70, 20, 80)

	//...其他指标

	return results
}
