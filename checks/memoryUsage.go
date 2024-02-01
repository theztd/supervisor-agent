package checks

import (
	"github.com/shirou/gopsutil/process"
)

func getMemoryUsageBytes(pid int) (uint64, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0, err
	}
	memInfo, err := p.MemoryInfo()
	if err != nil {
		return 0, err
	}
	return memInfo.RSS, nil
}
