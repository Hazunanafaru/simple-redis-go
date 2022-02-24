package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	//Pool
	var pool *redis.Pool
	redisPassword := os.Getenv("REDIS_PASS")
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = ":6379"
	}

	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisHost, redis.DialPassword(redisPassword))
			if err != nil {
				return nil, err
			}

			return conn, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// Stop Channel
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-stopChan

		pool.Close()

		os.Exit(0)
	}()

	// Start the pool
	if err := start(context.Background(), pool); err != nil {
		log.Fatal(err)
	}
}
