package storage

import (
	"context"
	"encoding/json"
	"fmt"

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

func (r *RedisClient) SaveTask(ctx context.Context, t task.Task) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("task:%v", t.ID)
	res := r.Client.Set(ctx, key, data, 0)
	return res.Err()
}

func (r *RedisClient) UpdateTask(ctx context.Context, t task.Task) error {
	key := fmt.Sprintf("task:%v", t.ID)
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	res := r.Client.Set(ctx, key, data, 0)
	return res.Err()
}

func (r *RedisClient) GetTask(ctx context.Context, id int) (task.Task, error) {
	key := fmt.Sprintf("task:%v", id)
	res := r.Client.Get(ctx, key)
	if res.Err() != nil {
		return task.Task{}, res.Err()
	}
	var t task.Task
	data := json.Unmarshal([]byte(res.Val()), &t)
	if data != nil {
		return task.Task{}, data
	}
	return t, nil
}

func (r *RedisClient) PushDeadTask(ctx context.Context, t task.Task) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	res := r.Client.LPush(ctx, "dead_tasks", string(data))
	return res.Err()
}

func (r *RedisClient) GetDeadTasks(ctx context.Context) ([]task.Task, error) {
	res := r.Client.LRange(ctx, "dead_tasks", 0, -1)
	if res.Err() != nil {
		return nil, res.Err()
	}
	var tasks []task.Task
	for _, item := range res.Val() {
		var t task.Task
		err := json.Unmarshal(
			[]byte(item),
			&t,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *RedisClient) RemoveDeadTask(
	ctx context.Context,
	t task.Task,
) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	res := r.Client.LRem(
		ctx,
		"dead_tasks",
		1,
		string(data),
	)
	return res.Err()
}

func (r *RedisClient) RetryDeadTask(
	ctx context.Context,
	t task.Task,
) error {

	err := r.RemoveDeadTask(
		ctx,
		t,
	)
	if err != nil {
		return err
	}
	t.Status = task.Pending
	t.Retries = 0
	err = r.UpdateTask(
		ctx,
		t,
	)
	if err != nil {
		return err
	}
	err = r.PushTask(
		ctx,
		t,
	)
	if err != nil {
		return err
	}
	return nil
}
