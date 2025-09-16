package util

import (
	"fmt"
	"math"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

// AppPath 获取当前应用程序所在的目录路径, 优先通过环境变量APP_PATH读取
// 函数会尝试获取可执行文件的真实路径（解析符号链接），
// 若解析失败则返回原始可执行文件所在目录
// 返回值:
//   - string: 应用程序所在目录路径
//   - panic: 获取可执行文件路径时发生的错误（符号链接解析错误会被忽略）
func AppPath() string {
	if path := os.Getenv("APP_PATH"); path != "" {
		return path
	}
	// 获取当前可执行文件路径
	ext, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("获取可执行文件路径失败: %v", err))
	}

	// 尝试解析符号链接获取真实路径
	// 若解析失败，使用原始路径的目录作为回退方案
	realPath, err := filepath.EvalSymlinks(ext)
	if err != nil {
		return filepath.Dir(ext)
	}
	return filepath.Dir(realPath)
}

func LogPath() string {
	if path := os.Getenv("LOG_PATH"); path != "" {
		return path
	}
	// 获取当前可执行文件路径
	ext, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("获取可执行文件路径失败: %v", err))
	}

	// 尝试解析符号链接获取真实路径
	// 若解析失败，使用原始路径的目录作为回退方案
	realPath, err := filepath.EvalSymlinks(ext)
	if err != nil {
		return filepath.Dir(ext)
	}
	return filepath.Dir(realPath)
}

// 获取本机所有ip
func GetAllLocalIPs() ([]string, error) {
	ipArr := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ipArr, err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipArr = append(ipArr, ipnet.IP.String())
			}
		}

	}
	return ipArr, nil
}

func DelayMicroseconds(us int64) {
	var tv syscall.Timeval
	_ = syscall.Gettimeofday(&tv)

	stratTick := int64(tv.Sec)*int64(1000000) + int64(tv.Usec) + us
	endTick := int64(0)
	for endTick < stratTick {
		_ = syscall.Gettimeofday(&tv)
		endTick = int64(tv.Sec)*int64(1000000) + int64(tv.Usec)
	}
}

func GetAndEqJionString(payload map[string]string) (formattedString string) {

	if len(payload) == 0 {
		return
	}

	keys := make([]string, 0, len(payload))
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sortedDict := make([]string, len(keys))
	for i, k := range keys {
		sortedDict[i] = k + "=" + payload[k]
	}
	formattedString = strings.Join(sortedDict, "&")
	return formattedString
}

// 一维数组转二维数组 size为二维数组中一维数组的长度
func Arr2Arr[T any](arr []T, size int) [][]T {

	if len(arr) < size {
		return [][]T{arr}
	}
	length := len(arr) / size
	slices := make([][]T, length+1)
	i := 0
	for {
		if len(arr)-size < 1 {
			slices[i] = arr
			break
		}
		slices[i] = arr[:size]
		i++
		arr = arr[size:]
	}
	return slices
}

func Bytes2String(b []byte) string {
	return string(b)
}

func Bytes2Float64(b []byte) float64 {
	return BytesToFloat64(b)
}

func Int64toFloat64WithExponent(v int64, exponent int) float64 {
	return float64(v) * math.Pow10(exponent)
}

// Int64toFloat64WithExponent 将int64数值转换为浮点数，并对最后一位（由exponent决定的小数位）四舍五入
// v: 原始整数值（如价格的整数表示）
// exponent: 指数（通常为负数，表示小数点后位数，如exponent=-2表示保留2位小数）
func RoundInt64toFloat64WithExponent(v int64, exponent int) float64 {
	// 计算需要保留的小数位数（d = -exponent）
	d := -exponent
	if d < 0 {
		d = 0 // 避免指数为正时出现负的小数位数
	}
	scale := math.Pow10(d) // 缩放因子（用于定位四舍五入的位置）

	// 原始转换值（未四舍五入）
	rawValue := float64(v) * math.Pow10(exponent)

	// 四舍五入到指定小数位
	return math.Round(rawValue*scale) / scale
}
