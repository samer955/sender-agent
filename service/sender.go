package service

import (
	"context"
	"errors"
	psub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/samer955/gomdnsdisco/node"
	"github.com/samer955/gomdnsdisco/pubsub"
	"log"
	"time"
)

var (
	ctx         context.Context
	topicsMap   = make(map[string]*psub.Topic)
	roomsTopics = []string{"system", "cpu", "bandwidth", "memory", "tcp"}
)

type Sender struct {
	node   node.Node
	pubsub pubsub.PubSubService
}

func NewSender() *Sender {

	ctx := context.Background()
	node := node.InitializeNode(ctx)
	pubsub := pubsub.NewPubSubService(ctx, node.Host)

	return &Sender{*node, *pubsub}

}

func (s *Sender) subscribeTopics() {

	for _, room := range roomsTopics {
		topic, err := s.pubsub.JoinTopic(room)
		if err != nil {
			panic(err)
		}

		_, perr := s.pubsub.Subscribe(topic)
		if perr != nil {
			panic(err)
		}
		topicsMap[room] = topic
	}
}

func checkTopic(topicName string) (*psub.Topic, error) {

	topic, ok := topicsMap[topicName]
	if ok {
		return topic, nil
	}
	return nil, errors.New("topic not found")
}

func (s *Sender) publish(data any, topic *psub.Topic) {
	err := s.pubsub.Publish(data, ctx, topic)
	if err != nil {
		log.Println("Error publishing data", err)
	}
	log.Println("Data published: ", data)
}

func (s *Sender) sendSystemInfo() {

	for {
		//it means the local node is the only node in the LAN
		if s.node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}
		s.node.Metrics.System.Time = time.Now()
		topic, err := checkTopic("system")
		if err != nil {
			panic(err)
		}
		s.publish(s.node.Metrics.System, topic)
		time.Sleep(10 * time.Second)
	}

}

func (s *Sender) sendCpuIfo() {

	for {
		//it means the local node is the only node in the LAN
		if s.node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}

		s.node.Metrics.Cpu.UpdateUsagePercent()
		s.node.Metrics.Cpu.Time = time.Now()
		topic, err := checkTopic("cpu")
		if err != nil {
			panic(err)
		}
		s.publish(s.node.Metrics.Cpu, topic)
		time.Sleep(10 * time.Second)
	}

}

func (s *Sender) sendRamInfo() {

	for {
		//it means the local node is the only node in the LAN
		if s.node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}

		s.node.Metrics.Ram.UpdateUsagePercent()
		s.node.Metrics.Ram.Time = time.Now()
		topic, err := checkTopic("memory")
		if err != nil {
			panic(err)
		}
		s.publish(s.node.Metrics.Ram, topic)
		time.Sleep(10 * time.Second)
	}

}

func (s *Sender) sendBandInfo() {

	for {
		//it means the local node is the only node in the LAN
		if s.node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}
		actual := s.node.BandCounter.GetBandwidthTotals()
		s.node.Metrics.Bandwidth.TotalIn = actual.TotalIn
		s.node.Metrics.Bandwidth.TotalOut = actual.TotalOut
		s.node.Metrics.Bandwidth.RateIn = int(actual.RateIn)
		s.node.Metrics.Bandwidth.RateOut = int(actual.RateOut)
		s.node.Metrics.Bandwidth.Time = time.Now()

		topic, err := checkTopic("bandwidth")
		if err != nil {
			panic(err)
		}

		s.publish(s.node.Metrics.Bandwidth, topic)
		time.Sleep(10 * time.Second)
	}

}

func (s *Sender) sendTcpInfo() {

	for {
		//it means the local node is the only node in the LAN
		if s.node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}

		s.node.Metrics.Tcp.UpdateQueueSize()
		s.node.Metrics.Tcp.Time = time.Now()
		topic, err := checkTopic("tcp")
		if err != nil {
			panic(err)
		}
		s.publish(s.node.Metrics.Tcp, topic)
		time.Sleep(10 * time.Second)
	}

}

func (s *Sender) Start() {

	s.subscribeTopics()
	go s.sendSystemInfo()
	go s.sendBandInfo()
	go s.sendCpuIfo()
	go s.sendRamInfo()
	go s.sendTcpInfo()

}
