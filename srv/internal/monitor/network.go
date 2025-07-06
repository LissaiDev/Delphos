package monitor

import "github.com/shirou/gopsutil/v4/net"

func GetNetworkInfo() ([]*Network, error) {
	var networks []*Network
	netStats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	for _, netStat := range netStats {
		network := &Network{
			InterfaceName:  netStat.Name,
			TotalBytesSent: netStat.BytesSent,
			TotalBytesRecv: netStat.BytesRecv,
		}
		networks = append(networks, network)
	}

	return networks, nil
}
