package storage

import (
	"context"
	"encoding/json"

	"github.com/Akash5106/distributed-task-queue/internal/task"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &RedisClient{
		Client: client,
	}
}

func (r *RedisClient) Ping() error {
	ctx := context.Background()
	result := r.Client.Ping(ctx)
	return result.Err()
}

func (r *RedisClient) PushTask(ctx context.Context, t task.Task) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	res := r.Client.LPush(ctx, "tasks", string(data))
	return res.Err()
}

func (r *RedisClient) PopTask(ctx context.Context) (task.Task, error) {
	res := r.Client.BRPop(ctx, 0, "tasks")
	if res.Err() != nil {
		return task.Task{}, res.Err()
	}
	data := res.Val()
	var t task.Task
	err := json.Unmarshal(
		[]byte(data[1]),
		&t,
	)
	if err != nil {
		return task.Task{}, err
	}
	return t, nil
}

func (r *RedisClient) GenerateID(ctx context.Context) (int, error) {
	id := r.Client.Incr(ctx, "task_id")
	if id.Err() != nil {
		return -1, id.Err()
	}
	return int(id.Val()), nil
}
