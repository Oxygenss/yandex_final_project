package repository

import (
	"github.com/Oxygenss/yandex_final_project/internal/models"
	"github.com/Oxygenss/yandex_final_project/internal/repository/migrations"
	"github.com/Oxygenss/yandex_final_project/internal/repository/sqlite"
)

type Repository interface {
	AddTask(task models.Task) (int64, error)
	GetTaskByID(id string) (models.Task, error)
	GetTasks() ([]models.Task, error)
	SearchTasksByString(search string) ([]models.Task, error)
	SearchTasksByDate(date string) ([]models.Task, error)
	EditTask(task models.Task) error
	DeleteByID(id string) error
}

func New(pathDB string) (Repository, error) {
	db, err := sqlite.NewSQLiteDB(pathDB)
	if err != nil {
		return nil, err
	}

	err = migrations.Migrations(db, pathDB)
	if err != nil {
		return nil, err
	}

	repo := sqlite.New(db)

	return repo, nil
}
