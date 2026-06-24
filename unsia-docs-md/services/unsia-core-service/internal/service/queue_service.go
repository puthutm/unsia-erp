package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/google/uuid"
)

const (
	QueueDefault = "default"
	QueueCritical = "critical"
	QueueLow     = "low"
	QueueEmail   = "email"
	QueueNotification = "notification"
)

// Task represents a background task
type Task struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	Queue     string    `json:"queue"`
	Retry     int       `json:"retry"`
	Retried   int       `json:"retried"`
	Error     string    `json:"error,omitempty"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Task) TableName() string {
	return "tasks"
}

// TaskService handles task management
type TaskService struct {
	asynqScheduler *asynq.Scheduler
	asynqQueue     *asynq.Queue
}

func NewTaskService(redisAddr string) (*TaskService, error) {
	r := asynq.RedisClientOpt{Addr: redisAddr}
	scheduler := asynq.NewScheduler(r, nil)
	queue := asynq.NewQueue(r)

	return &TaskService{
		asynqScheduler: scheduler,
		asynqQueue:    queue,
	}, nil
}

// EnqueueTask enqueues a task
func (s *TaskService) EnqueueTask(taskType string, payload interface{}, queue string, opts ...TaskOption) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(taskType, payloadBytes)

	options := defaultTaskOptions
	for _, opt := range opts {
		opt(&options)
	}

	// Enqueue with options
	if options.Delay > 0 {
		return s.asynqScheduler.Enqueue(task,
			asynq.ProcessIn(options.Delay),
			asynq.Queue(queue),
			asynq.MaxRetry(options.Retry),
		)
	}

	return s.asynqQueue.Enqueue(task,
		asynq.Queue(queue),
		asynq.MaxRetry(options.Retry),
	)
}

// TaskOption is a functional option
type TaskOption func(*TaskOptions)

type TaskOptions struct {
	Delay time.Duration
	Retry int
}

var defaultTaskOptions = TaskOptions{
	Delay: 0,
	Retry: 3,
}

// WithDelay sets delay
func WithDelay(delay time.Duration) TaskOption {
	return func(o *TaskOptions) {
		o.Delay = delay
	}
}

// WithRetry sets retry count
func WithRetry(retry int) TaskOption {
	return func(o *TaskOptions) {
		o.Retry = retry
	}
}

// EmailTaskPayload represents email task payload
type EmailTaskPayload struct {
	To      []string `json:"to"`
	Subject string  `json:"subject"`
	Body    string  `json:"body"`
	Type    string  `json:"type"` // text/html
}

// SendEmailTask creates a send email task
func (s *TaskService) SendEmailTask(payload EmailTaskPayload) error {
	return s.EnqueueTask("send_email", payload, QueueEmail)
}

// NotificationTaskPayload represents notification task payload
type NotificationTaskPayload struct {
	UserID  string `json:"user_id"`
	Title  string `json:"title"`
	Message string `json:"message"`
	Type   string `json:"type"`
}

// SendNotificationTask creates a notification task
func (s *TaskService) SendNotificationTask(payload NotificationTaskPayload) error {
	return s.EnqueueTask("send_notification", payload, QueueNotification)
}

// CleanupTaskPayload represents cleanup task payload
type CleanupTaskPayload struct {
	Type string `json:"type"`
	Days int    `json:"days"`
}

// ScheduleCleanupTask schedules a cleanup task
func (s *TaskService) ScheduleCleanupTask(payload CleanupTaskPayload) error {
	return s.EnqueueTask("cleanup", payload, QueueLow, WithDelay(24*time.Hour))
}

// InMemoryTaskQueue provides in-memory queue for simple use cases
type InMemoryTaskQueue struct {
	tasks   []Task
	mu      sync.Mutex
	handlers map[string]TaskHandler
}

type TaskHandler func(ctx context.Context, task Task) error

// NewInMemoryTaskQueue creates a new in-memory task queue
func NewInMemoryTaskQueue() *InMemoryTaskQueue {
	return &InMemoryTaskQueue{
		tasks:   make([]Task, 0),
		handlers: make(map[string]TaskHandler),
	}
}

// RegisterHandler registers a task handler
func (q *InMemoryTaskQueue) RegisterHandler(taskType string, handler TaskHandler) {
	q.handlers[taskType] = handler
}

// Enqueue adds a task to the queue
func (q *InMemoryTaskQueue) Enqueue(taskType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Payload:   string(payloadBytes),
		Queue:    QueueDefault,
		Retry:    3,
		Retried:  0,
		CreatedAt: time.Now(),
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, task)

	return nil
}

// Dequeue removes and returns a task
func (q *InMemoryTaskQueue) Dequeue() (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tasks) == 0 {
		return nil, nil
	}

	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return &task, nil
}

// Process processes all tasks in the queue
func (q *InMemoryTaskQueue) Process(ctx context.Context) error {
	for {
		task, err := q.Dequeue()
		if err != nil {
			return err
		}
		if task == nil {
			break
		}

		handler, exists := q.handlers[task.Type]
		if !exists {
			fmt.Printf("No handler for task type: %s\n", task.Type)
			continue
		}

		if err := handler(ctx, *task); err != nil {
			fmt.Printf("Task error: %v\n", err)
			// Could implement retry logic here
		}
	}

	return nil
}

// ScheduleTask schedules a task for later
type ScheduledTask struct {
	Task      Task
	ExecuteAt time.Time
}

type SchedulerService struct {
	scheduledTasks []ScheduledTask
	mu           sync.Mutex
}

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{
		scheduledTasks: make([]ScheduledTask, 0),
	}
}

// Schedule schedules a task for execution
func (s *SchedulerService) Schedule(task Task, executeAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scheduledTasks = append(s.scheduledTasks, ScheduledTask{
		Task:      task,
		ExecuteAt: executeAt,
	})
}

// GetDueTasks returns tasks that are due
func (s *SchedulerService) GetDueTasks() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	var dueTasks []Task
	var remaining []ScheduledTask

	for _, st := range s.scheduledTasks {
		if st.ExecuteAt.Before(now) || st.ExecuteAt.Equal(now) {
			dueTasks = append(dueTasks, st.Task)
		} else {
			remaining = append(remaining, st)
		}
	}

	s.scheduledTasks = remaining
	return dueTasks
}
