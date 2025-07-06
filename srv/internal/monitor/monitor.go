package monitor

// GetSystemStats retrieves comprehensive system statistics
// Collects data from all monitoring modules and returns a consolidated view
func GetSystemStats() (*Monitor, error) {
	// Collect host information
	host, err := GetHostInfo()
	if err != nil {
		return nil, err
	}

	// Collect memory information
	mem, err := GetMemoryInfo()
	if err != nil {
		return nil, err
	}

	// Collect CPU information
	cpu, err := GetCPUInfo()
	if err != nil {
		return nil, err
	}

	// Collect disk information
	disk, err := GetDiskInfo()
	if err != nil {
		return nil, err
	}

	// Collect network information
	net, err := GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	// Return consolidated system statistics
	return &Monitor{
		Host:    host,
		Memory:  mem,
		CPU:     cpu,
		Disk:    disk,
		Network: net,
	}, nil
}
