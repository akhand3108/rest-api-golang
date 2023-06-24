package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/akhand3108/restgo/models"
)

var (
	ErrTodoNotFound = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

type TodoService struct {
	DB *sql.DB
}

func (ts *TodoService) GetAllTodos(userID int) ([]models.Todo, error) {
	query := `
		SELECT t.id, t.title, t.done
		FROM todos t
		INNER JOIN users u ON t.user_id = u.id
		WHERE t.user_id = $1
		`
	rows, err := ts.DB.Query(query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo = make([]models.Todo, 0)
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Done)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (ts *TodoService) CreateTodo(todo *models.Todo) error {
	query := `
		INSERT INTO todos (title, done, user_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := ts.DB.QueryRow(query, todo.Title, todo.Done, todo.UserID).Scan(&todo.ID)
	return err
}

func (ts *TodoService) GetTodoByID(todoID, userID int) (*models.Todo, error) {
	query := `
		SELECT t.id, t.title, t.done
		FROM todos t
		INNER JOIN users u ON t.user_id = u.id
		WHERE t.id = $1 AND t.user_id = $2
	`

	var todo models.Todo
	row := ts.DB.QueryRow(query, todoID, userID)
	err := row.Scan(&todo.ID, &todo.Title, &todo.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTodoNotFound
		}
		fmt.Printf("Error: %v", err)
		return nil, err
	}
	return &todo, nil
}

func (ts *TodoService) UpdateTodoByID(todo *models.Todo) error {
	query := `
		UPDATE todos
		SET title = $1, done = $2
		WHERE id = $3 AND user_id = $4
	`
	result, err := ts.DB.Exec(query, todo.Title, todo.Done, todo.ID, todo.UserID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}

func (ts *TodoService) DeleteTodoByID(todoID, userID int) error {
	query := `
		DELETE FROM todos
		WHERE id = $1 AND user_id = $2
	`
	result, err := ts.DB.Exec(query, todoID, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}
