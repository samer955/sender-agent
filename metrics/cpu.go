package metrics

import (
	"bufio"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type Cpu struct {
	Hostname    string    `json:"hostname"`
	Model       string    `json:"model"`
	Utilization float64   `json:"utilization"`
	Time        time.Time `json:"time"`
}

func newCpu() *Cpu {
	var cpu Cpu
	cpu.Hostname = hostname
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
		log.Println("Unable to get Cpu percent usage. Set default to 0.00%")
		c.Utilization = 0.00
		return
	}

	//Round Float. 2 places after decimal --> e.g. 50.25
	res := math.Round(percent[0]*100) / 100

	c.Utilization = res

}
