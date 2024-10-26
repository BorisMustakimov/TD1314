package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/database"
	"go_final_project/models"
	"go_final_project/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	// Парсинг дат
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Неверный формат 'now'", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		http.Error(w, "Неверный формат 'date'", http.StatusBadRequest)
		return
	}

	nextDate, err := utils.NextDate(now, date, repeatStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := nextDate.Format("20060102")
	if response == "00010101" {
		response = ""
	}
	w.Write([]byte(response))
}

func TaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		handlePostTask(w, r, db)
	case http.MethodGet:
		handleGetTask(w, r, db)
	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

func handlePostTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}
	if !isValidRepeatFormat(task.Repeat) {
		http.Error(w, `{"error":"Неверный формат правила повторения"}`, http.StatusBadRequest)
		return
	}
	now := time.Now()

	var taskDate time.Time
	var err error
	if task.Date == "" || task.Date == now.Format("20060102") {
		taskDate = now
	} else {
		taskDate, err = time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Дата представлена в неверном формате"}`, http.StatusBadRequest)
			return
		}
	}

	if taskDate.Before(now) {
		if task.Repeat == "" {
			taskDate = now
		} else {
			nextDate, err := utils.NextDate(now, taskDate, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Неверный формат правила повторения"}`, http.StatusBadRequest)
				return
			}
			taskDate = nextDate
		}
	}

	taskID, err := database.InsertTask(db, taskDate, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"id": strconv.Itoa(taskID)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func isValidRepeatFormat(repeat string) bool {
	// Проверка на "y"
	if repeat == "y" {
		return true
	}

	// Проверка на "d <число>"
	parts := strings.Split(repeat, " ")
	if len(parts) == 2 && parts[0] == "d" {
		if _, err := strconv.Atoi(parts[1]); err == nil {
			return true
		}
	}
	if repeat == "" {
		return true
	}
	return false
}

func GetTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50")
	if err != nil {
		http.Error(w, `{"error": "Ошибка выполнения запроса"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		var date time.Time
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, `{"error": "Ошибка чтения данных"}`, http.StatusInternalServerError)
			return
		}

		task.Date = date.Format("20060102")
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	// Проверка ошибок после завершения перебора
	if err := rows.Err(); err != nil {
		http.Error(w, `{"error": "Ошибка при переборе результатов"}`, http.StatusInternalServerError)
		return
	}

	// Создание ответа

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(map[string]interface{}{"tasks": tasks})
	if err != nil {
		http.Error(w, `{"error": "Ошибка при переборе результатов"}`, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}

func handleGetTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := database.GetTaskByID(db, id)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":      strconv.Itoa(task.ID),
		"date":    task.Date.Format("20060102"),
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
