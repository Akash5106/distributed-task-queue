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
	t := <-w.Queue.Tasks
	fmt.Printf("ID : %v and Payload : %v", t.ID, t.Payload)
}
