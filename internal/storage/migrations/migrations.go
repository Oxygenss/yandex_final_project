package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func Migrations(db *sql.DB, pathDB string) error {

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date INTEGER NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat VARCHAR(128)
	);
	`

	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
	`

	_, err := os.Stat(pathDB)
	dbExists := !os.IsNotExist(err)

	if !dbExists {
		_, err = db.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("ошибка при создании таблицы: %w", err)
		}

		_, err = db.Exec(createIndexSQL)
		if err != nil {
			return fmt.Errorf("ошибка при создании индекса: %w", err)
		}

		log.Println("База данных была создана и успешно инициализирована.")
	} else {
		log.Println("База данных уже существует. Подключение выполнено.")
	}
	return nil
}
