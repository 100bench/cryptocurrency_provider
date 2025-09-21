package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewBroker(brokers []string, topic, group string) (*Broker, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no brokers provided")
	}
	if topic == "" {
		return nil, fmt.Errorf("topic is required")
	}
	if group == "" {
		return nil, fmt.Errorf("group ID is required")
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        group,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})

	return &Broker{writer: w, reader: r}, nil
}

func (b *Broker) Produce(ctx context.Context, rates []en.Rate) error {
	if len(rates) == 0 {
		return nil
	}
	msgs := make([]kafka.Message, 0, len(rates))
	for _, r := range rates {
		body, err := json.Marshal(r)
		if err != nil {
			return err
		}
		msgs = append(msgs, kafka.Message{
			Key:   []byte(r.Currency),
			Value: body,
			Time:  time.Now(),
		})
	}
	return b.writer.WriteMessages(ctx, msgs...)
}

func (b *Broker) Consume(ctx context.Context) (<-chan en.Rate, error) {
	out := make(chan en.Rate, 256)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msg, err := b.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				time.Sleep(200 * time.Millisecond)
				continue
			}
			var r en.Rate
			if err := json.Unmarshal(msg.Value, &r); err == nil {
				select {
				case out <- r:
				case <-ctx.Done():
					return
				}
			}
			_ = b.reader.CommitMessages(ctx, msg)
		}
	}()
	return out, nil
}

func (b *Broker) Close() error {
	_ = b.writer.Close()
	return b.reader.Close()
}
