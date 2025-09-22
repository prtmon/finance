package indicators

import "github.com/prtmon/finance/common"

// ThreeWhiteSoldier 红三兵识别 (Three White Soldiers)
func ThreeWhiteSoldier(series common.Candlesticks) []bool {
	output := make([]bool, len(series))
	if len(series) < 3 {
		return output
	}

	for i := 2; i < len(series); i++ {
		c1 := series[i-2]
		c2 := series[i-1]
		c3 := series[i]

		// 1. 连续三根阳线 (收盘 > 开盘)
		condition1 := (c1.Close > c1.Open) && (c2.Close > c2.Open) && (c3.Close > c3.Open)

		// 2. 每根K线收盘价高于前一根收盘价
		condition2 := c2.Close > c1.Close && c3.Close > c2.Close

		// 3. 每根K线开盘价在前一根实体范围内
		condition3 := c2.Open > c1.Open && c2.Open < c1.Close &&
			c3.Open > c2.Open && c3.Open < c2.Close

		// 4. 实体长度递增 (可选)
		condition4 := (c2.Close-c2.Open) > (c1.Close-c1.Open) &&
			(c3.Close-c3.Open) > (c2.Close-c2.Open)

		output[i] = condition1 && condition2 && condition3 && condition4
	}
	return output
}
