package main

import (
	"github.com/Akash5106/distributed-task-queue/internal/queue"
	"github.com/Akash5106/distributed-task-queue/internal/task"
	"github.com/Akash5106/distributed-task-queue/internal/worker"
)

func main() {
	q := queue.NewQueue(100)
	w := worker.Worker{
		ID:    1,
		Queue: q,
	}
	t := task.Task{
		ID:      1,
		Payload: "send email",
	}
	q.Enqueue(t)
	w.Start()
}
