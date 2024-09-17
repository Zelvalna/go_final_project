package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Zelvalna/go_final_project/internal/storage"
	dates "github.com/Zelvalna/go_final_project/internal/utils"
	"github.com/Zelvalna/go_final_project/model"
)

// setErrorResponse Функция для создания и отправки ответа об ошибке
func setErrorResponse(w http.ResponseWriter, s string, err error) {
	errorResponse := model.ErrorResponse{
		Error: fmt.Errorf("%s: %w", s, err).Error()}

	// Сериализация ответа об ошибке
	errorData, _ := json.Marshal(errorResponse)
	w.WriteHeader(http.StatusBadRequest)
	// Пишем ответ
	_, writeErr := w.Write(errorData)
	if writeErr != nil {
		http.Error(w, fmt.Errorf("error: %w", writeErr).Error(), http.StatusBadRequest)
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		TaskAddPost(w, r)
	case http.MethodGet:
		TasksReadGet(w, r)
	case http.MethodPut:
		TaskUpdatePut(w, r)
	case http.MethodDelete:
		TaskDelete(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
func TaskAddPost(w http.ResponseWriter, r *http.Request) {
	var taskData model.Task

	// Декодирование JSON тела запроса
	if err := json.NewDecoder(r.Body).Decode(&taskData); err != nil {
		setErrorResponse(w, "JSON deserialization error", err)
		return
	}
	// Установка даты по умолчанию или проверка формата даты
	if len(taskData.Date) == 0 {
		taskData.Date = time.Now().Format(model.DatePat)
	} else {
		date, err := time.Parse(model.DatePat, taskData.Date)
		if err != nil {
			setErrorResponse(w, "bad data format", err)
			return
		}

		if date.Before(time.Now()) {
			taskData.Date = time.Now().Format(model.DatePat)
		}
	}
	// Проверка заголовка задачи
	if len(taskData.Title) == 0 {
		setErrorResponse(w, "invalid title", errors.New("title is empty"))
		return
	}
	// Проверка формата повтора
	if len(taskData.Repeat) > 0 {
		if _, err := dates.GetNextDate(time.Now(), taskData.Date, taskData.Repeat); err != nil {
			setErrorResponse(w, "invalid repeat format", errors.New("no such format"))
			return
		}
	}
	// Добавление задачи в базу данных
	taskId, err := storage.InsertTask(taskData)
	if err != nil {
		setErrorResponse(w, "failed to create task", err)
		return
	}
	// Возвращение ID созданной задачи
	jsonResponse(w, http.StatusCreated)
	if err := json.NewEncoder(w).Encode(model.TaskIdResponse{Id: taskId}); err != nil {
		setErrorResponse(w, "failed to encode response", err)
		return
	}
	log.Println(fmt.Sprintf("Added task with id=%d", taskId))
}

func jsonResponse(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
}

func TasksReadGet(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := fetchTasks(search)
	if err != nil {
		setErrorResponse(w, "failed to get tasks", err)
		return
	}

	jsonResponse(w, http.StatusOK)
	if err := json.NewEncoder(w).Encode(model.Tasks{Tasks: tasks}); err != nil {
		setErrorResponse(w, "failed to encode response", err)
		return
	}

	log.Println(fmt.Sprintf("Read %d tasks", len(tasks)))
}

func fetchTasks(search string) ([]model.Task, error) {
	if len(search) > 0 {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			return storage.SearchTasksByDate(date.Format(model.DatePat))
		}
		return storage.SearchTasks(search)
	}
	return storage.ReadTasks()
}

func TaskByIdGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := storage.ReadTaskById(id)
	if err != nil {
		setErrorResponse(w, "failed to get task by id", err)
		return
	}
	jsonResponse(w, http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		setErrorResponse(w, "failed to encode response", err)
		return
	}

	log.Println(fmt.Sprintf("Read task with id=%s", id))
}
func TaskUpdatePut(w http.ResponseWriter, r *http.Request) {
	var task model.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		setErrorResponse(w, "JSON deserialization error", err)
		return
	}
	if len(task.ID) == 0 {
		setErrorResponse(w, "invalid id", errors.New("id is empty"))
		return
	}
	if _, err := strconv.Atoi(task.ID); err != nil {
		setErrorResponse(w, "invalid id", err)
		return
	}
	if _, err := time.Parse(model.DatePat, task.Date); err != nil {
		setErrorResponse(w, "invalid date format", err)
		return
	}
	if len(task.Title) == 0 {
		setErrorResponse(w, "invalid title", errors.New("title is empty"))
		return
	}
	if len(task.Repeat) > 0 {
		if _, err := dates.GetNextDate(time.Now(), task.Date, task.Repeat); err != nil {
			setErrorResponse(w, "invalid repeat format", errors.New("no such format"))
			return
		}
	}

	_, err := storage.UpdateTask(task)
	if err != nil {
		setErrorResponse(w, "failed to update task", errors.New("failed to update task"))
		return
	}

	jsonResponse(w, http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		setErrorResponse(w, "failed to encode response", err)
		return
	}

}
func TaskDonePost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := storage.ReadTaskById(id)
	if err != nil {
		setErrorResponse(w, "failed to get task by id", err)
		return
	}

	if task.Repeat == "" {
		err = storage.DeleteTask(task.ID)
		if err != nil {
			setErrorResponse(w, "failed to delete task", err)
			return
		}
		log.Println(fmt.Sprintf("task with id=%s was deleted", task.ID))
	} else {
		task.Date, err = dates.GetNextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			setErrorResponse(w, "failed to get next date", err)
			return
		}
		// Обновляем задачу с новой датой
		_, err = storage.UpdateTask(task)
		if err != nil {
			setErrorResponse(w, "failed to update task", err)
			return
		}
	}

	jsonResponse(w, http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		setErrorResponse(w, "failed to encode response", err)
		return
	}

	log.Println(fmt.Sprintf("Updated task with id=%s", task.ID))

}
func TaskDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := storage.DeleteTask(id)
	if err != nil {
		setErrorResponse(w, "failed to delete task", err)
		return
	}

	jsonResponse(w, http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		setErrorResponse(w, "failed to encode response", err)
		return
	}
	log.Println(fmt.Sprintf("Deleted task with id=%s", id))
}
