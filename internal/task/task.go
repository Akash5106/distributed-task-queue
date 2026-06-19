package task

type Task struct {
	ID      int
	Payload string
}

type TaskRequest struct {
	Payload string `json:"Payload"`
}
