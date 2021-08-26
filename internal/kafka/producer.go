package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
)

const (
	topic = "team"
)

var brokers = []string{"localhost:9094"}

type Producer interface {
	Send(message Message) error
}

type producer struct {
	actor sarama.SyncProducer
	topic string
}

func NewProducer() (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(brokers, config)

	return &producer{actor: p, topic: topic}, err
}

func (p *producer) Send(message Message) error {
	msg, err := prepareMessage(message)
	if err != nil {
		return err
	}

	_, _, err = p.actor.SendMessage(msg)

	return err
}

func prepareMessage(message Message) (*sarama.ProducerMessage, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(b),
	}

	return msg, nil
}
