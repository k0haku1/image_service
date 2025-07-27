package service

import (
	"github.com/google/uuid"
	"image_service/internal/model"
	"sync"
)

type TaskService struct {
	tasks map[string]*model.Task
	mu    sync.RWMutex
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[string]*model.Task),
	}
}

func (s *TaskService) CreateTask() *model.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()

	task := &model.Task{
		ID:     id,
		Status: model.StatusPending,
		Files:  []string{},
	}

	s.tasks[id] = task
	return task
}

func (s *TaskService) GetTask(id string) *model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.tasks[id]
}
