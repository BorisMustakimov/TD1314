package sqltable

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BorisMustakimov/TD1314/config"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Sql_table(cfg *config.Config) (*sqlx.DB, error) {

	dbPath := cfg.DBFile
	if dbPath == "" {
		appPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("ошибка пулучения пути к БД")
		}
		dbPath = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	// Проверка существования базы данных
	_, err := os.Stat(dbPath)
	install := os.IsNotExist(err)

	// Подключение к базе данных
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания таблицы")
	}

	// Если база данных не существует, создаём таблицу и индекс
	if install {
		_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK(length(repeat) <= 128)
	);
	CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания таблицы")
		}
		fmt.Println("таблица создана успешно")
	}

	return db, nil
}
