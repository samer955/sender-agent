package bootstrap

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// Initialize the MDNS service
func initMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	notifee := &discoveryNotifee{}
	notifee.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, notifee)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return notifee.PeerChan
}
