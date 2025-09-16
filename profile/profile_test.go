package profile

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
	"time"
)

// TestSystemMonitor 测试系统监控功能
func TestSystemMonitor(t *testing.T) {
	// 创建监控配置，设置较低的阈值以便于测试
	cfg := &MonitorConfig{
		CPUThreshold:         100, // 设置较低的CPU阈值
		MemThreshold:         500, // 设置较低的内存阈值（MB）
		DiskPercentThreshold: 10,  // 设置较低的硬盘阈值（GB）
		GCDurationUs:         100, // 设置较低的GC时长阈值（us）
		CheckIntervals:       5,   // 缩短检查间隔以便快速测试
	}

	// 创建监控器
	monitor := NewSystemMonitor(cfg)

	// 导出checkAndAlert方法供测试使用
	t.Log("开始手动触发系统检查...")
	monitor.checkSysStatus()

	// 创建一些内存压力
	t.Log("开始创建内存压力...")
	var data [][]byte
	for i := 0; i < 5; i++ {
		// 每次分配约100MB内存
		chunk := make([]byte, 100*1024*1024)
		for j := range chunk {
			chunk[j] = byte(j % 256)
		}
		data = append(data, chunk)
		t.Logf("已分配约 %d MB内存", (i+1)*100)
	}

	// 再次检查系统状态
	t.Log("再次触发系统检查...")
	monitor.checkSysStatus()

	// 模拟GC
	t.Log("触发GC...")
	for i := 0; i < 3; i++ {
		data = nil // 释放内存，触发GC
		t.Log("已释放内存，等待GC...")
		time.Sleep(time.Second)
	}

	// 最后检查一次
	t.Log("最终系统检查...")
	monitor.checkSysStatus()

	t.Log("测试完成")
}

// TestMonitorStart 测试监控器的启动和自动检查
func TestMonitorStart(t *testing.T) {
	// 使用更短的检查间隔
	cfg := &MonitorConfig{
		CheckIntervals: 2,
	}

	monitor := NewSystemMonitor(cfg)

	// 启动监控，但仅运行短时间
	go monitor.Start()

	// 运行一段时间后关闭
	t.Log("监控器已启动，等待自动检查...")
	time.Sleep(5 * time.Second)

	monitor.Close()
	t.Log("监控器已关闭")
}

// TestMemoryThreshold 专门测试内存阈值告警
func TestMemoryThreshold(t *testing.T) {
	// 获取当前内存使用情况
	memUsage, _, err := checkMemory()
	if err != nil {
		t.Fatalf("获取内存使用量失败: %v", err)
	}

	// 设置略低于当前内存使用量的阈值，确保触发告警
	cfg := &MonitorConfig{
		MemThreshold: memUsage - 100, // 设置比当前使用量低100MB的阈值
	}

	monitor := NewSystemMonitor(cfg)
	t.Logf("当前内存使用量: %.2fMB, 设置阈值: %.2fMB", memUsage, cfg.MemThreshold)

	// 触发检查，应该会产生告警
	t.Log("触发检查，预期产生内存告警...")
	alerts := monitor.checkSysStatus()
	slog.Info("检查结果:", slog.Any("alerts", alerts))
}

func TestCPUThreshold(t *testing.T) {
	// 获取当前内存使用情况
	cpuUsage, _, err := checkCPU()
	if err != nil {
		t.Fatalf("获取内存使用量失败: %v", err)
	}

	// 设置略低于当前内存使用量的阈值，确保触发告警
	cfg := &MonitorConfig{
		CPUThreshold: cpuUsage - 0.4, // 设置比当前使用量低100MB的阈值
	}

	monitor := NewSystemMonitor(cfg)
	t.Logf("当前CPU使用量: %.2f, 设置阈值: %.2f", cpuUsage, cfg.MemThreshold)

	// 触发检查，应该会产生告警
	t.Log("触发检查，预期产生内存告警...")
	alerts := monitor.checkSysStatus()
	slog.Info("检查结果:", slog.Any("alerts", alerts))
}

func TestCheckDisk(t *testing.T) {
	// 获取当前内存使用情况
	diskUsage, err := checkDiskPercent("/")
	assert.Nil(t, err)
	slog.Info("当前硬盘使用量:", slog.Float64("使用量", diskUsage))
}
