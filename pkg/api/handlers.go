package api

import (
	"go1f/pkg/db"
	"net/http"
	"time"
)

// Возвращаем следующую дату задачи или ошибку
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	//Получаем параметры из запроса
	var now time.Time
	nowstr := r.FormValue("now")
	if nowstr == "" {
		now = time.Now().UTC()
	} else {
		now, _ = time.Parse(layout, nowstr)
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	//Получаем следующую дату задачи или ошибку
	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		//Возвращаем ошибку
		w.Write([]byte(err.Error()))
	}
	//Возвращаем следующую дату
	w.Write([]byte(nextDate))
}

// Работа с задачей
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		output := "Не указан id задачи"
		writeJson(w, Out{Error: output})
		return
	}

	err := doneTask(id)
	if err != nil {
		output := "Ошибка изменения статуса задачи"
		writeJson(w, Out{Error: output})
		return
	}

	writeJson(w, struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		output := "Не указан id задачи"
		writeJson(w, Out{Error: output})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		output := "Ошибка удаления задачи"
		writeJson(w, Out{Error: output})
		return
	}
	writeJson(w, struct{}{})
}
