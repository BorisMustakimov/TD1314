package service

import (
	"fmt"
	"log"
	"time"

	"github.com/BorisMustakimov/TD1314/nextdate"
	"github.com/BorisMustakimov/TD1314/repository"
	"github.com/BorisMustakimov/TD1314/task"
)

type TaskService interface {
	AddTask(task *task.Task) (int64, error)
	TaskDone(id string, now time.Time) error
	GetTasks(search, id string) ([]task.Task, error)
	UpdateTask(task *task.Task) error
	DeleteTask(id string) error
}

type TaskServiceImpl struct {
	Repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{
		Repo: repo,
	}
}

func (s *TaskServiceImpl) AddTask(task *task.Task) (int64, error) {
	now := time.Now()
	var taskDate time.Time

	if task.Date == "" || task.Date == now.Format(nextdate.DateFormat) {
		taskDate = now
		task.Date = now.Format(nextdate.DateFormat)
	} else {
		var err error
		taskDate, err = time.Parse(nextdate.DateFormat, task.Date)
		if err != nil {
			log.Printf("неправильный формат даты: %v", err)
			return 0, fmt.Errorf("неправильный формат")
		}
	}

	if taskDate.Before(now) {
		if task.Repeat == "" || task.Repeat == "d 1" {
			task.Date = now.Format(nextdate.DateFormat)
		} else {
			nextDate, err := nextdate.NextDate(now, taskDate.Format(nextdate.DateFormat), task.Repeat)
			if err != nil {
				return 0, fmt.Errorf("невозможно расчитать дату: %v", err)
			}
			task.Date = nextDate
		}
	}

	id, err := s.Repo.Create(task)
	if err != nil {
		return 0, fmt.Errorf("ошибка сохранения задачи: %v", err)
	}

	return id, nil
}

func (s *TaskServiceImpl) GetTasks(search, id string) ([]task.Task, error) {
	filter := repository.Filter{
		Search: search,
	}

	tasks, err := s.Repo.SearchTasks(filter, id) // передаем id в репозиторий
	if err != nil {

		return nil, fmt.Errorf("не удается получить задачу: %v", err)
	}

	return tasks, nil
}

// редактирование задач
func (s *TaskServiceImpl) UpdateTask(task *task.Task) error {
	if task.ID == "" || task.Title == "" {
		return fmt.Errorf("ID или title пустое")
	}

	now := time.Now()
	var taskDate time.Time

	if task.Date == "" || task.Date == now.Format(nextdate.DateFormat) {
		taskDate = now
		task.Date = now.Format(nextdate.DateFormat)
	} else {
		var err error
		taskDate, err = time.Parse(nextdate.DateFormat, task.Date)
		if err != nil {
			return fmt.Errorf("неверный формат")
		}
	}

	if taskDate.Before(now) {
		if task.Repeat == "" || task.Repeat == "d 1" {
			task.Date = now.Format(nextdate.DateFormat)
		} else {
			nextDate, err := nextdate.NextDate(now, taskDate.Format(nextdate.DateFormat), task.Repeat)
			if err != nil {
				return fmt.Errorf("невозможно расчитать дату: %v", err)
			}
			task.Date = nextDate
		}
	}
	err := s.Repo.UpdateTask(task)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %v", err)
	}

	return nil
}

// Удаление задачи
func (s *TaskServiceImpl) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("нет ID")
	}

	err := s.Repo.Delete(id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %v", err)
	}

	return nil
}

// Выполнение задачи
func (s *TaskServiceImpl) TaskDone(id string, now time.Time) error {
	filter := repository.Filter{ID: []string{id}} // Используем фильтр для поиска по ID
	tasks, err := s.Repo.SearchTasks(filter, id)
	if err != nil {
		return fmt.Errorf("неудалось получить задание: %v", err)
	}
	if len(tasks) == 0 {
		return fmt.Errorf("задание не найдено")
	}

	task := tasks[0] // Получаем первую задачу (если она найдена)

	// Если задача не имеет повторений, удалить её
	if task.Repeat == "" {
		if err := s.Repo.Delete(task.ID); err != nil {
			return fmt.Errorf("ошибка удаления задания: %v", err)
		}
	} else {
		// Если задача имеет повторения, вычислить следующую дату
		nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("невозможно расчитать дату: %v", err)
		}

		// Обновляем дату задачи
		task.Date = nextDate
		if err := s.Repo.UpdateTask(&task); err != nil {
			return fmt.Errorf("failed to update task date: %v", err)
		}
	}

	return nil
}
