package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Router *chi.Mux
	server *http.Server
	db     *sqlx.DB
}

func NewRouter(version string, serverAddress string, taskService service.TaskService, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Статические файлы из директории web
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", http.StripPrefix("/", fileServer))

	taskHandler := hendlers.NewTaskHandler(taskService)

	// Маршрут для обработки правил повторения
	r.Get("/api/nextdate", nextdate.HandlerNextDate)

	r.Group(func(r chi.Router) {

		r.Get("/api/tasks", taskHandler.GetTasksHandler)
		r.Mount("/api/task", http.HandlerFunc(taskHandler.TaskHandler))
	})

	return r
}

func New() (*App, error) {
	cfg, err := config.MustLoad()
	if err != nil {
		return nil, err
	}

	database, err := sqltable.Sql_table(cfg)
	if err != nil {
		return nil, err
	}

	taskRepo := repository.NewTaskRepo(database)
	taskService := service.NewTaskService(taskRepo)

	router := NewRouter(cfg.Version, cfg.ServerAddress, taskService, cfg)
	if router == nil {
		return nil, fmt.Errorf("ошибка создания роутера")
	}

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	a := &App{
		Router: router,
		server: server,
		db:     database,
	}

	return a, nil
}

func (a *App) Run() error {
	fmt.Printf("Starting app: Listening on %s\n", a.server.Addr)

	// Запуск сервера
	err := a.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("App stopped")
	}

	return err
}
