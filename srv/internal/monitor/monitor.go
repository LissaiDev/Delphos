package monitor

import (
	"encoding/json"
	"time"

	"github.com/LissaiDev/Delphos/pkg/logger"
)

// StatsService handles system statistics collection and management
type StatsService struct {
	logger logger.BasicLogger
}

// NewStatsService creates a new stats service instance
func NewStatsService(log logger.BasicLogger) *StatsService {
	return &StatsService{
		logger: log,
	}
}

// GetStats retrieves comprehensive system statistics
func (s *StatsService) GetStats() (*Monitor, error) {
	startTime := time.Now()

	s.logger.Debug("Starting system statistics collection", map[string]interface{}{
		"timestamp": startTime.Format(time.RFC3339),
	})

	// Collect all system information
	host, err := s.collectHostInfo()
	if err != nil {
		return nil, err
	}

	mem, err := s.collectMemoryInfo()
	if err != nil {
		return nil, err
	}

	cpu, err := s.collectCPUInfo()
	if err != nil {
		return nil, err
	}

	disk, err := s.collectDiskInfo()
	if err != nil {
		return nil, err
	}

	net, err := s.collectNetworkInfo()
	if err != nil {
		return nil, err
	}

	// Create consolidated result
	result := &Monitor{
		Host:    host,
		Memory:  mem,
		CPU:     cpu,
		Disk:    disk,
		Network: net,
	}

	s.logCompletionStats(result, time.Since(startTime))

	return result, nil
}

// GetStatsJSON returns system statistics as JSON
func (s *StatsService) GetStatsJSON() ([]byte, error) {
	stats, err := s.GetStats()
	if err != nil {
		return nil, err
	}
	return json.Marshal(stats)
}

// collectWithLogging is a DRY helper for collecting system info with consistent logging
func (s *StatsService) collectHostInfo() (*Host, error) {
	s.logger.Debug("Collecting host information")
	result, err := GetHostInfo()
	if err != nil {
		s.logger.Error("Failed to collect host information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	s.logger.Debug("Host information collected successfully")
	return result, nil
}

func (s *StatsService) collectMemoryInfo() (*Memory, error) {
	s.logger.Debug("Collecting memory information")
	result, err := GetMemoryInfo()
	if err != nil {
		s.logger.Error("Failed to collect memory information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	s.logger.Debug("Memory information collected successfully")
	return result, nil
}

func (s *StatsService) collectCPUInfo() ([]*CPU, error) {
	s.logger.Debug("Collecting cpu information")
	result, err := GetCPUInfo()
	if err != nil {
		s.logger.Error("Failed to collect cpu information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	s.logger.Debug("CPU information collected successfully")
	return result, nil
}

func (s *StatsService) collectDiskInfo() ([]*Disk, error) {
	s.logger.Debug("Collecting disk information")
	result, err := GetDiskInfo()
	if err != nil {
		s.logger.Error("Failed to collect disk information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	s.logger.Debug("Disk information collected successfully")
	return result, nil
}

func (s *StatsService) collectNetworkInfo() ([]*Network, error) {
	s.logger.Debug("Collecting network information")
	result, err := GetNetworkInfo()
	if err != nil {
		s.logger.Error("Failed to collect network information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	s.logger.Debug("Network information collected successfully")
	return result, nil
}

// logCompletionStats logs the completion statistics
func (s *StatsService) logCompletionStats(result *Monitor, duration time.Duration) {
	s.logger.Info("System statistics collection completed", map[string]interface{}{
		"collection_time":    duration.String(),
		"hostname":           result.Host.Hostname,
		"cpu_cores":          len(result.CPU),
		"disk_partitions":    len(result.Disk),
		"network_interfaces": len(result.Network),
		"memory_usage_percent": func() float64 {
			if result.Memory.Total > 0 {
				return (result.Memory.Used / result.Memory.Total) * 100
			}
			return 0
		}(),
		"cpu_usage_percent": func() float64 {
			if len(result.CPU) > 0 {
				return result.CPU[0].Usage
			}
			return 0
		}(),
	})
}

// GetSystemStats provides backward compatibility
func GetSystemStats() (*Monitor, error) {
	service := NewStatsService(logger.Log)
	return service.GetStats()
}
