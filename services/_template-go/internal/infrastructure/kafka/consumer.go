package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

// KafkaProducer publishes events to Kafka
type KafkaProducer struct {
	producer sarama.AsyncProducer
	logger   *log.Logger
}

// NewKafkaProducer creates a new producer
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &KafkaProducer{
		producer: producer,
		logger:   log.New(nil, "kafka-producer: ", log.LstdFlags),
	}, nil
}

// Publish sends an event to Kafka
func (p *KafkaProducer) Publish(ctx context.Context, topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("content-type"),
				Value: []byte("application/json"),
			},
		},
	}

	p.producer.Input() <- msg

	return nil
}

// KafkaConsumer consumes events from Kafka
type KafkaConsumer struct {
	consumer      sarama.ConsumerGroup
	logger        *log.Logger
	handlers      map[string]EventHandler
	handlersLock  sync.RWMutex
}

// EventHandler processes Kafka events
type EventHandler interface {
	Handle(ctx context.Context, message []byte) error
}

// NewKafkaConsumer creates a new consumer group
func NewKafkaConsumer(brokers []string, groupID string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return &KafkaConsumer{
		consumer:     consumer,
		logger:       log.New(nil, "kafka-consumer: ", log.LstdFlags),
		handlers:     make(map[string]EventHandler),
	}, nil
}

// RegisterHandler registers event handler for topic
func (c *KafkaConsumer) RegisterHandler(topic string, handler EventHandler) {
	c.handlersLock.Lock()
	defer c.handlersLock.Unlock()

	c.handlers[topic] = handler
}

// Start begins consuming messages
func (c *KafkaConsumer) Start(ctx context.Context, topics []string) error {
	go func() {
		for err := range c.consumer.Errors() {
			c.logger.Printf("consumer error: %v", err)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := c.consumer.Consume(ctx, topics, c); err != nil {
					c.logger.Printf("consumer error: %v", err)
				}
			}
		}
	}()

	return nil
}

// Setup is run at the beginning of a new session
func (c *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session
func (c *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from partition
func (c *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			c.handlersLock.RLock()
			handler, exists := c.handlers[message.Topic]
			c.handlersLock.RUnlock()

			if exists {
				if err := handler.Handle(context.Background(), message.Value); err != nil {
					c.logger.Printf("error handling message from %s: %v", message.Topic, err)
				}
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// Close closes the consumer
func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}
