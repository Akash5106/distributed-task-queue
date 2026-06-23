# Distributed Task Queue

A fault-tolerant distributed task queue built in Go using Redis.

Designed to execute background jobs reliably across multiple workers while handling retries, worker crashes, task recovery, and failed task replay.

## Highlights

- Multi-worker concurrent task processing
- Automatic retry mechanism
- Dead Letter Queue (DLQ) support
- Worker crash recovery using processing timeouts
- Dockerized deployment with Docker Compose
- Redis-backed task storage and queue management

---

## Features

### Core Execution

- REST API for task submission
- Redis-backed task queue
- Multiple concurrent workers
- Task status tracking
- Persistent task storage

### Reliability

- Automatic task retries
- Dead Letter Queue (DLQ)
- Failed task replay endpoint
- Processing queue ownership
- Visibility timeout style recovery
- Worker crash recovery
- At-least-once task delivery

### Deployment

- Dockerized application
- Docker Compose support
- Redis container orchestration
- One-command startup

---

## Architecture

```text
                    +---------+
                    | Client  |
                    +---------+
                         |
                    POST /tasks
                         |
                         v

                +----------------+
                |   HTTP Server  |
                +----------------+
                         |
                         v

                +----------------+
                |     Redis      |
                +----------------+
                 |      |      |
                 |      |      |
             tasks  processing  dead_tasks
                 |
                 |
      +----------+----------+
      |          |          |
      v          v          v

 +---------+ +---------+ +---------+
 | Worker  | | Worker  | | Worker  |
 |    1    | |    2    | |    3    |
 +---------+ +---------+ +---------+

                 |
                 |
           Process Tasks
                 |
         +-------+-------+
         |               |
         v               v
     Completed       Failed
                         |
                         v
                        DLQ
```

---

## Task Lifecycle

```text
Pending
   |
   v
Processing
   |
   +----------------+
   |                |
   v                |
Completed      Worker Crash
                    |
                    v
            Processing Timeout
                    |
                    v
             Task Recovery
                    |
                    v
                Requeued
```

### Retry Flow

```text
Pending
   |
   v
Processing
   |
   v
Failure
   |
   v
Retry 1
   |
   v
Retry 2
   |
   v
Retry 3
   |
   +-------------+
   |             |
   v             v
Success        Failed
                  |
                  v
                 DLQ
```

---

## Crash Recovery

When a worker claims a task:

```text
tasks
  |
  v
processing
```

A timestamp is recorded:

```text
processing:<task_id>
```

The recovery service periodically scans the processing queue.

If a task remains in processing beyond the timeout window:

```text
Worker Crash
      |
      v
Task Stuck In Processing
      |
      v
Timeout Exceeded
      |
      v
Task Recovered
      |
      v
Moved Back To Queue
      |
      v
Another Worker Processes It
```

---

## API Endpoints

### Create Task

```http
POST /tasks
```

Request:

```json
{
  "payload": "send email"
}
```

Response:

```json
{
  "id": 1,
  "payload": "send email",
  "status": "pending",
  "retries": 0
}
```

---

### Get Task

```http
GET /tasks/{id}
```

Example:

```bash
curl http://localhost:8080/tasks/1
```

---

### View Dead Letter Queue

```http
GET /dead-tasks
```

Example:

```bash
curl http://localhost:8080/dead-tasks
```

---

### Retry Failed Task

```http
POST /dead-tasks/{id}/retry
```

Example:

```bash
curl -X POST http://localhost:8080/dead-tasks/1/retry
```

---

## Project Structure

```text
distributed-task-queue/
│
├── cmd/
│   └── worker/
│       └── main.go
│
├── internal/
│   ├── server/
│   │   └── server.go
│   │
│   ├── storage/
│   │   └── redis.go
│   │
│   ├── task/
│   │   └── task.go
│   │
│   └── worker/
│       └── worker.go
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## Running Locally

### Using Docker Compose

```bash
git clone https://github.com/Akash5106/distributed-task-queue.git

cd distributed-task-queue

docker compose up
```

Services started:

```text
Redis      -> localhost:6379
HTTP API   -> localhost:8080
Workers    -> Worker 1, Worker 2, Worker 3
```

---

## Technologies

- Go
- Redis
- Docker
- Docker Compose
- REST APIs

---

## Design Goals

This project was built to explore the core concepts behind distributed task processing systems:

- Reliable task execution
- Fault tolerance
- Retry mechanisms
- Dead Letter Queues
- Worker crash recovery
- Queue-based architectures
- Containerized deployments

---

## Future Improvements

- Metrics endpoint
- Worker health monitoring
- Worker auto-restart
- Task scheduling
- Priority queues
- Prometheus + Grafana monitoring
- Web dashboard

---
