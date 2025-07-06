package monitor

import (
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

// GetSystemStats retrieves comprehensive system statistics
// Collects data from all monitoring modules and returns a consolidated view
func GetSystemStats() (*Monitor, error) {
	startTime := time.Now()

	logger.Log.Debug("Starting system statistics collection", map[string]interface{}{
		"timestamp": startTime.Format(time.RFC3339),
	})

	// Collect host information
	logger.Log.Debug("Collecting host information")
	host, err := GetHostInfo()
	if err != nil {
		logger.Log.Error("Failed to collect host information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	logger.Log.Debug("Host information collected successfully", map[string]interface{}{
		"hostname": host.Hostname,
		"os":       host.OS,
		"uptime":   host.UpTime,
	})

	// Collect memory information
	logger.Log.Debug("Collecting memory information")
	mem, err := GetMemoryInfo()
	if err != nil {
		logger.Log.Error("Failed to collect memory information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	logger.Log.Debug("Memory information collected successfully", map[string]interface{}{
		"total_memory_gb":      mem.Total / 1024 / 1024 / 1024,
		"used_memory_gb":       mem.Used / 1024 / 1024 / 1024,
		"free_memory_gb":       mem.Free / 1024 / 1024 / 1024,
		"memory_usage_percent": (mem.Used / mem.Total) * 100,
	})

	// Collect CPU information
	logger.Log.Debug("Collecting CPU information")
	cpu, err := GetCPUInfo()
	if err != nil {
		logger.Log.Error("Failed to collect CPU information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	logger.Log.Debug("CPU information collected successfully", map[string]interface{}{
		"cpu_count": len(cpu),
		"cpu_model": cpu[0].Model,
		"cpu_cores": cpu[0].Cores,
		"cpu_usage": cpu[0].Usage,
	})

	// Collect disk information
	logger.Log.Debug("Collecting disk information")
	disk, err := GetDiskInfo()
	if err != nil {
		logger.Log.Error("Failed to collect disk information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	logger.Log.Debug("Disk information collected successfully", map[string]interface{}{
		"disk_count": len(disk),
		"partitions": func() []string {
			var mounts []string
			for _, d := range disk {
				mounts = append(mounts, d.Mountpoint)
			}
			return mounts
		}(),
	})

	// Collect network information
	logger.Log.Debug("Collecting network information")
	net, err := GetNetworkInfo()
	if err != nil {
		logger.Log.Error("Failed to collect network information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	logger.Log.Debug("Network information collected successfully", map[string]interface{}{
		"interface_count": len(net),
		"interfaces": func() []string {
			var ifaces []string
			for _, n := range net {
				ifaces = append(ifaces, n.InterfaceName)
			}
			return ifaces
		}(),
	})

	// Calculate total collection time
	collectionTime := time.Since(startTime)

	// Return consolidated system statistics
	result := &Monitor{
		Host:    host,
		Memory:  mem,
		CPU:     cpu,
		Disk:    disk,
		Network: net,
	}

	logger.Log.Info("System statistics collection completed", map[string]interface{}{
		"collection_time":      collectionTime.String(),
		"hostname":             host.Hostname,
		"cpu_cores":            len(cpu),
		"disk_partitions":      len(disk),
		"network_interfaces":   len(net),
		"memory_usage_percent": (mem.Used / mem.Total) * 100,
		"cpu_usage_percent":    cpu[0].Usage,
	})

	return result, nil
}
