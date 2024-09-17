package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Zelvalna/go_final_project/config"

	"github.com/Zelvalna/go_final_project/internal/handlers"
	"github.com/Zelvalna/go_final_project/internal/middleware"
	"github.com/Zelvalna/go_final_project/internal/storage"
	"github.com/Zelvalna/go_final_project/model"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Инициализация базы данных
	_, err := storage.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Проверяем, установлен ли пароль в переменной окружения TODO_PASSWORD
	cfg := config.Config{
		TodoPassword: os.Getenv("TODO_PASSWORD"),
		Port:         model.DefPort,
	}
	if cfg.TodoPassword == "" {
		log.Fatal("TODO_PASSWORD environment variable is required")
	}

	// Проверка переменной окружения TODO_PORT для переопределения порта
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) != 0 {
		cfg.Port = envPort
	}

	// Путь к директории с веб-файлами
	webDir := model.WebDir
	fs := http.FileServer(http.Dir(webDir))

	r := chi.NewRouter()

	r.Mount("/", fs)
	r.Get("/api/nextdate", handlers.NextDateHandler)
	r.Post("/api/task", middleware.Auth(handlers.TaskHandler, cfg))
	r.Get("/api/tasks", middleware.Auth(handlers.TaskHandler, cfg))
	r.Get("/api/task", middleware.Auth(handlers.TaskByIdGet, cfg))
	r.Put("/api/task", middleware.Auth(handlers.TaskHandler, cfg))
	r.Post("/api/task/done", middleware.Auth(handlers.TaskDonePost, cfg))
	r.Delete("/api/task", middleware.Auth(handlers.TaskHandler, cfg))
	r.Post("/api/signin", func(w http.ResponseWriter, r *http.Request) { handlers.SingInHandler(w, r, cfg) })

	// Запуск сервера
	log.Printf("Сервер запущен на порту %v", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
