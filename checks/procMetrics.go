package checks

import (
	"log"

	"github.com/shirou/gopsutil/process"
)

type ProcDetail struct {
	Memory     process.MemoryInfoStat
	CPUPercent float64
	IO         process.IOCountersStat
}

// returns struct of metrics
func GetProcessInfo(pid int) (ProcDetail, error) {
	ret := ProcDetail{}

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return ret, err
	}
	// Get memory
	if memInfoPtr, err := p.MemoryInfo(); err == nil {
		ret.Memory = *memInfoPtr
	} else {
		log.Println(err)
	}

	// Get CPU
	if CPUPtr, err := p.CPUPercent(); err == nil {
		ret.CPUPercent = CPUPtr
	} else {
		log.Println(err)
	}

	// Get IO
	if IOCountPtr, err := p.IOCounters(); err == nil {
		ret.IO = *IOCountPtr
	} else {
		log.Println(err)
	}

	return ret, nil
}
