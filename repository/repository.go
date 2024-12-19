package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/BorisMustakimov/TD1314/nextdate"
	"github.com/BorisMustakimov/TD1314/task"
	"github.com/jmoiron/sqlx"
)

const Limit = 50

// определяем интерфейс для работы с задачами 
type TaskRepository interface {
	Create(task *task.Task) (int64, error)
	SearchTasks(filter Filter, id string) ([]task.Task, error)
	UpdateTask(task *task.Task) error
	Delete(id string) error
}

type TaskRepo struct {
	db *sqlx.DB
}

// используем для фильтрации задач
type Filter struct {
	ID     []string
	Search string
	Date   string
}

func NewTaskRepo(db *sqlx.DB) TaskRepository {
	return &TaskRepo{db: db}
}

// создание задачи
func (r *TaskRepo) Create(task *task.Task) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?,?,?,?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Получение списка задач с фильтрацией
func (r *TaskRepo) SearchTasks(filter Filter, id string) ([]task.Task, error) {
	var tasks []task.Task

	// Начальное условие
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1"
	var params []interface{}

	// Если передан ID, то ищем по ID
	if id != "" {
		query += " AND id = ?"
		params = append(params, id)
	}

	// Выполняем фильтрацию по дате или заголовку/комментарию
	if filter.Search != "" {
		parsedDate, err := time.Parse("02.01.2006", filter.Search)
		if err == nil {
			query += " AND date = ?"
			params = append(params, parsedDate.Format(nextdate.DateFormat))
		} else {
			query += " AND (LOWER(title) LIKE LOWER(?) OR LOWER(comment) LIKE LOWER(?))"
			search := "%" + filter.Search + "%"
			params = append(params, search, search)
		}
	}

	// Добавляем сортировку и лимит
	query += " ORDER BY date LIMIT ?"
	params = append(params, Limit)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task task.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Если ищем по ID, то возвращаем ошибку, если задача не найдена
	if id != "" && len(tasks) == 0 {
		return nil, fmt.Errorf("задача не найдена")
	}

	return tasks, nil
}

// Обновление задачи
func (r *TaskRepo) UpdateTask(task *task.Task) error {
	res, err := r.db.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Удаление задачи
func (r *TaskRepo) Delete(id string) error {
	res, err := r.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
