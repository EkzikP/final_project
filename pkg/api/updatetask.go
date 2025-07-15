package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go1f/pkg/db"
	"net/http"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	//десериализуем полученный в запросе JSON
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err = fmt.Errorf("ошибка десериализации JSON")
		writeJson(w, Out{Error: err.Error()}, err)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		err = fmt.Errorf("ошибка десериализации JSON")
		writeJson(w, Out{Error: err.Error()}, err)
		return
	}

	//Проверяем, что поле ID не пустое.
	if task.ID == "" {
		err := fmt.Errorf("не указан ID задачи")
		writeJson(w, Out{Error: err.Error()}, err)
		return
	}

	//Проверяем, что поле Title не пустое.
	if task.Title == "" {
		err := fmt.Errorf("не указан заголовок задачи")
		writeJson(w, Out{Error: err.Error()}, err)
		return
	}

	if err := checkDate(&task); err != nil {
		err = fmt.Errorf("дата представлена в формате, отличном от 20060102")
		writeJson(w, Out{Error: err.Error()}, err)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		err = fmt.Errorf("Задача не найдена")
		writeJson(w, Out{Error: err.Error()}, err)
		return
	}

	writeJson(w, struct{}{}, nil)
}
