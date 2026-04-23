package util

import (
	"math"
)

func FloatCeil(val float64, n int) float64 {
	if n == 0 {
		return float64(int64(val))
	}
	p := math.Pow10(n)
	s := math.Pow10((n + 5) * -1) // Adjustment for floating-point precision issues
	return math.Ceil(val*p-s) / p
}

// Float64EqualWithEpsilon 判断两个float64是否近似相等（精度要求不高时）
// epsilon: 允许的误差范围（如1e-6表示百万分之一精度）
func Float64EqualWithEpsilon(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// fixme !!! 必须确保所有用到浮点数精度的地方的精度都大于这个最小的误差值
const float64Epsilon = 1e-9

func Float64Equal(a, b float64) bool {
	return math.Abs(a-b) < float64Epsilon
}

func Float64Greater(a, b float64) bool {
	return a-b > float64Epsilon
}

// Float64Less 判断 a 是否严格小于 b（考虑精度误差）
// 当 a < b - epsilon 时，认为 a < b
func Float64Less(a, b float64) bool {
	return b-a > float64Epsilon
}

// QtyFilter 策略层专用：在计算仓位前，先按数量步进和精度进行过滤
// 这里的 step 和 precision 是该币种/市场的原生属性
