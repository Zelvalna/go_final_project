package handlers

import (
	"fmt"
	"github.com/Zelvalna/go_final_project/constans"
	"github.com/Zelvalna/go_final_project/internal/utils"
	"net/http"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр "now" из запроса и парсим его
	now, err := time.Parse(constans.DatePat, r.FormValue("now"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Получаем параметры "date" и "repeat" из запроса
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Вычисляем следующую дату с помощью функции NextDate
	nextDate, err := utils.GetNextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат в ответе
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))

	if err != nil {
		http.Error(w, fmt.Errorf("writing tasks data error: %w", err).Error(), http.StatusInternalServerError)
	}
}
