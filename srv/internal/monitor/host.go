package monitor

import (
	"github.com/shirou/gopsutil/v4/host"
)

func GetHostInfo() (*Host, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	return &Host{
		Hostname: info.Hostname,
		OS:       info.OS,
		UpTime:   info.Uptime,
	}, nil
}
