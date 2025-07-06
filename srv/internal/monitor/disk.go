package monitor

import (
	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/shirou/gopsutil/v4/disk"
)

func GetDiskInfo() ([]*Disk, error) {
	logger.Log.Debug("Starting disk information collection")

	var disks []*Disk

	// Get disk partitions
	logger.Log.Debug("Collecting disk partition information")
	parts, err := disk.Partitions(false)
	if err != nil {
		logger.Log.Error("Failed to collect disk partition information", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	logger.Log.Debug("Disk partitions found", map[string]interface{}{
		"partition_count": len(parts),
		"partitions": func() []string {
			var mountpoints []string
			for _, part := range parts {
				mountpoints = append(mountpoints, part.Mountpoint)
			}
			return mountpoints
		}(),
	})

	// Get usage information for each partition
	for i, part := range parts {
		logger.Log.Debug("Collecting disk usage for partition", map[string]interface{}{
			"partition_index": i,
			"mountpoint":      part.Mountpoint,
			"filesystem_type": part.Fstype,
		})

		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			logger.Log.Error("Failed to collect disk usage for partition", map[string]interface{}{
				"partition_index": i,
				"mountpoint":      part.Mountpoint,
				"error":           err.Error(),
			})
			return nil, err
		}

		diskInfo := &Disk{
			Mountpoint:  part.Mountpoint,
			Type:        part.Fstype,
			Total:       float64(usage.Total),
			Used:        float64(usage.Used),
			Free:        float64(usage.Free),
			UsedPercent: float64(usage.UsedPercent),
		}

		logger.Log.Debug("Disk usage collected for partition", map[string]interface{}{
			"partition_index": i,
			"mountpoint":      diskInfo.Mountpoint,
			"filesystem_type": diskInfo.Type,
			"total_gb":        diskInfo.Total / 1024 / 1024 / 1024,
			"used_gb":         diskInfo.Used / 1024 / 1024 / 1024,
			"free_gb":         diskInfo.Free / 1024 / 1024 / 1024,
			"used_percent":    diskInfo.UsedPercent,
		})

		disks = append(disks, diskInfo)
	}

	logger.Log.Debug("Disk information collection completed", map[string]interface{}{
		"total_partitions": len(disks),
		"total_space_gb": func() float64 {
			total := 0.0
			for _, d := range disks {
				total += d.Total / 1024 / 1024 / 1024
			}
			return total
		}(),
		"total_used_gb": func() float64 {
			used := 0.0
			for _, d := range disks {
				used += d.Used / 1024 / 1024 / 1024
			}
			return used
		}(),
	})

	return disks, nil
}
