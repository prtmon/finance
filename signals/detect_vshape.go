// 代码由DeepSeek给出,结合本地kline更改相关定义

package signals

import (
	"errors"
	"github.com/prtmon/finance/common"
	"github.com/prtmon/finance/utils"
)

// VShapeSignal 检测结果结构体
type VShapeSignal struct {
	Date             int64
	IsValidV         bool    // 是否有效V形反转
	DropRatio        float64 // 累计跌幅
	ConfirmationDays int     // 趋势确认天数
}

// findMinClose 查找窗口期内的最低收盘价
// 参数：
//
//	klines: K线数据切片（需保证至少包含一个元素）
//
// 返回值：
//
//	窗口期内的最低收盘价
func findMinClose(klines common.Candlesticks) float64 {
	if len(klines) == 0 {
		return 0 // 返回0表示无效值（调用方需确保输入非空）
	}

	minVal := klines[0].Close // 初始化最小值为第一个元素的收盘价
	for _, k := range klines {
		if k.Close < minVal {
			minVal = k.Close // 更新最小值
		}
	}
	return minVal
}

// findMinCloseOptimized 滑动窗口优化版本
func findMinCloseOptimized(klines common.Candlesticks, window int) []float64 {
	if len(klines) == 0 || window <= 0 {
		return nil
	}

	result := make([]float64, 0, len(klines)-window+1)
	deque := make([]int, 0) // 双端队列，存储索引

	for i := 0; i < len(klines); i++ {
		// 移除队列中超出窗口范围的索引
		if len(deque) > 0 && deque[0] == i-window {
			deque = deque[1:]
		}

		// 移除队列中比当前元素大的索引
		for len(deque) > 0 && klines[deque[len(deque)-1]].Close >= klines[i].Close {
			deque = deque[:len(deque)-1]
		}

		// 将当前索引加入队列
		deque = append(deque, i)

		// 记录当前窗口的最小值
		if i >= window-1 {
			result = append(result, klines[deque[0]].Close)
		}
	}

	return result
}

// calculateVolumeMA 计算窗口期内的成交量移动平均值
// 参数：
//
//	klines: K线数据切片（需保证至少包含一个元素）
//	window: 计算均值的窗口大小
//
// 返回值：
//
//	窗口期内的成交量均值
func calculateVolumeMA(klines common.Candlesticks, window int) float64 {
	if len(klines) == 0 || window <= 0 {
		return 0 // 返回0表示无效值（调用方需确保输入合法）
	}

	sum := 0.0
	for _, k := range klines {
		sum += k.Volume // 累加窗口期内的成交量
	}
	return sum / float64(window) // 计算均值
}

// DetectVShape 增强版V形反转检测
//
//	confirmationDays: 趋势确认天数（默认1）
//	confirmationThreshold: 确认涨幅阈值（0.01表示1%）
func DetectVShape(klines common.Candlesticks, window int, dropThreshold, reboundThreshold, volumeRatio float64, confirmationDays int, confirmationThreshold float64) ([]VShapeSignal, error) {
	if len(klines) < window+confirmationDays+1 {
		return nil, errors.New("insufficient data length")
	}

	var signals []VShapeSignal

	for i := window; i < len(klines)-confirmationDays; i++ {
		// 1. 检查严格连续下跌
		strictDrop := true
		for j := i - window; j < i; j++ {
			if klines[j+1].Close >= klines[j].Close {
				strictDrop = false
				break
			}
		}

		// 2. 计算累计跌幅
		startPrice := klines[i-window].Close
		lowestPrice := findMinClose(klines[i-window : i])
		cumulativeDrop := (startPrice - lowestPrice) / startPrice

		// 3. 成交量验证
		volumeMA := calculateVolumeMA(klines[i-window:i-1], window-1)
		volumeSpike := klines[i].Volume > volumeRatio*volumeMA

		// 4. 反弹条件
		rebound := (klines[i].Close - klines[i-1].Close) / klines[i-1].Close

		// 5. 趋势确认（新增核心逻辑）
		confirmed := true
		for k := 1; k <= confirmationDays; k++ {
			// 检查后续交易日是否持续上涨
			if (klines[i+k].Close-klines[i+k-1].Close)/klines[i+k-1].Close < confirmationThreshold {
				confirmed = false
				break
			}
		}

		// 综合判断
		isValidV := strictDrop &&
			cumulativeDrop >= dropThreshold &&
			rebound > reboundThreshold &&
			volumeSpike &&
			confirmed

		signals = append(signals, VShapeSignal{
			Date:             utils.ToInt64(klines[i].Time),
			IsValidV:         isValidV,
			DropRatio:        cumulativeDrop,
			ConfirmationDays: confirmationDays,
		})
	}

	return signals, nil
}

// 辅助函数保持不变...

/*
改进说明：
1. 新增趋势确认机制：
   - 要求反弹后连续N日保持上涨趋势
   - 可配置确认天数（confirmationDays）
   - 可配置每日最小涨幅（confirmationThreshold）

2. 严格连续下跌判断：
   - 维持原有严格递减检查
   - 确保窗口期内无任何反弹

3. 信号有效性验证：
   - 原始V形条件 + 趋势确认 = 有效信号
   - 避免假突破后的再次下跌

算法流程图：
连续下跌检测 → 跌幅达标 → 反弹验证 → 成交量验证 → 趋势确认 → 有效信号

使用示例：
klines := []KLine{
	{"D1", 10.0, 1000},  // 下跌开始
	{"D2", 9.5, 1200},
	{"D3", 9.0, 1100},  // 窗口期结束
	{"D4", 9.3, 3000},  // 反弹日
	{"D5", 9.6, 2800},  // 确认日1（+3.2%）
	{"D6", 9.8, 2500},  // 确认日2（+2.1%）
}

// 要求2天确认，每日涨幅>1%
signals, _ := DetectVShape(klines, 3, 0.05, 0.03, 1.5, 2, 0.01)
*/

// 辅助函数保持不变...

// 建议参数组合（适用于A股市场）
/*signals, _ := DetectVShape(klines,
3,    // 3日观察窗口
0.08, // 8%累计跌幅
0.04, // 4%单日反弹
1.8,  // 1.8倍成交量
2,    // 2日趋势确认
0.015 // 每日1.5%涨幅阈值
)
*/
// 建议组合结果验证：
// 1. 3日严格连续下跌（-8%以上）
// 2. 放量反弹（成交量>1.8倍均量）
// 3. 后续2日持续上涨（每日>1.5%）
