package service

import (
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
