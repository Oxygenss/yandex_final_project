package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Oxygenss/yandex_go_final_project/internal/models"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddTask(task models.Task) (int64, error) {

	query := `INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?, ?, ?, ?)`

	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func (r *Repository) GetTaskByID(id string) (models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"

	var task models.Task

	res := r.db.QueryRow(query, id)
	err := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, fmt.Errorf("task with id %s not found: %w", id, err)
		}
		return models.Task{}, fmt.Errorf("error executing query: %w", err)
	}

	return task, nil
}

func (r *Repository) GetTasks() ([]models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func (r *Repository) SearchTasksByString(search string) ([]models.Task, error) {
	query := `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date`

	rows, err := r.db.Query(query, "%"+search+"%", "%"+search+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Comment, &task.Date, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func (r *Repository) SearchTasksByDate(date string) ([]models.Task, error) {
	query := `SELECT * FROM scheduler WHERE date = ?`

	rows, err := r.db.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Comment, &task.Date, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func (r *Repository) EditTask(task models.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	result, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to edit task with id %s: %w", task.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get the number of affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", task.ID)
	}

	return nil
}

func (r *Repository) DeleteByID(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get the number of affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}

	return nil
}
