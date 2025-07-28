package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"image_service/archiver"
	"image_service/internal/model"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const tasks = "data/tasks.json"

var allowedExtensions = []string{".pdf", ".jpeg", ".jpg"}

type TaskService struct {
	tasks    map[string]*model.Task
	mu       sync.RWMutex
	archiver *archiver.ArchiverService
}

func NewTaskService(arch *archiver.ArchiverService) *TaskService {
	s := &TaskService{
		tasks:    make(map[string]*model.Task),
		archiver: arch,
	}
	_ = s.loadTasks()
	return s
}

func (s *TaskService) CreateTask() (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isServerBusy() {
		return nil, errors.New("server busy")
	}

	id := uuid.New().String()

	task := &model.Task{
		ID:     id,
		Status: model.StatusPending,
		Files:  []string{},
	}

	s.tasks[id] = task
	_ = s.saveTasks()

	return task, nil
}

func (s *TaskService) saveTasks() error {
	if err := os.MkdirAll(filepath.Dir(tasks), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(tasks, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Tasks saved:", len(s.tasks))
	return nil
}

func (s *TaskService) loadTasks() error {
	data, err := os.ReadFile(tasks)

	if err != nil {
		if os.IsNotExist(err) {
			s.tasks = make(map[string]*model.Task)
			return nil
		}
		return err
	}

	s.tasks = make(map[string]*model.Task)
	err = json.Unmarshal(data, &s.tasks)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tasks json: %w", err)
	}

	fmt.Println("Tasks loaded")
	return nil
}
func (s *TaskService) AddFilesInTask(id string, files []string) (*model.Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("failed to find task %s", id)
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(task.Files) >= 3 {
		return nil, errors.New("task has maximum number of files")
	}

	if len(task.Files)+len(files) > 3 {
		return nil, errors.New("adding these files would exceed max of 3")
	}

	filteredFiles, err := filterFiles(files)
	if err != nil {
		return nil, err
	}

	task.Files = append(task.Files, filteredFiles...)

	if len(task.Files) == 3 {
		task.Status = model.StatusRunning
		go func(task *model.Task) {
			s.archiver.Archive(task)
			_ = s.saveTasks()
		}(task)
	} else {
		task.Status = model.StatusPending
	}

	_ = s.saveTasks()

	return task, nil
}

func (s *TaskService) GetTaskStatus(id string) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("failed to find task %s", id)
	}
	return task, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func isAllowedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(strings.Split(filename, "?")[0]))
	return contains(allowedExtensions, ext)
}

func filterFiles(files []string) ([]string, error) {
	var filtered []string
	for _, f := range files {
		if !isAllowedExtension(f) {
			return nil, fmt.Errorf("file extension %s not allowed", filepath.Ext(f))
		}
		filtered = append(filtered, f)
	}
	return filtered, nil
}
func (s *TaskService) isServerBusy() bool {
	active := 0
	for _, task := range s.tasks {
		if task.Status == model.StatusRunning || task.Status == model.StatusPending {
			active++
		}

		if active >= 3 {
			return true
		}
	}
	return false
}
