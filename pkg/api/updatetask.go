package api

import (
	"bytes"
	"encoding/json"
	"go1f/pkg/db"
	"net/http"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	//Проверяем, что поле ID не пустое.
	if task.ID == "" {
		output := "не указан ID задачи"
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

	err = db.UpdateTask(&task)
	if err != nil {
		output := "Задача не найдена"
		writeJson(w, Out{Error: output})
		return
	}

	writeJson(w, struct{}{})
}
