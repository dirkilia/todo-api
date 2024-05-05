package sqlite

import (
	"database/sql"
	"fmt"
	httpserver "todo/internal/http-server"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveTask(task string, is_finished bool) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO tasks(title,is_finished) values(?,?)")
	if err != nil {
		return 0, fmt.Errorf("prepare statement: %v", err)
	}

	res, err := stmt.Exec(task, is_finished)
	if err != nil {
		return 0, fmt.Errorf("execute statement: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

func (s *Storage) GetTasks() ([]httpserver.Task, error) {
	stmt, err := s.db.Prepare("SELECT title, is_finished FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %v", err)
	}

	resultTasks := []httpserver.Task{}
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("rows not found: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		task := httpserver.Task{}
		err = rows.Scan(&task.Title, &task.Is_finished)
		if err != nil {
			return nil, fmt.Errorf("failed to get rows: %v", err)
		}
		resultTasks = append(resultTasks, task)
	}

	return resultTasks, nil
}

func (s *Storage) DeleteTaskById(id int64) (int64, error) {
	stmt, err := s.db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("prepare statement: %v", err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("execute statement: %v", err)
	}
	rows_deleted, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to delete row: %v", err)
	}

	return rows_deleted, nil
}
