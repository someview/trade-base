package tracer

import (
	"sync"
)

type BucketIndexFunc func(latencyUs int64) int

// HistogramStats 直方图统计结构体
type HistogramStats struct {
	buckets         []int64
	bucketIndexFunc BucketIndexFunc
	slowThresholdUs int64
	slowRequests    int64
	totalRequests   int64
	totalLatencyUs  int64 // 添加总延迟时间用于计算平均值
	minLatency      int64
	maxLatency      int64
	mu              sync.RWMutex
}

// NewHistogramStats 创建新的直方图统计实例
func NewHistogramStats(bucketCount int, slowThresholdUs int64, bucketIndexFunc func(latencyUs int64) int) *HistogramStats {
	if slowThresholdUs <= 0 {
		slowThresholdUs = 600 // 默认600微秒
	}
	if bucketIndexFunc == nil {
		bucketIndexFunc = DefaultBucketIndexFunc
	}
	return &HistogramStats{
		buckets:         make([]int64, bucketCount),
		slowThresholdUs: slowThresholdUs,
		bucketIndexFunc: bucketIndexFunc,
	}
}

// DefaultBucketIndexFunc 默认桶索引计算函数
// 桶配置：<600µs(桶0), 600-800µs(桶1), 800-1000µs(桶2), ..., 1800-2000µs(桶7), >2000µs(桶8)
func DefaultBucketIndexFunc(latencyUs int64) int {
	if latencyUs < 600 {
		return 0
	}
	if latencyUs > 2000 {
		return 8
	}
	// 600-2000范围内，每200微秒一个桶
	return int((latencyUs-600)/200) + 1
}

// RecordLatency 记录延迟数据（微秒）
func (h *HistogramStats) RecordLatency(latencyUs int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.totalRequests++
	h.totalLatencyUs += latencyUs // 累计总延迟

	h.minLatency = min(h.minLatency, latencyUs)
	h.maxLatency = max(h.maxLatency, latencyUs)

	// 记录到直方图桶中
	bucketIndex := h.bucketIndexFunc(latencyUs)
	if bucketIndex >= 0 && bucketIndex < len(h.buckets) {
		h.buckets[bucketIndex]++
	}

	// 只统计慢请求
	if latencyUs >= h.slowThresholdUs {
		h.slowRequests++
	}
}

// GetStats 获取统计信息
func (h *HistogramStats) GetStats() HistogramStatsResult {
	h.mu.RLock()
	defer h.mu.RUnlock()

	buckets := make([]int64, len(h.buckets))
	copy(buckets, h.buckets)

	result := HistogramStatsResult{
		Buckets:       buckets,
		SlowThreshold: h.slowThresholdUs,
		SlowCount:     h.slowRequests,
		TotalRequests: h.totalRequests,
		MinLatency:    h.minLatency,
		MaxLatency:    h.maxLatency,
	}

	if h.totalRequests > 0 {
		result.AvgLatency = h.totalLatencyUs / h.totalRequests
		result.SlowRatio = float64(h.slowRequests) / float64(h.totalRequests)
	}

	return result
}

// GetStatsAndReset 获取统计信息快照并重置数据
func (h *HistogramStats) GetStatsAndReset() HistogramStatsResult {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 获取当前统计信息
	buckets := make([]int64, len(h.buckets))
	copy(buckets, h.buckets)

	result := HistogramStatsResult{
		Buckets:       buckets,
		SlowThreshold: h.slowThresholdUs,
		SlowCount:     h.slowRequests,
		TotalRequests: h.totalRequests,
		MinLatency:    h.minLatency,
		MaxLatency:    h.maxLatency,
	}

	if h.totalRequests > 0 {
		result.AvgLatency = h.totalLatencyUs / h.totalRequests
		result.SlowRatio = float64(h.slowRequests) / float64(h.totalRequests)
	}

	// 重置所有统计数据
	for i := range h.buckets {
		h.buckets[i] = 0
	}
	h.slowRequests = 0
	h.totalRequests = 0
	h.totalLatencyUs = 0
	h.minLatency = 0
	h.maxLatency = 0

	return result
}

// Reset 重置统计数据
func (h *HistogramStats) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := range h.buckets {
		h.buckets[i] = 0
	}
	h.slowRequests = 0
	h.totalRequests = 0
	h.totalLatencyUs = 0
	h.minLatency = 0
	h.maxLatency = 0
}

// HistogramStatsResult 统计结果结构体
type HistogramStatsResult struct {
	Buckets       []int64 `json:"buckets"`        // 各桶的计数
	SlowThreshold int64   `json:"slow_threshold"` // 慢请求阈值（微秒）
	SlowCount     int64   `json:"slow_count"`     // 慢请求计数
	SlowRatio     float64 `json:"slow_ratio"`     // 慢请求比例
	TotalRequests int64   `json:"total_requests"` // 总请求数
	MinLatency    int64   `json:"min_latency"`    // 最小延迟（微秒）
	MaxLatency    int64   `json:"max_latency"`    // 最大延迟（微秒）
	AvgLatency    int64   `json:"avg_latency"`    // 平均延迟（微秒）
}
