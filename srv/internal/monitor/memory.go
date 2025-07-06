package monitor

import (
	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetMemoryInfo() (*Memory, error) {
	logger.Log.Debug("Starting memory information collection")

	// Get virtual memory information
	logger.Log.Debug("Collecting virtual memory information")
	info, err := mem.VirtualMemory()
	if err != nil {
		logger.Log.Error("Failed to collect virtual memory information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	logger.Log.Debug("Virtual memory information collected", map[string]interface{}{
		"total_bytes":     info.Total,
		"available_bytes": info.Available,
		"used_bytes":      info.Used,
		"usage_percent":   info.UsedPercent,
	})

	// Get swap memory information
	logger.Log.Debug("Collecting swap memory information")
	swapInfo, swapErr := mem.SwapMemory()
	if swapErr != nil {
		logger.Log.Error("Failed to collect swap memory information", map[string]interface{}{
			"error": swapErr.Error(),
		})
		return nil, swapErr
	}

	logger.Log.Debug("Swap memory information collected", map[string]interface{}{
		"swap_total_bytes": swapInfo.Total,
		"swap_free_bytes":  swapInfo.Free,
		"swap_used_bytes":  swapInfo.Used,
	})

	// Create memory structure with converted values (MB)
	memoryInfo := &Memory{
		Total:     float64(info.Total) / 1024 / 1024,
		Free:      float64(info.Available) / 1024 / 1024,
		Used:      float64(info.Used) / 1024 / 1024,
		SwapTotal: float64(swapInfo.Total) / 1024 / 1024,
		SwapFree:  float64(swapInfo.Free) / 1024 / 1024,
		SwapUsed:  float64(swapInfo.Used) / 1024 / 1024,
	}

	logger.Log.Debug("Memory information collection completed", map[string]interface{}{
		"total_memory_mb": memoryInfo.Total,
		"free_memory_mb":  memoryInfo.Free,
		"used_memory_mb":  memoryInfo.Used,
		"memory_usage_percent": func() float64 {
			if memoryInfo.Total > 0 {
				return (memoryInfo.Used / memoryInfo.Total) * 100
			}
			return 0
		}(),
		"swap_total_mb": memoryInfo.SwapTotal,
		"swap_free_mb":  memoryInfo.SwapFree,
		"swap_used_mb":  memoryInfo.SwapUsed,
	})

	return memoryInfo, nil
}
