package metrics

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var execCommand = exec.Command

// Tcp QueueSize = number of open TCP connections
// Tcp SegmentsReceived = number of segments received
// Tcp SegmentsSent = number of segments sent
type Tcp struct {
	Hostname         string    `json:"hostname"`
	QueueSize        int       `json:"queue_size"`
	SegmentsReceived int       `json:"segments_received"`
	SegmentsSent     int       `json:"segments_sent"`
	Time             time.Time `json:"time"`
}

func newTcp() *Tcp {
	var tcp Tcp
	tcp.Hostname = hostname
	return &tcp
}

// Execution of the "netstat -na" Command in order to get all the ESTABLISHED Queue
func (t *Tcp) GetConnectionsQueueSize() {
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

func (t *Tcp) GetSegments() {

	pr, err := execCommand("netstat", "-st").Output()
	if err != nil {
		log.Println(err)
		t.SegmentsSent, t.SegmentsReceived = 0, 0
		return
	}
	received, sent, err := numbersOfSegments(string(pr))

	if err != nil {
		log.Println(err)
	}
	t.SegmentsReceived = received
	t.SegmentsSent = sent
	return

}

// Format the output of "netstat -na" to find the ESTAB tcp queue
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

func numbersOfSegments(s string) (int, int, error) {

	var segmentsReceived = 0
	var segmentsSent = 0

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		if len(words) == 3 {
			if strings.Contains(words[1], "segments") && strings.Contains(words[2], "received") {
				value, err := strconv.Atoi(words[0])
				if err != nil {
					fmt.Println(err)
				} else {
					segmentsReceived = value
				}
			}
		}
		if len(words) == 4 {
			if strings.Contains(words[1], "segments") && strings.Contains(words[2], "sent") {
				value, err := strconv.Atoi(words[0])
				if err != nil {
					fmt.Println(err)
				} else {
					segmentsSent = value
				}
			}
		}
	}
	err := scanner.Err()
	return segmentsReceived, segmentsSent, err

}
