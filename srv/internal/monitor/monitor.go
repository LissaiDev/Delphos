package monitor

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/LissaiDev/Delphos/internal/config"
	"github.com/LissaiDev/Delphos/pkg/echo"
	"github.com/LissaiDev/Delphos/pkg/logger"
)

// StatsService handles system statistics collection and management
type StatsService struct {
	logger   logger.BasicLogger
	notifier echo.Notifier
}

var (
	StatsServiceInstance *StatsService
	once                 sync.Once
)

// New creates a new stats service instance
func New(log logger.BasicLogger, notifier echo.Notifier) *StatsService {
	return &StatsService{
		logger:   log,
		notifier: notifier,
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

	// ALERT: Check thresholds and notify if necessary
	cfg := config.Env

	// CPU: average across all cores
	if len(cpu) > 0 {
		sum := 0.0
		for _, c := range cpu {
			sum += c.Usage
		}
		avg := sum / float64(len(cpu))
		if avg > cfg.CPUThreshold {
			_ = s.notifier.Notify(
				fmt.Sprintf("ALERT: CPU usage above threshold (%.1f%% > %.1f%%)", avg, cfg.CPUThreshold),
			)
		}
	}

	// Memory
	if mem.Total > 0 {
		memPercent := (mem.Used / mem.Total) * 100
		if memPercent > cfg.MemoryThreshold {
			_ = s.notifier.Notify(
				fmt.Sprintf("ALERT: Memory usage above threshold (%.1f%% > %.1f%%)", memPercent, cfg.MemoryThreshold),
			)
		}
	}

	// Disk: any partition above the threshold
	for _, d := range disk {
		if d.UsedPercent > cfg.DiskThreshold {
			_ = s.notifier.Notify(
				fmt.Sprintf("ALERT: Disk usage on %s above threshold (%.1f%% > %.1f%%)", d.Mountpoint, d.UsedPercent, cfg.DiskThreshold),
			)
		}
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

func GetInstance() *StatsService {
	log := logger.GetInstance()
	echo := echo.GetInstance()
	once.Do(func() {
		StatsServiceInstance = New(log, echo)
	})
	return StatsServiceInstance
}

// GetSystemStats provides backward compatibility
func GetSystemStats() (*Monitor, error) {
	service := GetInstance()
	return service.GetStats()
}
