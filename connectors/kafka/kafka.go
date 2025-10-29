package kafka

import (
	"encoding/json"
	"fmt"
	"sync"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaConnector sends sensor data to a Kafka topic
type KafkaConnector struct {
	Producer *KafkaProducer
	Topic    string
}

type KafkaProducer struct {
	Producer  *confluentKafka.Producer
	WaitGroup *sync.WaitGroup
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers string) (*KafkaProducer, error) {
	p, err := confluentKafka.NewProducer(&confluentKafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, err
	}

	producer := &KafkaProducer{
		Producer:  p,
		WaitGroup: &sync.WaitGroup{},
	}

	// Handle events (like delivery reports) in a background goroutine
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *confluentKafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}

				producer.WaitGroup.Done()
			}
		}
	}()

	return producer, nil
}

// Send marshals the data and sends it to Kafka
func (kc *KafkaConnector) Send(data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	err = kc.Producer.Producer.Produce(&confluentKafka.Message{
		TopicPartition: confluentKafka.TopicPartition{Topic: &kc.Topic, Partition: confluentKafka.PartitionAny},
		Value:          bytes,
	}, nil)

	if err != nil {
		return fmt.Errorf("error producing message: %w", err)
	}

	kc.Producer.WaitGroup.Add(1)

	return nil
}

func (kp *KafkaProducer) Wait() {
	kp.WaitGroup.Wait()
}
