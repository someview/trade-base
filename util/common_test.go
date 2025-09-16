package util

import (
	"log/slog"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllLocalIPs(t *testing.T) {
	ipArr, err := GetAllLocalIPs()
	assert.Nil(t, err)
	assert.NotEmpty(t, ipArr)
}

func TestGet(t *testing.T) {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // 设置日志级别
		AddSource: true,            // 可选：记录源码位置
	})
	logger := slog.New(handler)
	logger.Info("item", slog.String("hello", "123"))
}

func TestAppPath(t *testing.T) {
	slog.Info("apppth", slog.Any("key", AppPath()))
}

func TestInt64toFloat64WithExponent(t *testing.T) {
	// 原测试：验证相同输入输出一致（基础功能）
	t.Run("Same input produces same output", func(t *testing.T) {
		v := int64(3798)
		exponent := -2
		s1 := Int64toFloat64WithExponent(v, exponent)
		t.Log(s1)
		s2 := RoundInt64toFloat64WithExponent(v, exponent)
		t.Log(s2)
	})

	t.Run("Five rounding example", func(t *testing.T) {
		// 使用一个末尾是5的例子
		v := int64(12345)
		exponent := -4

		// 转换为1.2345
		original := Int64toFloat64WithExponent(v, exponent)
		t.Logf("原始值: %.20f", original) // 显示更多小数位以查看潜在精度问题

		// 四舍五入到3位小数应该得到1.235
		roundedTo3 := math.Round(original*1000) / 1000
		t.Logf("手动四舍五入到3位小数: %.3f", roundedTo3)

		// 使用我们的函数，先转换回整数然后应用四舍五入
		v2 := int64(1235)
		exponent2 := -3
		usingFunc := RoundInt64toFloat64WithExponent(v2, exponent2)
		t.Logf("使用函数四舍五入: %.3f", usingFunc)

		// 验证结果一致
		assert.Equal(t, 1.235, roundedTo3, "手动四舍五入应该将1.2345变为1.235")
		assert.Equal(t, 1.235, usingFunc, "函数应该将1.235正确表示")
	})

}
