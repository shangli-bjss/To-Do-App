package store

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Completed *bool `json:"completed"`
}

type TodoStore struct {
	db *sql.DB
	mu sync.Mutex
}

func NewTodoStore(dbPath string) (*TodoStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &TodoStore{db: db}
	if err := store.initDB(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *TodoStore) initDB() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	)`)
	return err
}

func (s *TodoStore) Close() error {
	return s.db.Close()
}

func (s *TodoStore) GetAllTodos() ([]Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.Id, &t.Title, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (s *TodoStore) CreateTodo(title string, completed *bool) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := Todo{
		Id:        uuid.New().String(),
		Title:     title,
		Completed: completed,
	}

	tx, err := s.db.Begin()
	if err != nil {
		return Todo{}, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO todos (id, title, completed) VALUES (?, ?, ?)",
		todo.Id, todo.Title, todo.Completed)
	if err != nil {
		return Todo{}, err
	}

	if err := tx.Commit(); err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *TodoStore) GetTodo(id string) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var todo Todo
	err := s.db.QueryRow("SELECT id, title, completed FROM todos WHERE id = ?", id).
		Scan(&todo.Id, &todo.Title, &todo.Completed)
	if err == sql.ErrNoRows {
		return Todo{}, errors.New("todo not found")
	} else if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s *TodoStore) UpdateTodo(id string, title string, completed bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("UPDATE todos SET title = ?, completed = ? WHERE id = ?",
		title, completed, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("todo not found")
	}

	return tx.Commit()
}

func (s *TodoStore) DeleteTodo(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("todo not found")
	}

	return tx.Commit()
}