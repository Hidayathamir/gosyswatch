package proc

import (
	"github.com/prometheus/procfs"
)

// GetMemoryUsagePercentage retrieves the memory usage percentage.
func GetMemoryUsagePercentage() (float64, error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return 0, err
	}

	memInfo, err := fs.Meminfo()
	if err != nil {
		return 0, err
	}

	memTotal := *memInfo.MemTotal
	memAvailable := *memInfo.MemAvailable
	memUsage := memTotal - memAvailable
	memUsagePercentage := float64(memUsage) / float64(memTotal) * 100

	return memUsagePercentage, nil
}
