package monitor

import "fmt"

type Host struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	UpTime   uint64 `json:"uptime"`
}

type Memory struct {
	Total     float64 `json:"total"`
	Used      float64 `json:"used"`
	Free      float64 `json:"free"`
	SwapTotal float64 `json:"swapTotal"`
	SwapUsed  float64 `json:"swapUsed"`
	SwapFree  float64 `json:"swapFree"`
}

type CPU struct {
	Usage float64 `json:"usage"`
	Model string  `json:"model"`
	Cores int     `json:"cores"`
}

type Disk struct {
	Mountpoint  string  `json:"mountpoint"`
	Type        string  `json:"type"`
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	Free        float64 `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type Network struct {
	TotalBytesSent uint64 `json:"totalBytesSent"`
	TotalBytesRecv uint64 `json:"totalBytesRecv"`
	InterfaceName  string `json:"interfaceName"`
}

type Monitor struct {
	Host    *Host      `json:"host"`
	Memory  *Memory    `json:"memory"`
	CPU     []*CPU     `json:"cpu"`
	Disk    []*Disk    `json:"disk"`
	Network []*Network `json:"network"`
}

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
