package profile

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/process"
	"github.com/someview/trade-base/jcdtg"
	"github.com/someview/trade-base/log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// MonitorConfig 监控阈值配置（新增 DiskPath 字段）
const MonitorConfigKey = "monitor_config"

type MonitorConfig struct {
	CPUThreshold         float64 `mapstructure:"cpu_threshold"`          // CPU使用绝对值
	MemThreshold         float64 `mapstructure:"mem_threshold"`          // 内存使用绝对值
	DiskPercentThreshold float64 `mapstructure:"disk_percent_threshold"` // 硬盘使用绝对值
	DiskPath             string  `mapstructure:"disk_path"`              // 硬盘监控路径
	GCDurationUs         int64   `mapstructure:"gc_duration_us"`         // GC 时长阈值（微秒）
	CheckIntervals       int64   `mapstructure:"check_interval_s"`       // 检查间隔
}

func (m MonitorConfig) CheckValidity() error {
	return nil
}

// SystemInfo 系统硬件信息
type SystemInfo struct {
	CPUCores       int     // CPU核心数
	TotalMemoryMB  float64 // 总内存(MB)
	AllocatedMemMB float64 // 程序已分配内存(MB)
	SysMemoryMB    float64 // 程序从系统获取的内存(MB)
	TotalDiskGB    float64 // 总硬盘空间(GB)
	FreeDiskGB     float64 // 可用硬盘空间(GB)
	HostOS         string  // 操作系统信息
	GoVersion      string  // Go版本
}

// SystemMonitor 系统监控器
type SystemMonitor struct {
	cfg        *MonitorConfig
	ticker     *time.Ticker
	systemInfo *SystemInfo
}

// NewSystemMonitor 创建监控器实例（优雅化配置合并逻辑）
func NewSystemMonitor(cfg *MonitorConfig) *SystemMonitor {
	// 定义默认配置（附设计说明）
	defaultConfig := &MonitorConfig{
		CPUThreshold:         2,        // 默认值 200%
		MemThreshold:         2048 * 3, // 默认内存为2048M
		DiskPercentThreshold: 0.8,      // 默认使用10240M
		DiskPath:             "/",      // 默认监控根目录（覆盖系统主要存储）
		GCDurationUs:         300,      // 默认GC时长300us(超过300us就需要告警)
		CheckIntervals:       60,       // 默认检查间隔60秒（平衡实时性与性能开销）
	}

	// 初始化合并后的配置：优先使用用户配置，无则用默认
	if cfg.CPUThreshold == 0 {
		cfg.CPUThreshold = defaultConfig.CPUThreshold
	}
	if cfg.MemThreshold == 0 {
		cfg.MemThreshold = defaultConfig.MemThreshold
	}
	if cfg.DiskPercentThreshold == 0 {
		cfg.DiskPercentThreshold = defaultConfig.DiskPercentThreshold
	}
	if cfg.DiskPath == "" {
		cfg.DiskPath = defaultConfig.DiskPath
	}

	if cfg.GCDurationUs == 0 {
		cfg.GCDurationUs = defaultConfig.GCDurationUs
	}

	if cfg.CheckIntervals == 0 {
		cfg.CheckIntervals = defaultConfig.CheckIntervals
	}

	// 收集系统信息
	sysInfo := collectSystemInfo(cfg.DiskPath)

	// 日志记录系统配置信息（只记录日志，不发送TG）
	logSystemInfo(sysInfo, cfg)

	return &SystemMonitor{
		cfg:        cfg,
		systemInfo: sysInfo,
	}
}

// collectSystemInfo 收集系统硬件信息
func collectSystemInfo(diskPath string) *SystemInfo {
	info := &SystemInfo{
		GoVersion: runtime.Version(),
	}

	// 获取CPU信息
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.CPUCores = len(cpuInfo)
	} else {
		info.CPUCores = runtime.NumCPU()
	}

	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.TotalMemoryMB = float64(memInfo.Total) / (1024 * 1024)
	}

	// 获取程序内存信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	info.AllocatedMemMB = float64(memStats.Alloc) / (1024 * 1024)
	info.SysMemoryMB = float64(memStats.Sys) / (1024 * 1024)

	// 获取硬盘信息
	diskInfo, err := disk.Usage(diskPath)
	if err == nil {
		info.TotalDiskGB = float64(diskInfo.Total) / (1024 * 1024 * 1024)
		info.FreeDiskGB = float64(diskInfo.Free) / (1024 * 1024 * 1024)
	}

	// 获取操作系统信息
	hostInfo, err := host.Info()
	if err == nil {
		info.HostOS = fmt.Sprintf("%s %s (%s)", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.OS)
	}

	return info
}

// logSystemInfo 记录系统信息（只记录日志，不发送TG）
func logSystemInfo(sysInfo *SystemInfo, cfg *MonitorConfig) {
	record := log.NewRecord(log.InfoLevel, "系统硬件信息")
	record.AddAttr(log.Int("CPU核心数", sysInfo.CPUCores))
	record.AddAttr(log.Float64("总内存(MB)", sysInfo.TotalMemoryMB))
	record.AddAttr(log.Float64("程序已分配内存(MB)", sysInfo.AllocatedMemMB))
	record.AddAttr(log.Float64("程序系统内存(MB)", sysInfo.SysMemoryMB))
	record.AddAttr(log.Float64("总硬盘空间(GB)", sysInfo.TotalDiskGB))
	record.AddAttr(log.Float64("可用硬盘空间(GB)", sysInfo.FreeDiskGB))
	record.AddAttr(log.String("操作系统", sysInfo.HostOS))
	record.AddAttr(log.String("Go版本", sysInfo.GoVersion))

	// 记录监控配置
	record.AddAttr(log.Float64("CPU阈值", cfg.CPUThreshold))
	record.AddAttr(log.Float64("内存阈值MB", cfg.MemThreshold))
	record.AddAttr(log.Float64("硬盘阈值比例", cfg.DiskPercentThreshold))
	record.AddAttr(log.String("硬盘监控路径", cfg.DiskPath))
	record.AddAttr(log.Int64("GC时长阈值us", cfg.GCDurationUs))
	record.AddAttr(log.Int64("检查间隔s", cfg.CheckIntervals))

	log.LogRecord(record)
}

// Start 启动监控循环
func (m *SystemMonitor) Start() {
	ticker := time.NewTicker(time.Duration(m.cfg.CheckIntervals) * time.Second)
	m.ticker = ticker
	defer ticker.Stop()
	for range ticker.C {
		alerts := m.checkSysStatus()
		if len(alerts) > 0 {
			m.alert(alerts)
		}
	}
}

func (m *SystemMonitor) Close() {
	m.ticker.Stop()
}

// checkAndAlert 检查所有指标并发送警报
func (m *SystemMonitor) checkSysStatus() (alerts []string) {
	// 检查CPU
	sysCPUUsage, processCPUUsage, err := checkCPU()
	if err != nil {
		alerts = append(alerts, "获取CPU使用量失败: "+err.Error())
	} else if sysCPUUsage > m.cfg.CPUThreshold {
		alerts = append(alerts, fmt.Sprintf("系统CPU使用量异常: %.2f,阈值: %.2f,程序CPU: %.2f",
			sysCPUUsage, m.cfg.CPUThreshold, processCPUUsage))
	}

	// 检查内存
	sysMemUsageMB, processMemUsageMB, err := checkMemory()
	if err != nil {
		alerts = append(alerts, "获取内存使用量失败: "+err.Error())
	} else if sysMemUsageMB > m.cfg.MemThreshold {
		alerts = append(alerts, fmt.Sprintf("系统内存使用量异常MB: %.2f,阈值: %.2f, 程序内存: %.2f",
			sysMemUsageMB, m.cfg.MemThreshold, processMemUsageMB))
	}

	// 检查硬盘（使用可配置路径）
	diskUsagePercent, err := checkDiskPercent(m.cfg.DiskPath)
	if err != nil {
		alerts = append(alerts, "获取硬盘使用率失败: "+err.Error())
	} else if diskUsagePercent > m.cfg.DiskPercentThreshold {
		alerts = append(alerts, fmt.Sprintf("硬盘使用率异常: %.2f,阈值比例: %.2f",
			diskUsagePercent, m.cfg.DiskPercentThreshold))
	}

	// 检查GC时长
	gcDurationUS, err := checkGCDuration()
	if err != nil {
		alerts = append(alerts, "获取GC时长失败: "+err.Error())
	} else if gcDurationUS > m.cfg.GCDurationUs {
		alerts = append(alerts, fmt.Sprintf("GC时长异常: %v us,阈值: %v us",
			gcDurationUS, m.cfg.GCDurationUs))
	}
	// 无论是否有告警，都记录监控状态到日志
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	currentAllocMB := float64(memStats.Alloc) / (1024 * 1024)

	record := log.NewRecord(log.InfoLevel, "系统监控状态")
	record.AddAttr(log.Float64("系统CPU使用量", sysCPUUsage))
	record.AddAttr(log.Float64("程序CPU使用量", processCPUUsage))
	record.AddAttr(log.Float64("内存使用量MB", sysMemUsageMB))
	record.AddAttr(log.Float64("程序内存MB", currentAllocMB))
	record.AddAttr(log.Float64("硬盘使用量GB", diskUsagePercent))
	record.AddAttr(log.Int64("GC时长us", gcDurationUS))
	log.LogRecord(record)
	return alerts
}

func (m *SystemMonitor) alert(alerts []string) {
	// 仅在有告警时发送TG，并附带系统信息
	// 构建告警信息
	alertMsg := "系统监控异常提醒:\n" + strings.Join(alerts, "\n")

	// 附加系统信息
	sysInfoMsg := fmt.Sprintf("\n\n系统信息:\nCPU: %d核\n总内存: %.2fMB\n,硬盘空间: %.2fGB (可用: %.2fGB)\n操作系统: %s\nGo版本: %s",
		m.systemInfo.CPUCores,
		m.systemInfo.TotalMemoryMB,
		m.systemInfo.TotalDiskGB,
		m.systemInfo.FreeDiskGB,
		m.systemInfo.HostOS,
		m.systemInfo.GoVersion)
	// 发送完整信息
	jcdtg.TGSyncWarning(alertMsg + sysInfoMsg)

}

// checkCPU 获取CPU使用量（所有核心的总使用量）
// checkCPU 获取CPU使用量（返回系统总使用量 + 当前进程使用量）
func checkCPU() (systemUsage float64, processUsage float64, err error) {
	// 获取系统整体CPU使用率（保持原有逻辑）
	perCPU, err := cpu.Percent(1*time.Second, true)
	if err != nil {
		return 0, 0, err
	}
	if len(perCPU) == 0 {
		return 0, 0, fmt.Errorf("no CPU cores detected")
	}
	var total float64
	for _, p := range perCPU {
		total += p
	}
	systemUsage = total / 100

	// 新增获取当前进程CPU使用率
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return systemUsage, 0, fmt.Errorf("process.NewProcess failed: %v", err)
	}

	percent, err := p.Percent(1 * time.Second)
	if err != nil {
		return systemUsage, 0, fmt.Errorf("p.Percent failed: %v", err)
	}
	processUsage = percent / 100 // 转换为小数形式（如0.85表示85%）

	return systemUsage, processUsage, nil
}

// checkMemory 获取内存使用量（MB）
func checkMemory() (sysUsedMem float64, processUsedMem float64, err error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}
	// 返回内存使用量（MB）
	usedMemoryMB := float64(v.Used) / (1024 * 1024)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	currentAllocMB := float64(memStats.Alloc) / (1024 * 1024)
	return usedMemoryMB, currentAllocMB, nil
}

// checkDisk 获取指定路径的硬盘使用率（GB）
func checkDiskPercent(path string) (float64, error) {
	d, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}
	return d.UsedPercent / 100, nil
}

// checkGCDuration 获取最近一次GC的暂停时长（微秒us）
func checkGCDuration() (int64, error) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	// PauseNs 是环形缓冲区（长度256），(stats.NumGC+255)%256 用于获取最新的GC暂停时间索引
	lastGCPause := stats.PauseNs[(stats.NumGC+255)%256]
	return int64(lastGCPause / 1e3), nil // 转换为us
}
