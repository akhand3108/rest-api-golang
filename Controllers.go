package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TodoController struct {
	DB *sql.DB
}

func (tdc *TodoController) GetAllTodos(w http.ResponseWriter, r *http.Request) {

	rows, err := tdc.DB.Query("SELECT id,title,done FROM todos;")

	if err != nil {
		http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
		fmt.Printf("Error: %v", err)
		return
	}
	defer rows.Close()

	var todos []Todo = make([]Todo, 0)
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Done)
		if err != nil {
			http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
			fmt.Printf("Error: %v", err)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)

}

func (tdc *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Printf("Error: %v", err)
		return
	}

	err = tdc.DB.QueryRow("INSERT INTO todos (title, done) VALUES ($1, $2) RETURNING id;", todo.Title, todo.Done).Scan(&todo.ID)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		fmt.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (tdc *TodoController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var todo Todo
	row := tdc.DB.QueryRow("SELECT * FROM todos WHERE id=$1", id)
	err = row.Scan(&todo.ID, &todo.Title, &todo.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve todo", http.StatusInternalServerError)
			fmt.Printf("Error: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)

}

func (tdc *TodoController) UpdateTodoByID(w http.ResponseWriter, r *http.Request) {

	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var todo Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var updateQuery string
	params := make([]interface{}, 0)

	if todo.Title != "" {
		updateQuery = "UPDATE todos SET title=$1"
		params = append(params, todo.Title)

		if todo.Done {
			updateQuery += ", done=$2"
			params = append(params, todo.Done)
		}
	} else if todo.Done {
		updateQuery = "UPDATE todos SET done=$1"
		params = append(params, todo.Done)
	} else {
		http.Error(w, "Invalid update parameters", http.StatusBadRequest)
		return
	}
	params = append(params, id)
	updateQuery = updateQuery + " WHERE id=$" + strconv.Itoa(len(params))

	result, err := tdc.DB.Exec(updateQuery, params...)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		fmt.Printf("Error: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	row := tdc.DB.QueryRow("SELECT id, title, done FROM todos WHERE id=$1", id)
	err = row.Scan(&todo.ID, &todo.Title, &todo.Done)
	if err != nil {
		http.Error(w, "Failed to retrieve updated todo", http.StatusInternalServerError)
		fmt.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (tdc *TodoController) DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		fmt.Printf("Error: %v", err)
		return
	}

	result, err := tdc.DB.Exec("DELETE FROM todos where id=$1", id)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		fmt.Printf("Error: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	http.NotFound(w, r)
}
