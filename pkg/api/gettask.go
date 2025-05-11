package api

import (
	"errors"
	"fmt"
	"go1f/pkg/db"
	"net/http"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		err := errors.New("Задача не найдена")
		writeJson(w, Out{Error: err.Error()}, err)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		err = fmt.Errorf("Задача не найдена")
		writeJson(w, Out{Error: err.Error()}, err)
		return
	}
	writeJson(w, task, nil)
}
