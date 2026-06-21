package task

type Task struct {
	ID      int    `json:"id"`
	Payload string `json:"payload"`
	Status  string `json:"status"`
}

type TaskRequest struct {
	Payload string `json:"payload"`
}

const (
	Pending   = "pending"
	Running   = "running"
	Completed = "completed"
	Failed    = "failed"
)
