package metrics

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Memory struct {
	Hostname    string    `json:"hostname"`
	Total       float64   `json:"total"`
	Utilization float64   `json:"utilization,omitempty"`
	Time        time.Time `json:"time"`
}

func newMemory() *Memory {

	var ram Memory
	ram.Hostname = hostname
	totalRam, err := getTotalMemory()
	if err != nil {
		log.Println(err)
	}
	ram.Total = totalRam
	return &ram

}

func (mem *Memory) GetMemoryUtilization() {

	totalRAM := mem.Total
	freeRAM, err := getFreeMemory()
	if err != nil {
		log.Println(err)
		mem.Utilization = 0
		return
	}
	usedRAM := totalRAM - freeRAM
	usedRAMPercent := (usedRAM / totalRAM) * 100.0

	mem.Utilization = usedRAMPercent

}

func getTotalMemory() (float64, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		log.Println("Unable to open /proc/meminfo")
		return 0, err
	}

	lines := strings.Split(string(data), "\n")

	var totalRAM uint64

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		if fields[0] == "MemTotal:" {
			totalRAM, err = strconv.ParseUint(fields[1], 10, 64)
			return float64(totalRAM), err
		}
	}

	return float64(0), errors.New("unable to find Total Memory in /proc/meminfo")
}

func getFreeMemory() (float64, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		log.Println("Unable to open /proc/meminfo")
		return 0, err
	}

	lines := strings.Split(string(data), "\n")

	var freeRAM uint64

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		if fields[0] == "MemFree:" {
			freeRAM, err = strconv.ParseUint(fields[1], 10, 64)
			return float64(freeRAM), err
		}
	}

	return float64(0), errors.New("unable to find Free Memory in /proc/meminfo")
}
