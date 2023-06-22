package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TodoController struct {
	// todoDB *sql.DB
	todos []Todo
}

func (tdc *TodoController) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tdc.todos)
}

func (tdc *TodoController) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tdc.todos = append(tdc.todos, todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (tdc *TodoController) getTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	for _, todo := range tdc.todos {
		if todo.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	http.NotFound(w, r)

}

func (tdc *TodoController) updateTodoHandler(w http.ResponseWriter, r *http.Request) {

	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	updatedTodo.ID = id
	for i, todo := range tdc.todos {
		if todo.ID == id {
			tdc.todos[i] = updatedTodo
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.NotFound(w, r)

}

func (tdc *TodoController) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	for i, todo := range tdc.todos {
		if todo.ID == id {
			tdc.todos = append(tdc.todos[:i], tdc.todos[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.NotFound(w, r)
}
