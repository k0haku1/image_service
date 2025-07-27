package main

import (
	"image_service/internal/handlers"
	"image_service/internal/service"
	"net/http"
)

func InitRouter() *http.ServeMux {
	taskService := service.NewTaskService()
	taskHandler := handlers.NewTaskHandler(taskService)

	router := http.NewServeMux()

	router.HandleFunc("/task/create", taskHandler.CreateTask)
	router.HandleFunc("/task", taskHandler.GetTask)
	return router
}
