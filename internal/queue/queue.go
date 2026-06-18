package queue

import "github.com/Akash5106/distributed-task-queue/internal/task"

type Queue struct {
	Tasks chan task.Task
}

func NewQueue(size int) *Queue {
	return &Queue{
		Tasks: make(chan task.Task, size),
	}
}

func (q *Queue) Enqueue(t task.Task) {
	q.Tasks <- t
}
