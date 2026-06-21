package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/Akash5106/distributed-task-queue/internal/storage"
	"github.com/Akash5106/distributed-task-queue/internal/task"
)

type Worker struct {
	Redis *storage.RedisClient
	ID    int
}

func (w *Worker) Start() {
	for {
		t, err := w.Redis.PopTask(context.Background())
		if err != nil {
			panic(err)
		}
		t.Status = task.Running
		res := w.Redis.UpdateTask(context.Background(), t)
		if res != nil {
			panic(res)
		}
		fmt.Printf("Worker %v Task %v Payload %v Status %v\n", w.ID, t.ID, t.Payload, t.Status)
		time.Sleep(5 * time.Second)
		t.Status = task.Completed
		res = w.Redis.UpdateTask(context.Background(), t)
		if res != nil {
			panic(res)
		}
		fmt.Printf("Worker %v Task %v Payload %v Status %v\n", w.ID, t.ID, t.Payload, t.Status)
	}
}
