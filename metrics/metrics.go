package metrics

import (
	"os"
	"time"
)

var ip string

func (m *Metrics) InitializeMetrics() {

	os.Hostname()

	ip = m.Ip
	m.System = newSystem()
	m.Cpu = newCpu()
	m.Memory = newMemory()
	m.Bandwidth = newBandwidth()
	m.Tcp = newTcp()

}

type Metrics struct {
	Ip        string
	System    *System
	Cpu       *Cpu
	Memory    *Memory
	Bandwidth *Bandwidth
	Tcp       *Tcp
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

func newBandwidth() *Bandwidth {
	var bandwidth Bandwidth
	bandwidth.Ip = ip
	return &bandwidth
}
