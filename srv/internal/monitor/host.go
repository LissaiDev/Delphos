package monitor

import (
	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/shirou/gopsutil/v4/host"
)

func GetHostInfo() (*Host, error) {
	log := logger.GetInstance()

	log.Debug("Starting host information collection")

	// Get host information
	log.Debug("Collecting host system information")
	info, err := host.Info()
	if err != nil {
		log.Error("Failed to collect host information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	log.Debug("Host information collected", map[string]interface{}{
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

	log.Debug("Host information collection completed", map[string]interface{}{
		"hostname":       hostInfo.Hostname,
		"os":             hostInfo.OS,
		"uptime_seconds": hostInfo.UpTime,
		"uptime_hours":   hostInfo.UpTime / 3600,
	})

	return hostInfo, nil
}
