package monitor

import "github.com/shirou/gopsutil/v4/disk"

func GetDiskInfo() ([]*Disk, error) {
	var disks []*Disk
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			return nil, err
		}
		disk := &Disk{
			Mountpoint:  part.Mountpoint,
			Type:        part.Fstype,
			Total:       float64(usage.Total),
			Used:        float64(usage.Used),
			Free:        float64(usage.Free),
			UsedPercent: float64(usage.UsedPercent),
		}
		disks = append(disks, disk)
	}
	return disks, nil
}
