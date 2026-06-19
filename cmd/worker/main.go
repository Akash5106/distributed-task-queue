package main

import (
	"github.com/Akash5106/distributed-task-queue/internal/queue"
	"github.com/Akash5106/distributed-task-queue/internal/server"
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
	s := server.NewServer(q)
	go w1.Start()
	go w2.Start()
	go w3.Start()
	s.Start()
}
