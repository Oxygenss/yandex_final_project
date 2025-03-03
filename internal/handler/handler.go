package handler

import (
	"net/http"

	"github.com/Oxygenss/yandex_final_project/internal/config"
	"github.com/Oxygenss/yandex_final_project/internal/service"
)

type Task interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	DoneTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	EditTask(w http.ResponseWriter, r *http.Request)
	GetTaskByID(w http.ResponseWriter, r *http.Request)
	GetTasks(w http.ResponseWriter, r *http.Request)
	AddTask(w http.ResponseWriter, r *http.Request)
	NextDateHandler(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	Task
}

func NewHandler(service service.Service, cfg config.Config) *Handler {
	return &Handler{Task: NewTaskHandler(service, cfg)}
}
