package metrics

import (
	"bufio"
	"log"
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
	system.Ip = ip
	system.GetSystemInformation()
	return &system
}

func (s *System) GetSystemInformation() {

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Unable to get the hostname")
		hostname = ""
	}
	s.Hostname = hostname
	s.Os = runtime.GOOS
	s.Architecture = runtime.GOARCH

	s.Platform, s.Version, err = getPlatformAndVersion()
	if err != nil {
		log.Println(err)
	}

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
