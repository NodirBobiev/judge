package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

type Submission struct {
	Filename string
	Content  []byte
}

type Client struct {
	client   sarama.Client
	producer sarama.SyncProducer
	consumer sarama.Consumer
	offset   int
}

func newProducerAndConsumerFromClient(client sarama.Client) (producer sarama.SyncProducer, consumer sarama.Consumer, err error) {
	producer, err = sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return
	}
	consumer, err = sarama.NewConsumerFromClient(client)
	if err != nil {
		return
	}
	return
}

func NewKafkaClient(brokers []string) (*Client, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.Return.Successes = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}
	producer, consumer, err := newProducerAndConsumerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &Client{
		client:   client,
		producer: producer,
		consumer: consumer,
	}, nil
}

func (kc *Client) Close() error {
	if kc.producer != nil {
		if err := kc.producer.Close(); err != nil {
			return err
		}
	}
	if kc.consumer != nil {
		if err := kc.consumer.Close(); err != nil {
			return err
		}
	}
	if kc.client != nil {
		return kc.client.Close()
	}
	return nil
}

func (kc *Client) SendMessage(topic string, message Submission) error {
	valueEncoder := sarama.ByteEncoder(message.Content)

	_, _, err := kc.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(message.Filename),
		Value: valueEncoder,
	})
	return err
}

func (kc *Client) ConsumeMessageWithContext(ctx context.Context, topic string) (*Submission, error) {

	partitionConsumer, err := kc.consumer.ConsumePartition(topic, 0, int64(kc.offset))
	if err != nil {
		return nil, err
	}
	defer partitionConsumer.Close()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled. Exiting...")
			return nil, nil
		case msg := <-partitionConsumer.Messages():
			message := &Submission{
				Filename: string(msg.Key),
				Content:  msg.Value,
			}
			kc.offset++
			return message, nil
		case err := <-partitionConsumer.Errors():
			return nil, err
		}
	}
}
