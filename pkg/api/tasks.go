package api

import (
	"go1f/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	if search != "" {
		searchDate, err := time.Parse("02.01.2006", search)
		if err != nil {
			tasks, err := db.SearchString("%"+search+"%", 50)
			if err != nil {
				output := "ошибка запроса к базе данных"
				writeJson(w, Out{Error: output})
				return
			}
			writeJson(w, TasksResp{Tasks: tasks})
			return
		}
		tasks, err := db.SearchDate(searchDate.Format(layout), 50)
		if err != nil {
			output := "ошибка запроса к базе данных"
			writeJson(w, Out{Error: output})
			return
		}
		writeJson(w, TasksResp{Tasks: tasks})
		return
	}
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		output := "ошибка запроса к базе данных"
		writeJson(w, Out{Error: output})
		return
	}
	writeJson(w, TasksResp{Tasks: tasks})
}
