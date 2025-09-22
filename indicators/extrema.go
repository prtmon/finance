package indicators

import "github.com/prtmon/finance/common"

// DetectExtrema 滑动窗口极值检测
func DetectExtrema(candles common.Candlesticks, window int) []int64 {
	n := len(candles)
	output := make([]int64, n)
	for i := window; i < n-window; i++ {
		current := candles[i]
		isPeak, isTrough := true, true

		// 检查窗口范围内极值
		for j := i - window; j <= i+window; j++ {
			if j == i {
				continue
			}
			if current.High <= candles[j].High {
				isPeak = false
			}
			if current.Low >= candles[j].Low {
				isTrough = false
			}
		}
		if isPeak {
			//峰顶
			output[i] = 1
		}
		if isTrough {
			//谷底
			output[i] = -1
		}
	}
	return output
}
