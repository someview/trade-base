package tracer

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/someview/trade-base/global"
)

// 仅包含字母数字的字符集，便于日志和正则处理
var charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
var counter uint32
var baseTimeMicros = time.Date(global.StartTime().Year(), 1, 1, 0, 0,
	0, 0, global.StartTime().Location()).UnixMicro() // 使用微秒基准时间

func init() {
	counter = 0
}

// NewTraceId 生成唯一的跟踪ID,
//
//	同一台机器,同一年内,同一微妙内，只有同时超过4095个traceId,才会出现重复的情况
//
// 2年时间范围内的计数不会超过10位数
func NewTraceId() string {
	// 使用微秒而不是纳秒，仍然有很高精度但数值小得多
	nowMicros := time.Now().UTC().UnixMicro() - baseTimeMicros

	// 获取计数器值
	count := atomic.AddUint32(&counter, 1) & 0xFFF // 12位计数器

	// 组合值
	uniqueValue := (nowMicros << 12) | int64(count)

	// 编码为Base62，不进行反转
	result := toBase62(uniqueValue)

	if len(result) > 10 {
		result = result[:10]
	}

	return result
}

// toBase62 将数值编码为Base62，使用strings.Builder避免内存分配
func toBase62(num int64) string {
	if num < 0 {
		num = -num
	}

	if num == 0 {
		return string(charset[0])
	}

	// 预先分配足够的容量，避免动态扩容
	var builder strings.Builder
	builder.Grow(12) // 预分配12个字符空间，足够容纳大多数ID

	for num > 0 {
		builder.WriteByte(charset[num%62])
		num /= 62
	}

	return builder.String()
}
