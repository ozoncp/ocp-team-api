package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/ozoncp/ocp-team-api/internal/config"
)

type Producer interface {
	Send(message Message) error
}

type producer struct {
	actor sarama.SyncProducer
	topic string
}

func NewProducer() (Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(config.GetInstance().Kafka.Brokers, saramaConfig)

	return &producer{actor: p, topic: config.GetInstance().Kafka.Topic}, err
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
		Topic:     config.GetInstance().Kafka.Topic,
		Partition: -1,
		Value:     sarama.StringEncoder(b),
	}

	return msg, nil
}
