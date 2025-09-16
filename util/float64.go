// Package util 提供了处理浮点数精度和格式化的工具函数
package util

import (
	"math"
	"strconv"
)

// floatEpsilon 是一个极小值，用于修正 float64 在计算中产生的微小二進制精度误差。
const floatEpsilon = 1e-9

// todo 详尽的测试
// truncateFloat (辅助函数)
// 一个健壮的、使用微调法进行指定小数位截断（不四舍五入）的函数。
func truncateFloat(value float64, precision int) float64 {
	if precision == 0 {
		return math.Floor(value + floatEpsilon)
	}
	p := math.Pow10(precision)
	// 关键：乘以p后，加上微调值，再Floor，最后除以p，完成精确截断
	return math.Floor(value*p+floatEpsilon) / p
}

// --- 核心计算函数 (已重构，更清晰) ---

// FormatFloatUpByMultiplier 将 Float 向上舍入到最接近的 Multiplier 的倍数。
func FormatFloatUpByMultiplier(Float float64, FloatMultiplier float64, FloatPrecision int) float64 {
	if FloatMultiplier <= 0 {
		return truncateFloat(Float, FloatPrecision)
	}
	// 向上舍入时，应减去微调值，防止 x.000...001 这样的值被错误地进位
	adjustedFloat := math.Ceil(Float/FloatMultiplier-floatEpsilon) * FloatMultiplier
	return truncateFloat(adjustedFloat, FloatPrecision)
}

// FormatFloatDownByMultiplier 将 Float 向下舍入到最接近的 Multiplier 的倍数。
func FormatFloatDownByMultiplier(Float float64, FloatMultiplier float64, FloatPrecision int) float64 {
	if FloatMultiplier <= 0 {
		return truncateFloat(Float, FloatPrecision)
	}
	// 核心修正：在除法后加上微调值，再进行Floor操作
	adjustedFloat := math.Floor(Float/FloatMultiplier+floatEpsilon) * FloatMultiplier
	return truncateFloat(adjustedFloat, FloatPrecision)
}

// FormatFloatRoundByMultiplier 将 Float 四舍五入到最接近的 Multiplier 的倍数。
func FormatFloatRoundByMultiplier(Float float64, FloatMultiplier float64, FloatPrecision int) float64 {
	if FloatMultiplier <= 0 {
		return truncateFloat(Float, FloatPrecision)
	}
	adjustedFloat := math.Round(Float/FloatMultiplier) * FloatMultiplier
	return truncateFloat(adjustedFloat, FloatPrecision)
}

// --- 字符串输出封装函数 (保持不变) ---

func FormatFloatStrUpByMultiplier(Float float64, FloatMultiplier float64, FloatPrecision int) (float64, string) {
	res := FormatFloatUpByMultiplier(Float, FloatMultiplier, FloatPrecision)
	return res, strconv.FormatFloat(res, 'f', FloatPrecision, 64)
}

func FormatFloatStrDownByMultiplier(Float float64, FloatMultiplier float64, FloatPrecision int) (float64, string) {
	res := FormatFloatDownByMultiplier(Float, FloatMultiplier, FloatPrecision)
	return res, strconv.FormatFloat(res, 'f', FloatPrecision, 64)
}

func FormatFloatStrRoundByMultiplier(Float float64, FloatMultiplier float64, FloatPrecision int) (float64, string) {
	res := FormatFloatRoundByMultiplier(Float, FloatMultiplier, FloatPrecision)
	return res, strconv.FormatFloat(res, 'f', FloatPrecision, 64)
}

func QtyToSize(qty float64, contractValue float64) int64 {
	return int64(math.Round(qty/contractValue + floatEpsilon))
}

func SizeToQty(size int64, contractValue float64, qtyPrecision int) float64 {
	return truncateFloat(float64(size)*contractValue, qtyPrecision)
}
