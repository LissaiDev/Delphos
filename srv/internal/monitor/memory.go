package monitor

import (
	"github.com/shirou/gopsutil/v4/mem"
)

func GetMemoryInfo() (*Memory, error) {
	info, err := mem.VirtualMemory()
	swapInfo, swapErr := mem.SwapMemory()
	if err != nil || swapErr != nil {
		return nil, err
	}
	return &Memory{
		Total:     float64(info.Total) / 1024 / 1024,
		Free:      float64(info.Available) / 1024 / 1024,
		Used:      float64(info.Used) / 1024 / 1024,
		SwapTotal: float64(swapInfo.Total) / 1024 / 1024,
		SwapFree:  float64(swapInfo.Free) / 1024 / 1024,
		SwapUsed:  float64(swapInfo.Used) / 1024 / 1024,
	}, nil
}
