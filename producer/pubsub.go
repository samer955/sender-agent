package producer

import (
	"context"
	"encoding/json"
	"errors"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"log"
)

type PubSubService struct {
	psub   *pubsub.PubSub
	Topics []*pubsub.Topic
}

// NewPubSubService return a new PubSubService Service using the GossipSub Service
func NewPubSubService(ctx context.Context, host host.Host) *PubSubService {
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		log.Println("unable to create the producer service")
		return nil
	}

	var topics []*pubsub.Topic

	return &PubSubService{
		psub:   ps,
		Topics: topics,
	}
}

// JoinTopic allow the Peers to join a Topic on Pubsub
func (p *PubSubService) JoinTopic(room string) (*pubsub.Topic, error) {

	topic, err := p.psub.Join(room)
	if err != nil {
		log.Println("Error while subscribing in", room)
		return nil, err
	}
	log.Println("Joined room:", room)

	p.Topics = append(p.Topics, topic)

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

// publish any data on a specific topic
func (p *PubSubService) Publish(data any, context context.Context, topic *pubsub.Topic) error {

	msgBytes, err := json.Marshal(data)

	if err != nil {
		log.Println("cannot convert to Bytes ", data)
		return err
	}

	return topic.Publish(context, msgBytes)
}

func (p *PubSubService) GetTopic(topicName string) (*pubsub.Topic, error) {

	for _, topic := range p.Topics {
		if topic.String() == topicName {
			return topic, nil
		}
	}
	return nil, errors.New("Topic:" + " " + topicName + " " + "not found")
}
