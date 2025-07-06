package monitor

import "fmt"

// Host represents system host information
type Host struct {
	Hostname string `json:"hostname"` // System hostname
	OS       string `json:"os"`       // Operating system name
	UpTime   uint64 `json:"uptime"`   // System uptime in seconds
}

// Memory represents system memory statistics
type Memory struct {
	Total     float64 `json:"total"`     // Total physical memory in bytes
	Used      float64 `json:"used"`      // Used physical memory in bytes
	Free      float64 `json:"free"`      // Free physical memory in bytes
	SwapTotal float64 `json:"swapTotal"` // Total swap memory in bytes
	SwapUsed  float64 `json:"swapUsed"`  // Used swap memory in bytes
	SwapFree  float64 `json:"swapFree"`  // Free swap memory in bytes
}

// CPU represents CPU information and usage statistics
type CPU struct {
	Usage float64 `json:"usage"` // CPU usage percentage
	Model string  `json:"model"` // CPU model name
	Cores int     `json:"cores"` // Number of CPU cores
}

// Disk represents disk partition information and usage statistics
type Disk struct {
	Mountpoint  string  `json:"mountpoint"`  // Mount point path
	Type        string  `json:"type"`        // File system type
	Total       float64 `json:"total"`       // Total disk space in bytes
	Used        float64 `json:"used"`        // Used disk space in bytes
	Free        float64 `json:"free"`        // Free disk space in bytes
	UsedPercent float64 `json:"usedPercent"` // Disk usage percentage
}

// Network represents network interface statistics
type Network struct {
	TotalBytesSent uint64 `json:"totalBytesSent"` // Total bytes sent
	TotalBytesRecv uint64 `json:"totalBytesRecv"` // Total bytes received
	InterfaceName  string `json:"interfaceName"`  // Network interface name
}

// Monitor represents comprehensive system monitoring data
// Contains all system statistics in a consolidated structure
type Monitor struct {
	Host    *Host      `json:"host"`    // Host information
	Memory  *Memory    `json:"memory"`  // Memory statistics
	CPU     []*CPU     `json:"cpu"`     // CPU information (one per core)
	Disk    []*Disk    `json:"disk"`    // Disk information (one per partition)
	Network []*Network `json:"network"` // Network statistics (one per interface)
}

// String methods for pretty printing

func (h *Host) String() string {
	return fmt.Sprintf("Hostname: %s\nOS: %s\nUptime: %d",
		h.Hostname, h.OS, h.UpTime)
}

func (m *Memory) String() string {
	return fmt.Sprintf("Total: %.2f GB\nUsed: %.2f GB\nFree: %.2f GB\nSwap Total: %.2f GB\nSwap Used: %.2f GB\nSwap Free: %.2f GB",
		m.Total/1024/1024/1024, m.Used/1024/1024/1024, m.Free/1024/1024/1024,
		m.SwapTotal/1024/1024/1024, m.SwapUsed/1024/1024/1024, m.SwapFree/1024/1024/1024)
}

func (c *CPU) String() string {
	return fmt.Sprintf("Model: %s\nCores: %d\nUsage: %.2f%%",
		c.Model, c.Cores, c.Usage)
}

func (d *Disk) String() string {
	return fmt.Sprintf("Mountpoint: %s\nType: %s\nTotal: %.2f GB\nUsed: %.2f GB\nFree: %.2f GB\nUsed Percent: %.2f%%",
		d.Mountpoint, d.Type, d.Total/1024/1024/1024, d.Used/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
}

func (n *Network) String() string {
	return fmt.Sprintf("Total Bytes Sent: %d\nTotal Bytes Received: %d\nInterface Name: %s",
		n.TotalBytesSent, n.TotalBytesRecv, n.InterfaceName)
}

func (m *Monitor) String() string {
	return fmt.Sprintf("Host: %s\nOS: %s\nUptime: %d\nMemory: %s\nCPU: %s\nDisk: %s\nNetwork: %s",
		m.Host.Hostname, m.Host.OS, m.Host.UpTime,
		m.Memory.String(), m.CPU[0].String(), m.Disk[0].String(), m.Network[0].String())
}
