package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zelvalna/go_final_project/constans"
	"github.com/Zelvalna/go_final_project/internal/storage"
	"github.com/Zelvalna/go_final_project/internal/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

// setErrorResponse Функция для создания и отправки ответа об ошибке
func setErrorResponse(w http.ResponseWriter, s string, err error) {
	errorResponse := constans.ErrorResponse{
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
	var taskData constans.Task

	var buffer bytes.Buffer
	// Чтение тела запроса
	if _, err := buffer.ReadFrom(r.Body); err != nil {
		setErrorResponse(w, "body getting error", err)
		return
	}
	// Десериализация JSON в структуру Task
	if err := json.Unmarshal(buffer.Bytes(), &taskData); err != nil {
		setErrorResponse(w, "JSON deserialization error", err)
		return
	}
	// Установка даты по умолчанию или проверка формата даты
	if len(taskData.Date) == 0 {
		taskData.Date = time.Now().Format(constans.DatePat)
	} else {
		date, err := time.Parse(constans.DatePat, taskData.Date)
		if err != nil {
			setErrorResponse(w, "bad data format", err)
			return
		}

		if date.Before(time.Now()) {
			taskData.Date = time.Now().Format(constans.DatePat)
		}
	}
	// Проверка заголовка задачи
	if len(taskData.Title) == 0 {
		setErrorResponse(w, "invalid title", errors.New("title is empty"))
		return
	}
	// Проверка формата повтора
	if len(taskData.Repeat) > 0 {
		if _, err := utils.GetNextDate(time.Now(), taskData.Date, taskData.Repeat); err != nil {
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
	taskIdData, err := json.Marshal(constans.TaskIdResponse{Id: taskId})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(taskIdData)
	if err != nil {
		setErrorResponse(w, "writing task id error", err)
		return
	}

	log.Println(fmt.Sprintf("Added task with id=%d", taskId))
}

func TasksReadGet(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []constans.Task

	if len(search) > 0 {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			tasks, err = storage.SearchTasks(search)
		} else {
			tasks, err = storage.SearchTasksByDate(date.Format(constans.DatePat))
		}
	} else {
		err := errors.New("")
		tasks, err = storage.ReadTasks()
		if err != nil {
			setErrorResponse(w, "failed to get tasks", err)
			return
		}
	}

	tasksData, err := json.Marshal(constans.Tasks{Tasks: tasks})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	if err != nil {
		setErrorResponse(w, "writing tasks error", err)
		return
	}

	log.Println(fmt.Sprintf("Read %d tasks", len(tasks)))
}

func TaskByIdGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := storage.ReadTaskById(id)
	if err != nil {
		setErrorResponse(w, "failed to get task by id", err)
		return
	}
	taskData, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskData)
	if err != nil {
		setErrorResponse(w, "writing task error", err)
		return
	}

	log.Println(fmt.Sprintf("Read task with id=%s", id))
}
func TaskUpdatePut(w http.ResponseWriter, r *http.Request) {
	var task constans.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		setErrorResponse(w, "body getting error", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &task); err != nil {
		setErrorResponse(w, "JSON deserialization error", err)
	}
	if len(task.ID) == 0 {
		setErrorResponse(w, "invalid id", errors.New("id is empty"))
		return
	}
	if _, err := strconv.Atoi(task.ID); err != nil {
		setErrorResponse(w, "invalid id", err)
		return
	}
	if _, err := time.Parse(constans.DatePat, task.Date); err != nil {
		setErrorResponse(w, "invalid date format", err)
		return
	}
	if len(task.Title) == 0 {
		setErrorResponse(w, "invalid title", errors.New("title is empty"))
		return
	}
	if len(task.Repeat) > 0 {
		if _, err := utils.GetNextDate(time.Now(), task.Date, task.Repeat); err != nil {
			setErrorResponse(w, "invalid repeat format", errors.New("no such format"))
			return
		}
	}

	_, err := storage.UpdateTask(task)
	if err != nil {
		setErrorResponse(w, "failed to update task", errors.New("failed to update task"))
		return
	}
	taskData, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskData)
	if err != nil {
		setErrorResponse(w, "updating task error", err)
		return
	}

	log.Println(fmt.Sprintf("Updated task with id=%s", task.ID))

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
		task.Date, err = utils.GetNextDate(time.Now(), task.Date, task.Repeat)
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

	// Возвращаем обновлённую задачу
	taskData, err := json.Marshal(struct{}{})
	if err != nil {
		setErrorResponse(w, "JSON serialization error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskData)
	if err != nil {
		setErrorResponse(w, "writing task error", err)
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

	taskData, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskData)
	if err != nil {
		setErrorResponse(w, "writing task error", err)
	}
	log.Println(fmt.Sprintf("Deleted task with id=%s", id))
}
