package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() *pgxpool.Pool {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 5; i++ {
		pool, err = pgxpool.New(context.Background(), dsn)

		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				log.Println("✅ Successfully connected to Postgres")
				return pool
			}
		}

		log.Printf("⚠️  DB connection attempt %d failed: %v. Retrying in 2s...", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("❌ Could not connect to Postgres after 5 attempts: %v", err)
	return nil
}
