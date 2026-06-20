package task

type Task struct {
	ID      int    `json:"id"`
	Payload string `json:"payload"`
}

type TaskRequest struct {
	Payload string `json:"payload"`
}
