package service

import (
	"context"
	psub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/samer955/sender-agent/bootstrap"
	"github.com/samer955/sender-agent/config"
	"github.com/samer955/sender-agent/producer"

	"log"
	"time"
)

type Sender struct {
	Context       context.Context
	Node          bootstrap.Node
	PubSubService *producer.PubSubService
	Config        config.SenderConfig
}

func NewSender() *Sender {

	ctx := context.Background()
	cfg := config.GetConfig()
	node := bootstrap.InitializeNode(ctx, cfg.DiscoveryTag())
	ps := producer.NewPubSubService(ctx, node.Host)

	return &Sender{
		Context:       ctx,
		Node:          *node,
		PubSubService: ps,
		Config:        cfg,
	}

}

func (s *Sender) subscribeTopics() {

	for _, room := range s.Config.Topics() {
		topic, err := s.PubSubService.JoinTopic(room)
		if err != nil {
			panic(err)
		}

		_, perr := s.PubSubService.Subscribe(topic)
		if perr != nil {
			panic(err)
		}
	}
}

func (s *Sender) publish(data any, topic *psub.Topic) {
	err := s.PubSubService.Publish(data, s.Context, topic)
	if err != nil {
		log.Println("Error publishing data" + err.Error())
	}
	log.Println("Data published: ", data)
}

func (s *Sender) sendSystemInfo() {

	for {
		//it means the local bootstrap is the only bootstrap in the LAN
		if s.Node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}
		s.Node.Metrics.System.Time = time.Now()
		s.Node.Metrics.System.GetOnlineUsers()
		topic, err := s.PubSubService.GetTopic("SYSTEM")
		if err != nil {
			panic(err)
		}
		s.publish(s.Node.Metrics.System, topic)

		time.Sleep(time.Duration(s.Config.Frequency()) * time.Second)
	}

}

func (s *Sender) sendCpuIfo() {

	for {
		//it means the local bootstrap is the only bootstrap in the LAN
		if s.Node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}

		s.Node.Metrics.Cpu.UpdateUtilization()
		s.Node.Metrics.Cpu.Time = time.Now()
		topic, err := s.PubSubService.GetTopic("CPU")
		if err != nil {
			panic(err)
		}
		s.publish(s.Node.Metrics.Cpu, topic)
		time.Sleep(time.Duration(s.Config.Frequency()) * time.Second)
	}

}

func (s *Sender) sendRamInfo() {

	for {
		//it means the local bootstrap is the only bootstrap in the LAN
		if s.Node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}

		//s.bootstrap.Metrics.Memory.UpdateUsagePercent()
		s.Node.Metrics.Memory.Time = time.Now()
		topic, err := s.PubSubService.GetTopic("MEMORY")
		if err != nil {
			panic(err)
		}
		s.publish(s.Node.Metrics.Memory, topic)
		time.Sleep(time.Duration(s.Config.Frequency()) * time.Second)
	}

}

func (s *Sender) sendBandInfo() {

	for {
		//it means the local bootstrap is the only bootstrap in the LAN
		if s.Node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}
		actual := s.Node.BandCounter.GetBandwidthTotals()
		s.Node.Metrics.Bandwidth.TotalIn = actual.TotalIn
		s.Node.Metrics.Bandwidth.TotalOut = actual.TotalOut
		s.Node.Metrics.Bandwidth.RateIn = int(actual.RateIn)
		s.Node.Metrics.Bandwidth.RateOut = int(actual.RateOut)
		s.Node.Metrics.Bandwidth.Time = time.Now()

		topic, err := s.PubSubService.GetTopic("BANDWIDTH")
		if err != nil {
			panic(err)
		}

		s.publish(s.Node.Metrics.Bandwidth, topic)
		time.Sleep(time.Duration(s.Config.Frequency()) * time.Second)
	}

}

func (s *Sender) sendTcpInfo() {

	for {
		//it means the local bootstrap is the only bootstrap in the LAN
		if s.Node.Host.Peerstore().Peers().Len() == 0 {
			continue
		}

		s.Node.Metrics.Tcp.GetConnectionsQueueSize()
		s.Node.Metrics.Tcp.Time = time.Now()
		topic, err := s.PubSubService.GetTopic("TCP")
		if err != nil {
			panic(err)
		}
		s.publish(s.Node.Metrics.Tcp, topic)
		time.Sleep(time.Duration(s.Config.Frequency()) * time.Second)
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
