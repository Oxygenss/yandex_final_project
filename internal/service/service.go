package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Oxygenss/yandex_go_final_project/internal/models"
	repository "github.com/Oxygenss/yandex_go_final_project/internal/storage"
)

type Service struct {
	repository repository.Repository
}

func New(repository repository.Repository) *Service {
	return &Service{repository: repository}
}

// Title - обязательное поле
// Если date пустая или не указанная, то берется сегодняшнее число
// Date - в формате 20060101 (парсится с помощью  time.Parse())

// Если date < now, то
// - Если repeat пустой или не указан, то берется сегодняшнее число
// - Если repeat указан, то с помощью nextDate считаем дату, которая больше сегодняшней
func (s *Service) AddTask(task models.Task) (int64, error) {
	if task.Title == "" {
		return 0, fmt.Errorf("title is required")
	}

	now := time.Now()
	nowFormatted := now.Format("20060102")
	nowDate, _ := time.Parse("20060102", nowFormatted)

	if task.Date == "" {
		task.Date = nowFormatted
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			return 0, fmt.Errorf("invalid date format. Expected format is YYYYMMDD: %w", err)
		}

		if parsedDate.Equal(nowDate) {
			task.Date = parsedDate.Format("20060102")
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

func (s *Service) EditTask(task models.Task) error {
	_, err := strconv.Atoi(task.ID)
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}

	if task.Title == "" {
		return fmt.Errorf("title is required")
	}

	now := time.Now()
	nowFormatted := now.Format("20060102")
	nowDate, _ := time.Parse("20060102", nowFormatted)

	if task.Date == "" {
		task.Date = nowFormatted
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			return fmt.Errorf("invalid date format. Expected format is YYYYMMDD: %w", err)
		}

		if parsedDate.Equal(nowDate) {
			task.Date = parsedDate.Format("20060102")
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

func (s *Service) DeleteTask(id string) error {
	return s.repository.DeleteByID(id)
}

func (s *Service) DoneTask(id string) error {
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

func (s *Service) GetTaskByID(id string) (models.Task, error) {
	return s.repository.GetTaskByID(id)
}

func (s *Service) GetTasks() ([]models.Task, error) {
	return s.repository.GetTasks()
}

func (s *Service) SearchTasks(search string) ([]models.Task, error) {
    time, err := time.Parse("02.01.2006", search)
    if err == nil {
        dateFormatted := time.Format("20060102")
        return s.repository.SearchTasksByDate(dateFormatted)
    }

    return s.repository.SearchTasksByString(search)
}

func (s *Service) NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("нет правила повторения")
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат времени: %w", err)
	}

	if strings.HasPrefix(repeat, "d ") {
		daysStr := strings.TrimPrefix(repeat, "d ")

		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return "", fmt.Errorf("неверный интервал дня: %w", err)
		}

		if days <= 0 || days > 400 {
			return "", fmt.Errorf("интервал дня может быть в пределах от 1 до 400: %w", err)
		}

		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				return date.Format("20060102"), nil
			}
		}
	} else if repeat == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format("20060102"), nil
			}
		}
	} else {
		return "", fmt.Errorf("неподдерживаемый формат повторения")
	}
}
