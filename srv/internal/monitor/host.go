package monitor

import (
	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/shirou/gopsutil/v4/host"
)

func GetHostInfo() (*Host, error) {
	logger.Log.Debug("Starting host information collection")

	// Get host information
	logger.Log.Debug("Collecting host system information")
	info, err := host.Info()
	if err != nil {
		logger.Log.Error("Failed to collect host information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	logger.Log.Debug("Host information collected", map[string]interface{}{
		"hostname":         info.Hostname,
		"os":               info.OS,
		"platform":         info.Platform,
		"platform_family":  info.PlatformFamily,
		"platform_version": info.PlatformVersion,
		"kernel_version":   info.KernelVersion,
		"uptime":           info.Uptime,
	})

	hostInfo := &Host{
		Hostname: info.Hostname,
		OS:       info.OS,
		UpTime:   info.Uptime,
	}

	logger.Log.Debug("Host information collection completed", map[string]interface{}{
		"hostname":       hostInfo.Hostname,
		"os":             hostInfo.OS,
		"uptime_seconds": hostInfo.UpTime,
		"uptime_hours":   hostInfo.UpTime / 3600,
	})

	return hostInfo, nil
}
