package metrics

import (
	"bufio"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	execCommand = exec.Command
	ip          string
)

func (m *Metrics) InitializeMetrics() {

	ip = m.Ip
	m.System = newSystem()
	m.Cpu = newCpu()
	m.Ram = newRam()
	m.Bandwidth = newBandwidth()
	m.Tcp = newTcp()

}

type Metrics struct {
	Ip        string
	System    *System
	Cpu       *Cpu
	Ram       *Ram
	Bandwidth *Bandwidth
	Tcp       *Tcp
}

type System struct {
	Ip       string    `json:"ip"`
	Hostname string    `json:"hostname"`
	Os       string    `json:"os"`
	Platform string    `json:"platform"`
	Version  string    `json:"version"`
	Time     time.Time `json:"time"`
}

type Cpu struct {
	Ip    string    `json:"ip"`
	Model string    `json:"model"`
	Usage int       `json:"usage"`
	Time  time.Time `json:"time"`
}

type Ram struct {
	Ip    string    `json:"ip"`
	Usage int       `json:"usage, omitempty"`
	Time  time.Time `json:"time"`
}

// incoming and outgoing data transferred by the local peer.
type Bandwidth struct {
	Ip       string    `json:"ip"`
	TotalIn  int64     `json:"total_in"`
	TotalOut int64     `json:"total_out"`
	RateIn   int       `json:"rate_in"`
	RateOut  int       `json:"rate_out"`
	Time     time.Time `json:"time"`
}

// Queue Size = number of open TCP connections
// Received = number of segments received
// Sent = number of segments sent
type Tcp struct {
	Ip               string    `json:"ip"`
	QueueSize        int       `json:"tcp_queue_size"`
	SegmentsReceived int       `json:"segments_received"`
	SegmentsSent     int       `json:"segments_sent"`
	Time             time.Time `json:"time"`
}

// Working on Windows and Linux in order to get the number of open tcp queue of the node.
// Execution of the "netstat -na" Command in order to get all the ESTABLISHED Queue
func (t *Tcp) UpdateQueueSize() {
	out, err := execCommand("netstat", "-na").Output()
	if err != nil {
		log.Println("Unable to execute netstat -na command: ", err)
		return
	}
	output := string(out)
	tcpQueue, err := numberOfTcpQueue(output)
	if err != nil {
		log.Println("Error processing the data: ", err)
		return
	}
	t.QueueSize = tcpQueue
}

func numberOfTcpQueue(s string) (tcpConn int, err error) {

	var lines [][]string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		if (strings.HasPrefix(words[0], "TCP") || strings.HasPrefix(words[0], "tcp")) &&
			strings.HasPrefix(words[len(words)-1], "ESTAB") {
			lines = append(lines, words)
		}
	}
	err = scanner.Err()
	return len(lines), err
}

// Get the actual RAM Percentage from the system
func (ram *Ram) UpdateUsagePercent() {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Unable to get Memory percent usage")
		ram.Usage = 00
		return
	}
	ram.Usage = int(vmStat.UsedPercent)
}

func (s *System) GetSystemInformation() {

	hostname, _ := os.Hostname()
	s.Hostname = hostname

	hostStat, err := host.Info()

	if err != nil {
		s.Os, s.Version, s.Platform = "", "", ""
		return
	}
	s.Os, s.Platform, s.Version = hostStat.OS, hostStat.Platform, hostStat.PlatformVersion

}

func (c *Cpu) UpdateUsagePercent() {

	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Unable to get Cpu percent usage")
		c.Usage = 00
		return
	}

	c.Usage = int(percent[0])

}

func (c *Cpu) GetModel() {

	cpuStat, err := cpu.Info()

	if err != nil {
		c.Model = ""
		return
	}
	c.Model = cpuStat[0].ModelName

}

func newCpu() *Cpu {
	var cpu Cpu
	cpu.Ip = ip
	cpu.GetModel()
	return &cpu
}

func newSystem() *System {
	var system System
	system.Ip = ip
	system.GetSystemInformation()
	return &system
}

func newBandwidth() *Bandwidth {
	var bandwidth Bandwidth
	bandwidth.Ip = ip
	return &bandwidth
}

func newRam() *Ram {
	var ram Ram
	ram.Ip = ip
	return &ram
}

func newTcp() *Tcp {
	var tcp Tcp
	tcp.Ip = ip
	return &tcp
}
