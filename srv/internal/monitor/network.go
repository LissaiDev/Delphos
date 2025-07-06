package monitor

import (
	"github.com/LissaiDev/Delphos/pkg/logger"
	"github.com/shirou/gopsutil/v4/net"
)

func GetNetworkInfo() ([]*Network, error) {
	logger.Log.Debug("Starting network information collection")

	var networks []*Network

	// Get network I/O counters
	logger.Log.Debug("Collecting network I/O counters")
	netStats, err := net.IOCounters(true)
	if err != nil {
		logger.Log.Error("Failed to collect network I/O counters", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	logger.Log.Debug("Network interfaces found", map[string]interface{}{
		"interface_count": len(netStats),
		"interfaces": func() []string {
			var names []string
			for _, stat := range netStats {
				names = append(names, stat.Name)
			}
			return names
		}(),
	})

	// Process each network interface
	for i, netStat := range netStats {
		logger.Log.Debug("Processing network interface", map[string]interface{}{
			"interface_index": i,
			"interface_name":  netStat.Name,
			"bytes_sent":      netStat.BytesSent,
			"bytes_recv":      netStat.BytesRecv,
			"packets_sent":    netStat.PacketsSent,
			"packets_recv":    netStat.PacketsRecv,
		})

		network := &Network{
			InterfaceName:  netStat.Name,
			TotalBytesSent: netStat.BytesSent,
			TotalBytesRecv: netStat.BytesRecv,
		}

		networks = append(networks, network)
	}

	logger.Log.Debug("Network information collection completed", map[string]interface{}{
		"total_interfaces": len(networks),
		"total_bytes_sent": func() uint64 {
			total := uint64(0)
			for _, n := range networks {
				total += n.TotalBytesSent
			}
			return total
		}(),
		"total_bytes_recv": func() uint64 {
			total := uint64(0)
			for _, n := range networks {
				total += n.TotalBytesRecv
			}
			return total
		}(),
	})

	return networks, nil
}
