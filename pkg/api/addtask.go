package api

import (
	"bytes"
	"encoding/json"
	"go1f/pkg/db"
	"net/http"
	"strconv"
	"time"
)

type Out struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	//десериализуем полученный в запросе JSON
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		output := "ошибка десериализации JSON"
		writeJson(w, Out{Error: output})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		output := "ошибка десериализации JSON"
		writeJson(w, Out{Error: output})
		return
	}

	//Проверяем, что поле Title не пустое.
	if task.Title == "" {
		output := "не указан заголовок задачи"
		writeJson(w, Out{Error: output})
		return
	}

	if err := checkDate(&task); err != nil {
		output := "дата представлена в формате, отличном от 20060102"
		writeJson(w, Out{Error: output})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		output := err.Error()
		writeJson(w, Out{Error: output})
		return
	}

	writeJson(w, Out{ID: strconv.FormatInt(id, 10)})
}

// Возвращаем ответный JSON клиенту

func writeJson(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//// так как все успешно, то статус OK
	//w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

func checkDate(task *db.Task) error {
	var next string

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// если сегодня (now) больше task.Date (t)
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format("20060102")
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}
	return nil
}
