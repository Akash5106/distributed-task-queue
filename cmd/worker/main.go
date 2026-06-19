package main

import (
	"fmt"

	"github.com/Akash5106/distributed-task-queue/internal/queue"
	"github.com/Akash5106/distributed-task-queue/internal/task"
	"github.com/Akash5106/distributed-task-queue/internal/worker"
)

func main() {
	q := queue.NewQueue(100)
	w1 := worker.Worker{
		ID:    1,
		Queue: q,
	}
	w2 := worker.Worker{
		ID:    2,
		Queue: q,
	}
	w3 := worker.Worker{
		ID:    3,
		Queue: q,
	}
	t := task.Task{
		ID:      1,
		Payload: "send email",
	}
	t2 := task.Task{
		ID:      2,
		Payload: "send message",
	}
	t3 := task.Task{
		ID:      3,
		Payload: "open youtube",
	}
	t4 := task.Task{
		ID:      4,
		Payload: "open google",
	}
	q.Enqueue(t)
	q.Enqueue(t2)
	q.Enqueue(t3)
	q.Enqueue(t4)
	go w1.Start()
	go w2.Start()
	go w3.Start()
	select {}
	fmt.Println("Main ends")
}
