package db

import (
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	query := `SELECT * FROM scheduler ORDER BY date LIMIT ?`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		t := &Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`некорректный id для изменения задачи`)
	}
	return nil
}

func GetTask(id string) (*Task, error) {
	query := `SELECT * FROM scheduler WHERE id=?`
	row := DB.QueryRow(query, id)
	t := &Task{}
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return &Task{}, err
	}
	return t, nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id=?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`некорректный id для удаления задачи`)
	}
	return nil
}

func SearchString(search string, limit int) ([]*Task, error) {
	query := `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
	rows, err := DB.Query(query, search, search, limit)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()
	tasks := []*Task{}
	for rows.Next() {
		t := &Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func SearchDate(search string, limit int) ([]*Task, error) {
	query := `SELECT * FROM scheduler WHERE date LIKE ? ORDER BY date LIMIT ?`
	rows, err := DB.Query(query, search, limit)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()
	tasks := []*Task{}
	for rows.Next() {
		t := &Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
