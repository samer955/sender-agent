package node

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	band "github.com/libp2p/go-libp2p/core/metrics"
	"github.com/samer955/gomdnsdisco/metrics"
	"log"
	"net"
)

const discoveryTag = "discovery"

type Node struct {
	Ip          string
	Host        host.Host
	BandCounter *band.BandwidthCounter
	Metrics     *metrics.Metrics
}

func InitializeNode(ctx context.Context) *Node {

	n := new(Node)
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

	host, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"), libp2p.BandwidthReporter(n.BandCounter))

	if err != nil {
		panic(err)
	}
	n.Host = host
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

	peerChan := initMDNS(n.Host, discoveryTag)
	for {
		peer := <-peerChan
		fmt.Println("Found peer:", peer.ID.ShortString(), ", connecting")

		if err := n.Host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}
	}
}
