package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type ProductConsumer struct {
	reader *kafka.Reader
}

func NewProductConsumer(brokers []string, topic string, groupID string) *ProductConsumer {
	return &ProductConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *ProductConsumer) Start(ctx context.Context, processFunc func(id int64) error) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error while receiving message: %s", err.Error())
			continue
		}

		var data struct {
			ProductID int64 `json:"product_id"`
		}

		if err := json.Unmarshal(msg.Value, &data); err != nil {
			log.Printf("error unmarshaling message: %s", err.Error())
			continue
		}

		if err := processFunc(data.ProductID); err != nil {
			log.Printf("error processing product %d: %s", data.ProductID, err.Error())
		}
	}
}

func (c *ProductConsumer) Close() error {
	return c.reader.Close()
}
