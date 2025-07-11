package monitor

import (
	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/shirou/gopsutil/v4/cpu"
)

// GetCPUInfo returns CPU information for the system
func GetCPUInfo() ([]*CPU, error) {
	return getCPUInfoWithLogger()
}

// getCPUInfoWithLogger returns CPU information using the provided logger
func getCPUInfoWithLogger() ([]*CPU, error) {
	log := logger.GetInstance()

	log.Debug("Starting CPU information collection")

	var cpus []*CPU

	// Get CPU usage percentages
	log.Debug("Collecting CPU usage percentages")
	percent, err := cpu.Percent(0, true)
	if err != nil {
		log.Error("Failed to collect CPU usage percentages", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	log.Debug("CPU usage percentages collected", map[string]interface{}{
		"cpu_count":    len(percent),
		"usage_values": percent,
	})

	for _, usage := range percent {
		cpus = append(cpus, &CPU{Usage: usage})
	}

	// Get CPU detailed information
	log.Debug("Collecting CPU detailed information")
	info, err := cpu.Info()
	if err != nil {
		log.Error("Failed to collect CPU detailed information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	log.Debug("CPU detailed information collected", map[string]interface{}{
		"cpu_info_count": len(info),
	})

	for idx, cpuInfo := range info {
		if idx < len(cpus) {
			cpus[idx].Model = cpuInfo.ModelName
			cpus[idx].Cores = int(cpuInfo.Cores)

			log.Debug("CPU core information processed", map[string]interface{}{
				"core_index": idx,
				"model":      cpuInfo.ModelName,
				"cores":      cpuInfo.Cores,
				"usage":      cpus[idx].Usage,
			})
		}
	}

	log.Debug("CPU information collection completed", map[string]interface{}{
		"total_cpus": len(cpus),
		"avg_usage": func() float64 {
			if len(percent) == 0 {
				return 0
			}
			sum := 0.0
			for _, p := range percent {
				sum += p
			}
			return sum / float64(len(percent))
		}(),
	})

	return cpus, nil
}
