package util

import (
	"math"
	"strconv"
	"strings"
)

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
func FixPreFloat64ToString(f float64, prec int) string {
	return strconv.FormatFloat(FloatFixPre(f, prec), 'f', -1, 64)
}

func FixFormatFloat64(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func IsZeroWithPrec(f float64, prec int) bool {
	return FloatFixPre(f, prec) == 0
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToFloat64(num string) float64 {
	fnum, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0
	}
	return fnum
}

func BytesToFloat64(buf []byte) float64 {
	fnum, err := strconv.ParseFloat(UnsafeBytesToString(buf), 64)
	if err != nil {
		return 0
	}
	return fnum
}

// 获取有效位数（将"0.001"转换为3）
// "0.001" => 3
// "0.01" => 2
// "1" => 0
// "100" => -2
func GetDecimalPlaces(str string) int {

	// 如果字符串包含小数点
	if strings.Contains(str, ".") {
		for strings.HasSuffix(str, "0") {
			str = str[:len(str)-1]
		}

		// 分割字符串，获取小数点后的部分
		parts := strings.Split(str, ".")
		// 返回小数点后的位数
		return len(parts[1])

	} else {
		// 如果字符串不包含小数点，计算末尾的零的数量
		zeroCount := 0
		for i := len(str) - 1; i >= 0; i-- {
			if str[i] == '0' {
				zeroCount++
			} else {
				break
			}
		}
		// 返回零的数量的负值
		return -zeroCount
	}
}

// 根据有效位数获取最小步进（将3转换为0.001）
func GetMinStep(decimalPlaces int32) float64 {
	return math.Pow(10, float64(-decimalPlaces))
}

// 返回指定精度的浮点数，不四舍五入
func FloatFixPre(val float64, n int) float64 {
	if n == 0 {
		return float64(int64(val))
	}
	p := math.Pow10(n)
	s := math.Pow10((n + 5) * -1) //后面加s的原因是有的时候浮点数的精度问题，比如0.3,保存时可能是以0.2999999999999998来保存的
	return math.Floor(val*p+s) / p
}

// 返回指定精度的浮点数，四舍五入
func RoundFloat(val float64, n int) float64 {
	p := math.Pow10(n)
	s := math.Pow10((n + 5) * -1) //后面加s的原因是有的时候浮点数的精度问题，比如0.3,保存时可能是以0.2999999999999998来保存的
	return math.Round(val*p+s) / p
}

// 调整有效位数（舍去）
// value: 312.12345999999998, decimalPlaces: -2 => 300
// value: 312.12345999999998, decimalPlaces: -1 => 310
// value: 312.12345999999998, decimalPlaces: 0 => 312
// value: 312.12345999999998, decimalPlaces: 1 => 312.1
// value: 312.12345999999998, decimalPlaces: 2 => 312.12
// value: 312.12345999999998, decimalPlaces: 5 => 312.12346
func FormatFloat(value float64, decimalPlaces int32) float64 {
	str := strconv.FormatFloat(value, 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += "."
	}

	if str == "NaN" || str == "Inf" || str == "-Inf" {
		return 0
	}

	if len(str) > 18 {
		str = str[:18]
	}
	padding := strings.Repeat("0", 18-len(str))
	formatted := str + padding

	// 计算小数点部分的位数
	parts := strings.Split(formatted, ".")
	floatLen := 0
	if len(parts) > 1 {
		floatLen = len(parts[1])
	}

	value += math.Pow10(-floatLen + 1)

	str = strconv.FormatFloat(value, 'f', -1, 64)
	parts = strings.Split(str, ".")
	intStr := parts[0]
	floatStr := parts[1]
	floatStr += strings.Repeat("0", 18)

	if decimalPlaces > 0 {
		v, _ := strconv.ParseFloat(intStr+"."+floatStr[:decimalPlaces], 64)
		return v
	} else if decimalPlaces < 0 {

		if len(intStr) <= -int(decimalPlaces) {
			return 0
		}

		intStr2 := intStr[:len(intStr)+int(decimalPlaces)]
		v, _ := strconv.ParseFloat(intStr2+strings.Repeat("0", len(intStr)-len(intStr2)), 64)
		return v
	} else {
		v, _ := strconv.ParseFloat(intStr, 64)
		return v
	}
}

func GetMaxFloat64(v ...float64) float64 {
	max := v[0]
	for _, val := range v {
		if max < val {
			max = val
		}
	}
	return max
}

func GetMinFloat64(v ...float64) float64 {
	min := v[0]
	for _, val := range v {
		if min > val {
			min = val
		}
	}
	return min
}
func Int64toStringWithExponent(v int64, exponent int) string {
	return strconv.FormatFloat(float64(v)*math.Pow10(exponent), 'f', -1, 64)
}

func BytesToInt64(buf []byte) int64 {
	fnum, err := strconv.ParseInt(UnsafeBytesToString(buf), 10, 64)
	if err != nil {
		return 0
	}
	return fnum
}
