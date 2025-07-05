package monitor

import (
	"github.com/shirou/gopsutil/v4/cpu"
)

func GetCPUInfo() ([]*CPU, error) {
	// Obter informações detalhadas das CPUs
	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	// Obter percentual de uso geral da CPU
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	var cpus []*CPU
	usage := 0.0
	if len(percent) > 0 {
		usage = percent[0]
	}

	// Criar uma entrada para cada CPU física
	for _, cpuInfo := range info {
		cpus = append(cpus, &CPU{
			Usage: usage,
			Model: cpuInfo.ModelName,
			Cores: int(cpuInfo.Cores),
		})
	}

	return cpus, nil
}
