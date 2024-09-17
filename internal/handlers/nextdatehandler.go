package handlers

import (
	"log"
	"net/http"
	"time"

	dates "github.com/Zelvalna/go_final_project/internal/utils"
	"github.com/Zelvalna/go_final_project/model"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр "now" из запроса и парсим его
	now, err := time.Parse(model.DatePat, r.FormValue("now"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Получаем параметры "date" и "repeat" из запроса
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Вычисляем следующую дату с помощью функции NextDate
	nextDate, err := dates.GetNextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат в ответе
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))

	if err != nil {
		log.Printf("writing tasks data error: %v", err)
	}
}
