package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Akash5106/distributed-task-queue/internal/server"
	"github.com/Akash5106/distributed-task-queue/internal/storage"
	"github.com/Akash5106/distributed-task-queue/internal/worker"
)

func main() {
	RedisClient := storage.NewRedisClient()
	go func() {
		for {
			err := RedisClient.RecoverStuckTasks(
				context.Background(),
			)

			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(2 * time.Second)
		}
	}()
	s := server.NewServer(RedisClient)
	w1 := worker.Worker{
		ID:    1,
		Redis: RedisClient,
	}
	w2 := worker.Worker{
		ID:    2,
		Redis: RedisClient,
	}
	w3 := worker.Worker{
		ID:    3,
		Redis: RedisClient,
	}
	go s.Start()
	go w1.Start()
	go w2.Start()
	go w3.Start()
	select {}
}
