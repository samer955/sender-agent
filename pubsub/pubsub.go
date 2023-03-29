package pubsub

import (
	"context"
	"encoding/json"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"log"
)

type PubSubService struct {
	psub *pubsub.PubSub
}

// NewPubSubService return a new PubSub Service using the GossipSub Service
func NewPubSubService(ctx context.Context, host host.Host) *PubSubService {
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Println("unable to create the pubsub service")
		return nil
	}
	return &PubSubService{psub: ps}
}

// JoinTopic allow the Peers to join a Topic on Pubsub
func (p *PubSubService) JoinTopic(room string) (*pubsub.Topic, error) {

	topic, err := p.psub.Join(room)
	if err != nil {
		log.Println("Error while subscribing in", room)
		return nil, err
	}
	log.Println("Joined room:", room)

	return topic, nil

}

// Subscribe returns a new Subscription for the topic.
func (p *PubSubService) Subscribe(topic *pubsub.Topic) (*pubsub.Subscription, error) {

	subscribe, err := topic.Subscribe()

	if err != nil {
		log.Println("cannot subscribe to: ", topic.String())
		return nil, err
	}
	log.Println("Subscribed to topic: " + subscribe.Topic())

	return subscribe, nil
}

func (p *PubSubService) Publish(data any, context context.Context, topic *pubsub.Topic) error {

	//JSON encoding of cpu in order to send the data as []byte.
	msgBytes, err := json.Marshal(data)

	if err != nil {
		log.Println("cannot convert to Bytes ", data)
		return err
	}
	//public the data in the topic
	return topic.Publish(context, msgBytes)
}
