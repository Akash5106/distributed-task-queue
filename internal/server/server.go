package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Akash5106/distributed-task-queue/internal/queue"
	"github.com/Akash5106/distributed-task-queue/internal/task"
)

type Server struct {
	Queue *queue.Queue
}

func NewServer(q *queue.Queue) *Server {
	return &Server{
		Queue: q,
	}
}

func (s *Server) Start() {
	fmt.Println("Server listening on port : 8080")
	http.HandleFunc("/tasks", s.HandleTasks)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) HandleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	t := task.Task{
		ID:      1,
		Payload: "PAKSHIGOD",
	}
	s.Queue.Enqueue(t)
	w.WriteHeader(http.StatusCreated)
}
