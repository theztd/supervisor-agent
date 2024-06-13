package checks

/*
	Deprecated and already not used.

	ToDo: Should be removed in next version
*/

import (
	"github.com/shirou/gopsutil/process"
)

func GetMemoryUsageBytes(pid int) (uint64, error) {
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
