package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/derkres11/price-pulse/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(addr string) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (c *Cache) SetPrice(ctx context.Context, productID int64, price float64) error {
	key := fmt.Sprintf("product_price:%d", productID)
	return c.client.Set(ctx, key, price, 24*time.Hour).Err()
}

func (c *Cache) GetPrice(ctx context.Context, productID int64) (float64, error) {
	key := fmt.Sprintf("product_price:%d", productID)
	return c.client.Get(ctx, key).Float64()
}

func (c *Cache) Get(ctx context.Context, id int64) (*domain.Product, error) {
	key := fmt.Sprintf("product:%d", id)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var p domain.Product
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (c *Cache) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("product:%d", id)
	priceKey := fmt.Sprintf("product_price:%d", id)

	// Удаляем и сам объект, и закэшированную цену
	return c.client.Del(ctx, key, priceKey).Err()
}
