package service

import (
	"time"

	"github.com/Oxygenss/yandex_final_project/internal/models"
	"github.com/Oxygenss/yandex_final_project/internal/repository"
)

type Task interface {
	AddTask(task models.Task) (int64, error)
	EditTask(task models.Task) error
	DeleteTask(id string) error
	DoneTask(id string) error
	GetTaskByID(id string) (models.Task, error)
	GetTasks() ([]models.Task, error)
	SearchTasks(search string) ([]models.Task, error)
	NextDate(now time.Time, dateStr string, repeat string) (string, error)
}

type Service struct {
	Task
}

func NewService(repository repository.Repository) *Service {
	return &Service{Task: NewTaskService(repository)}
}
