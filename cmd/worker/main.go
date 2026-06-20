package main

import (
	//"github.com/Akash5106/distributed-task-queue/internal/queue"
	//"github.com/Akash5106/distributed-task-queue/internal/server"

	"context"
	"encoding/json"

	"github.com/Akash5106/distributed-task-queue/internal/storage"
	"github.com/Akash5106/distributed-task-queue/internal/task"

	//"github.com/Akash5106/distributed-task-queue/internal/worker"
	"fmt"
)

func main() {
	redisClient := storage.NewRedisClient()
	t := task.Task{
		ID:      1,
		Payload: "send email",
	}
	redisClient.PushTask(context.Background(), t)
	res, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
	data, err := redisClient.PopTask(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v, %v", data.ID, data.Payload)
}
