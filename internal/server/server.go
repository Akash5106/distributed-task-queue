package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Akash5106/distributed-task-queue/internal/storage"
	"github.com/Akash5106/distributed-task-queue/internal/task"
)

type Server struct {
	Redis *storage.RedisClient
}

func NewServer(redis *storage.RedisClient) *Server {
	return &Server{
		Redis: redis,
	}
}

func (s *Server) Start() {
	fmt.Println("Server listening on port : 8080")
	http.HandleFunc("/tasks", s.HandleTasks)
	http.HandleFunc("/tasks/", s.GetTask)
	http.HandleFunc("/dead-tasks", s.GetDeadTasks)
	http.HandleFunc("/dead-tasks/", s.RetryDeadTask)
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
	id, err := s.Redis.GenerateID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := task.Task{
		ID:      id,
		Payload: req.Payload,
		Status:  task.Pending,
		Retries: 0,
	}
	err = s.Redis.SaveTask(r.Context(), t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.Redis.PushTask(r.Context(), t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (s *Server) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	t, err := s.Redis.GetTask(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (s *Server) GetDeadTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	tasks, err := s.Redis.GetDeadTasks(
		r.Context(),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	json.NewEncoder(w).Encode(tasks)
}

func (s *Server) RetryDeadTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		w.WriteHeader(
			http.StatusMethodNotAllowed,
		)
		return
	}
	parts := strings.Split(
		r.URL.Path,
		"/",
	)
	if len(parts) < 4 {
		w.WriteHeader(
			http.StatusBadRequest,
		)
		return
	}
	id, err := strconv.Atoi(
		parts[2],
	)
	if err != nil {
		w.WriteHeader(
			http.StatusBadRequest,
		)
		return
	}
	t, err := s.Redis.GetTask(
		r.Context(),
		id,
	)
	if err != nil {
		w.WriteHeader(
			http.StatusNotFound,
		)
		return
	}
	err = s.Redis.RetryDeadTask(
		r.Context(),
		t,
	)
	if err != nil {
		w.WriteHeader(
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(
		http.StatusOK,
	)
}
