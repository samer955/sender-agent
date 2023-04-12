package metrics

import (
	"log"
	"os"
	"time"
)

var hostname string

func (m *Metrics) InitializeMetrics() {

	hostname = getHostname()

	m.System = newSystem()
	m.Cpu = newCpu()
	m.Memory = newMemory()
	m.Bandwidth = newBandwidth()
	m.Tcp = newTcp()

}

type Metrics struct {
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
	return &bandwidth
}

func getHostname() string {

	host, err := os.Hostname()
	if err != nil {
		log.Println("Unable to get the hostname")
		host = ""
	}
	return host

}
