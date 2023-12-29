package proc

import (
	"slices"
	"time"

	"github.com/prometheus/procfs"
)

type cpuUsageStats struct {
	TotalTicks float64
	IdleTicks  float64
}

// GetCPUsUsageInPercentage calculates CPU usage percentage for each CPU core.
func GetCPUsUsageInPercentage() (cpusUsageInPercentage []float64, err error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return nil, err
	}

	initialCPUStats, err := fs.Stat()
	if err != nil {
		return nil, err
	}

	initialStatsMap := map[int64]cpuUsageStats{}

	for cpuID, stat := range initialCPUStats.CPU {
		totalTicks := stat.User + stat.Nice + stat.System +
			stat.Idle + stat.Iowait + stat.IRQ + stat.SoftIRQ +
			stat.Steal + stat.Guest + stat.GuestNice

		idleTicks := stat.Idle

		initialStatsMap[cpuID] = cpuUsageStats{
			TotalTicks: totalTicks,
			IdleTicks:  idleTicks,
		}
	}

	time.Sleep(time.Millisecond * 1000)

	fs, err = procfs.NewFS("/proc")
	if err != nil {
		return nil, err
	}

	currentCPUStats, err := fs.Stat()
	if err != nil {
		return nil, err
	}

	currentStatsMap := map[int64]float64{}
	for cpuID, stat := range currentCPUStats.CPU {
		totalTicks := stat.User + stat.Nice + stat.System +
			stat.Idle + stat.Iowait + stat.IRQ + stat.SoftIRQ +
			stat.Steal + stat.Guest + stat.GuestNice

		idleTicks := stat.Idle

		deltaTotalTicks := totalTicks - initialStatsMap[cpuID].TotalTicks
		deltaIdleTicks := idleTicks - initialStatsMap[cpuID].IdleTicks

		var usagePercentage float64
		if deltaTotalTicks != 0 {
			usagePercentage = (deltaTotalTicks - deltaIdleTicks) / deltaTotalTicks * float64(100)
		}

		currentStatsMap[cpuID] = usagePercentage
	}

	sortedCPUIDs := []int64{}
	for k := range currentStatsMap {
		sortedCPUIDs = append(sortedCPUIDs, k)
	}

	slices.Sort(sortedCPUIDs)

	for _, cpuID := range sortedCPUIDs {
		cpusUsageInPercentage = append(cpusUsageInPercentage, currentStatsMap[cpuID])
	}

	return cpusUsageInPercentage, nil
}
