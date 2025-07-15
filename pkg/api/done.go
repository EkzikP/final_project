package api

import (
	"go1f/pkg/db"
	"time"
)

func doneTask(id string) error {
	task, err := db.GetTask(id)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		now := time.Now()
		date, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		task.Date = date
		err = db.UpdateTask(task)
		if err != nil {
			return err
		}
		return nil
	}

	err = db.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil
}
