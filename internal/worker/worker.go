package worker

import (
	"fmt"

	"github.com/Akash5106/distributed-task-queue/internal/queue"
)

type Worker struct {
	ID    int
	Queue *queue.Queue
}

func (w *Worker) Start() {
	for {
		t := <-w.Queue.Tasks
		fmt.Printf("Worker : %v ID : %v and Payload : %v\n", w.ID, t.ID, t.Payload)
	}
}
