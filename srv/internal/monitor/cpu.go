package monitor

import (
	"github.com/shirou/gopsutil/v4/cpu"
)

func GetCPUInfo() ([]*CPU, error) {
	var cpus []*CPU
	percent, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}
	for _, usage := range percent {
		cpus = append(cpus, &CPU{Usage: usage})
	}

	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	for idx, cpuInfo := range info {
		cpus[idx].Model = cpuInfo.ModelName
		cpus[idx].Cores = int(cpuInfo.Cores)
	}

	return cpus, nil
}
