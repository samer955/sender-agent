package bootstrap

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	band "github.com/libp2p/go-libp2p/core/metrics"
	"github.com/samer955/sender-agent/metrics"
	"log"
	"net"
)

type Node struct {
	Ip           string
	Host         host.Host
	BandCounter  *band.BandwidthCounter
	Metrics      *metrics.Metrics
	discoveryTag string
}

func InitializeNode(ctx context.Context, discoveryTag string) *Node {

	n := new(Node)
	n.discoveryTag = discoveryTag
	n.initializeBandCounter()
	n.createLibp2pHost()
	n.getIp()
	n.initMetrics()

	go n.findPeers(ctx)

	return n

}

func (n *Node) initMetrics() {

	metrics := metrics.Metrics{Ip: n.Ip}
	metrics.InitializeMetrics()
	n.Metrics = &metrics

}

// initialize Node using Libp2p, listening all ip4 address and default tcp port
func (n *Node) createLibp2pHost() {

	libp2phost, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"), libp2p.BandwidthReporter(n.BandCounter))

	if err != nil {
		log.Println("Unable to create a Libp2p-Host")
		panic(err)
	}
	n.Host = libp2phost
	log.Printf("New node initialized with host-ID %s\n", n.Host.ID().ShortString())

}

func (n *Node) initializeBandCounter() {
	n.BandCounter = band.NewBandwidthCounter()
}

func (n *Node) getIp() {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		n.Ip = ""
		return
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				n.Ip = ipnet.IP.String()
				return
			}
		}
	}
	n.Ip = ""

}

func (n *Node) findPeers(ctx context.Context) {

	peerChan := initMDNS(n.Host, n.discoveryTag)
	for {
		peer := <-peerChan
		fmt.Println("Found peer:", peer.ID.ShortString(), ", connecting")

		if err := n.Host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}
	}
}
