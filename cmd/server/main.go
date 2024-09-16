package main

import (
	"github.com/Zelvalna/go_final_project/constans"
	"github.com/Zelvalna/go_final_project/internal/handlers"
	"github.com/Zelvalna/go_final_project/internal/middleware"
	"github.com/Zelvalna/go_final_project/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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
	password := os.Getenv("TODO_PASSWORD")
	log.Println(password)
	if password == "" {
		log.Fatal("TODO_PASSWORD environment variable is required")
	}
	// установка порта по умолчанию
	port := constans.DefPort
	// Проверка переменной окружения TODO_PORT для переопределения порта
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) == 0 {
		envPort = port
	} else {
		port = envPort
	}

	// Путь к директории с веб-файлами
	webDir := constans.WebDir
	fs := http.FileServer(http.Dir(webDir))

	r := chi.NewRouter()

	r.Mount("/", fs)
	r.Get("/api/nextdate", handlers.NextDateHandler)
	r.Post("/api/task", middleware.Auth(handlers.TaskHandler))
	r.Get("/api/tasks", middleware.Auth(handlers.TaskHandler))
	r.Get("/api/task", middleware.Auth(handlers.TaskByIdGet))
	r.Put("/api/task", middleware.Auth(handlers.TaskHandler))
	r.Post("/api/task/done", middleware.Auth(handlers.TaskDonePost))
	r.Delete("/api/task", middleware.Auth(handlers.TaskHandler))
	r.Post("/api/signin", handlers.SingInHandler)

	// Запуск сервера
	log.Printf("Сервер запущен на порту %v", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
