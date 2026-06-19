package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Akash5106/distributed-task-queue/internal/queue"
	"github.com/Akash5106/distributed-task-queue/internal/task"
)

type Server struct {
	ID    int
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
	var req task.TaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	t := task.Task{
		ID:      s.ID,
		Payload: req.Payload,
	}
	s.ID++
	s.Queue.Enqueue(t)
	w.WriteHeader(http.StatusCreated)
}
