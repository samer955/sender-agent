package metrics

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
	"time"
)

var execCommand = exec.Command

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

func newTcp() *Tcp {
	var tcp Tcp
	tcp.Ip = ip
	return &tcp
}

// Working on Windows and Linux in order to get the number of open tcp queue of the bootstrap.
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
