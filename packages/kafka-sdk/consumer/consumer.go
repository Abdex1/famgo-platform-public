
// STEP 10 — CREATE CONSUMER

 // packages/kafka-sdk/consumer/consumer.go

//IMPLEMENTATION

package consumer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type HandlerFunc func(ctx context.Context, message kafka.Message) error

type Consumer struct {
	reader *kafka.Reader
}

func New(
	brokers []string,
	groupID string,
	topic string,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader: reader,
	}
}

func (c *Consumer) Consume(
	ctx context.Context,
	handler HandlerFunc,
) {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("consumer read error: %v", err)
			continue
		}

		if err := handler(ctx, message); err != nil {
			log.Printf("handler error: %v", err)
		}
	}
}
