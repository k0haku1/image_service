package main

import (
	"image_service/archiver"
	"image_service/internal/handlers"
	"image_service/internal/service"
	"net/http"
)

func InitRouter() *http.ServeMux {
	archiverService := archiver.NewArchiverService()
	taskService := service.NewTaskService(archiverService)
	taskHandler := handlers.NewTaskHandler(taskService)

	router := http.NewServeMux()

	router.HandleFunc("/task/create", taskHandler.CreateTask)
	router.HandleFunc("/task/add", taskHandler.AddFiles)
	router.HandleFunc("/task/status", taskHandler.GetStatus)

	return router
}
