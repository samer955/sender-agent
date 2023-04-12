package metrics

import (
	"bufio"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

type System struct {
	Ip           string    `json:"ip"`
	Hostname     string    `json:"hostname"`
	Os           string    `json:"os"`
	Architecture string    `json:"architecture"`
	Platform     string    `json:"platform"`
	Version      string    `json:"version"`
	Time         time.Time `json:"time"`
}

func newSystem() *System {
	var system System
	system.GetSystemInformation()
	system.getIp()
	return &system
}

func (s *System) GetSystemInformation() {

	s.Hostname = hostname
	s.Os = runtime.GOOS
	s.Architecture = runtime.GOARCH

	platform, version, err := getPlatformAndVersion()
	if err != nil {
		log.Println(err)
	}

	s.Platform, s.Version = platform, version

}

func getPlatformAndVersion() (string, string, error) {

	file, err := os.Open("/etc/os-release")
	if err != nil {
		log.Println("unable to open /etc/os-release")
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var platform, version string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "=")

		switch fields[0] {
		case "VERSION_ID":
			version = strings.Trim(fields[1], "\"")
		case "ID":
			platform = strings.Trim(fields[1], "\"")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		return "", "", err
	}

	return platform, version, nil
}

func (s *System) getIp() {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		s.Ip = ""
		return
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				s.Ip = ipnet.IP.String()
				return
			}
		}
	}
	s.Ip = ""

}
