package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Oxygenss/yandex_final_project/internal/models"
	"github.com/Oxygenss/yandex_final_project/internal/repository"
)

const (
	DateFormat = "20060102"
)

type TaskService struct {
	repository repository.Repository
}

func NewTaskService(repository repository.Repository) *TaskService {
	return &TaskService{repository: repository}
}

// Title - обязательное поле
// Если date пустая или не указанная, то берется сегодняшнее число
// Date - в формате 20060101 (парсится с помощью  time.Parse())

// Если date < now, то
// - Если repeat пустой или не указан, то берется сегодняшнее число
// - Если repeat указан, то с помощью nextDate считаем дату, которая больше сегодняшней
func (s *TaskService) AddTask(task models.Task) (int64, error) {
	if task.Title == "" {
		return 0, fmt.Errorf("title is required")
	}

	now := time.Now()
	nowFormatted := now.Format(DateFormat)
	nowDate, _ := time.Parse(DateFormat, nowFormatted)

	if task.Date == "" {
		task.Date = nowFormatted
	} else {
		parsedDate, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			return 0, fmt.Errorf("invalid date format. Expected format is YYYYMMDD: %w", err)
		}

		if parsedDate.Equal(nowDate) {
			task.Date = parsedDate.Format(DateFormat)
		} else if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = nowFormatted
			} else {
				nextDate, err := s.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return 0, fmt.Errorf("invalid repeat format or error calculating next date: %w", err)
				}
				task.Date = nextDate
			}
		}
	}

	id, err := s.repository.AddTask(task)
	if err != nil {
		return 0, fmt.Errorf("failed to add task to repository: %w", err)
	}

	return id, nil
}

func (s *TaskService) EditTask(task models.Task) error {
	_, err := strconv.Atoi(task.ID)
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}

	if task.Title == "" {
		return fmt.Errorf("title is required")
	}

	now := time.Now()
	nowFormatted := now.Format(DateFormat)
	nowDate, _ := time.Parse(DateFormat, nowFormatted)

	if task.Date == "" {
		task.Date = nowFormatted
	} else {
		parsedDate, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			return fmt.Errorf("invalid date format. Expected format is YYYYMMDD: %w", err)
		}

		if parsedDate.Equal(nowDate) {
			task.Date = parsedDate.Format(DateFormat)
		} else if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = nowFormatted
			} else {
				nextDate, err := s.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return fmt.Errorf("invalid repeat format or error calculating next date: %w", err)
				}
				task.Date = nextDate
			}
		}
	}

	return s.repository.EditTask(task)
}

func (s *TaskService) DeleteTask(id string) error {
	return s.repository.DeleteByID(id)
}

func (s *TaskService) DoneTask(id string) error {
	task, err := s.GetTaskByID(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		err = s.repository.DeleteByID(id)
		if err != nil {
			return err
		}
	} else {
		nextDate, err := s.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = nextDate

		err = s.repository.EditTask(task)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TaskService) GetTaskByID(id string) (models.Task, error) {
	return s.repository.GetTaskByID(id)
}

func (s *TaskService) GetTasks() ([]models.Task, error) {
	return s.repository.GetTasks()
}

func (s *TaskService) SearchTasks(search string) ([]models.Task, error) {
	time, err := time.Parse("02.01.2006", search)
	if err == nil {
		dateFormatted := time.Format(DateFormat)
		return s.repository.SearchTasksByDate(dateFormatted)
	}

	return s.repository.SearchTasksByString(search)
}

func (s *TaskService) NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat rule is missing")
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid time format: %w", err)
	}

	if strings.HasPrefix(repeat, "d ") {
		daysStr := strings.TrimPrefix(repeat, "d ")

		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return "", fmt.Errorf("invalid day interval: %w", err)
		}

		if days <= 0 || days > 400 {
			return "", fmt.Errorf("day interval must be between 1 and 400: %w", err)
		}

		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				return date.Format(DateFormat), nil
			}
		}
	} else if repeat == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(DateFormat), nil
			}
		}
	} else {
		return "", fmt.Errorf("unsupported repeat format")
	}
}
