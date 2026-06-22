# 📨 KAFKA PATTERNS
## Extracted from uber-master

**Status:** Pattern 3/8  
**Location:** `_patterns/03-kafka-patterns/`

---

## Producer Pattern

```go
type Producer struct {
    writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
    return &Producer{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Balancer: &kafka.LeastBytes{},
        },
    }
}

func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
    return p.writer.WriteMessages(ctx, kafka.Message{
        Topic: topic,
        Key:   key,
        Value: value,
    })
}
```

## Consumer Pattern

```go
type Consumer struct {
    reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
    return &Consumer{
        reader: kafka.NewReader(kafka.ReaderConfig{
            Brokers:        brokers,
            Topic:          topic,
            GroupID:        groupID,
            CommitInterval: time.Second,
            StartOffset:    kafka.LastOffset,
        }),
    }
}

func (c *Consumer) Start(ctx context.Context, handler func(msg *kafka.Message) error) {
    for {
        msg, err := c.reader.ReadMessage(ctx)
        if err != nil {
            continue
        }
        
        if err := handler(&msg); err != nil {
            continue // Retry logic
        }
    }
}
```

## Event Envelope Pattern

```go
type EventEnvelope struct {
    EventID    string          `json:"event_id"`
    EventType  string          `json:"event_type"`
    AggregateID string         `json:"aggregate_id"`
    Timestamp  int64           `json:"timestamp"`
    Payload    json.RawMessage `json:"payload"`
}

// All Kafka messages must use this envelope
```

**Pattern 3 Status:** READY FOR USE

---
