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
	const workerCount = 10

	for i := 1; i <= workerCount; i++ {
		w := worker.Worker{
			ID:    i,
			Redis: RedisClient,
		}
		go w.Start()
	}
	go s.Start()
	select {}
}
