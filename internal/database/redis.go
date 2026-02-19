package database

import (
	"context"
	"fmt"
	"time"

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
