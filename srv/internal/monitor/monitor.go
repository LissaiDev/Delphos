package monitor

func GetStats() (*Monitor, error) {

	host, err := GetHostInfo()
	if err != nil {
		return nil, err
	}

	mem, err := GetMemoryInfo()
	if err != nil {
		return nil, err
	}

	cpu, err := GetCPUInfo()
	if err != nil {
		return nil, err
	}

	disk, err := GetDiskInfo()
	if err != nil {
		return nil, err
	}

	net, err := GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	return &Monitor{
		Host:    host,
		Memory:  mem,
		CPU:     cpu,
		Disk:    disk,
		Network: net,
	}, nil

}
