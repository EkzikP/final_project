package api

import (
	"go1f/pkg/db"
	"net/http"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		output := "Задача не найдена"
		writeJson(w, Out{Error: output})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		output := "Задача не найдена"
		writeJson(w, Out{Error: output})
		return
	}
	writeJson(w, task)
}
