package util

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
)

func GetCpuPercent(d time.Duration) float64 {
	percent, err := cpu.Percent(d, true)
	if err != nil || len(percent) == 0 {
		return 0
	}
	total := 0.0
	for _, p := range percent {
		total += p
	}
	return total / float64(len(percent))
}

func GetDiskPercent() float64 {
	status, err := disk.Usage(AppPath())
	if err != nil {
		return 0.0
	}
	return status.UsedPercent
}
