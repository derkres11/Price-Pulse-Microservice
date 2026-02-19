package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type ProductProducer struct {
	writer *kafka.Writer
}

func NewProductProducer(brokers []string, topic string) *ProductProducer {
	return &ProductProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *ProductProducer) SendProductUpdate(ctx context.Context, productID int64) error {
	message := map[string]interface{}{
		"product_id": productID,
		"action":     "check_price",
	}

	payload, _ := json.Marshal(message)

	err := p.writer.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})

	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *ProductProducer) Close() error {
	return p.writer.Close()
}
