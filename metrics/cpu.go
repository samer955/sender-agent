package metrics

import (
	"bufio"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"os"
	"strings"
	"time"
)

type Cpu struct {
	Ip          string    `json:"ip"`
	Model       string    `json:"model"`
	Utilization float64   `json:"utilization"`
	Time        time.Time `json:"time"`
}

func newCpu() *Cpu {
	var cpu Cpu
	cpu.Ip = ip
	cpu.getCPUModel()
	return &cpu
}

func (c *Cpu) getCPUModel() {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		log.Println("Unable to open /proc/cpuinfo")
		c.Model = ""
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			fields := strings.Split(line, ":")
			if len(fields) > 1 {
				c.Model = strings.TrimSpace(fields[1])
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error scanning")
		c.Model = ""
		return
	}

	c.Model = ""
	log.Println("model name not found in /proc/cpuinfo")
}

func (c *Cpu) UpdateUtilization() {

	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Unable to get Cpu percent usage")
		c.Utilization = 00
		return
	}

	c.Utilization = percent[0]

}
