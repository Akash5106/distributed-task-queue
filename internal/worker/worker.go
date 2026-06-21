package worker

import (
	"context"
	"fmt"

	"github.com/Akash5106/distributed-task-queue/internal/storage"
)

type Worker struct {
	Redis *storage.RedisClient
	ID    int
}

func (w *Worker) Start() {
	for {
		data, err := w.Redis.PopTask(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Worker %v Task %v Payload %v\n", w.ID, data.ID, data.Payload)
	}
}
