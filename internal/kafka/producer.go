package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/ozoncp/ocp-team-api/internal/config"
)

// IProducer is the interface for sending messages to broker.
type IProducer interface {
	Send(message Message) error
}

// Producer is the struct that implements IProducer interface.
type Producer struct {
	actor sarama.SyncProducer
	topic string
}

// NewProducer is the constructor method for Producer struct.
// It returns error if such occurred during constructing.
func NewProducer() (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(config.GetInstance().Kafka.Brokers, saramaConfig)

	return &Producer{actor: p, topic: config.GetInstance().Kafka.Topic}, err
}

// Send is the method that sends message to the broker.
// It returns error if such occurred during either
// message preparing or sending.
func (p *Producer) Send(message Message) error {
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
