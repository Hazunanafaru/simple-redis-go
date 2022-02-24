package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Pool
	// Make sure to have a redis server or container
	// You can pull a redis container with this code
	// docker pull redis:6.2-alpine
	// and run it with
	// docker run -d -p 6379:6379 redis:6.2-alpine
	var pool *redis.Pool
	redisHost := ":6379"

	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 120 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisHost)
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

func start(ctx context.Context, pool *redis.Pool) error {
	var err error

	for i := 0; i < 10; i++ {
		// Retrieve connection to Reddis Pool
		conn := pool.Get()

		// Example of SET command
		_, err := redis.Bytes(conn.Do("SET", fmt.Sprintf("key-is-%v", i), fmt.Sprintf("value-is-%v", i)))
		if err != nil {
			log.Error(err)
		}

		// Example of GET command
		data, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("key-is-%v", i)))
		if err != nil {
			log.Error(err)
		} else {
			log.Info("Got value: %v", string(data))
		}
		conn.Close()
	}

	return err
}
