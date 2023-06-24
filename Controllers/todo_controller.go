package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akhand3108/restgo/models"
	"github.com/akhand3108/restgo/services"
	"github.com/go-chi/chi/v5"
)

type TodoController struct {
	TodoService *services.TodoService
}

func NewTodoController(db *sql.DB) *TodoController {
	todoService := &services.TodoService{
		DB: db,
	}

	return &TodoController{
		TodoService: todoService,
	}
}

func (tc *TodoController) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	todos, err := tc.TodoService.GetAllTodos(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve todos", http.StatusInternalServerError)
		fmt.Printf("Error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)

}

func (tc *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Printf("Error: %v", err)
		return
	}

	todo.UserID = userID
	err = tc.TodoService.CreateTodo(&todo)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		fmt.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (tc *TodoController) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := tc.TodoService.GetTodoByID(id, userID)
	if err != nil {
		if err == services.ErrTodoNotFound {
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

func (tc *TodoController) UpdateTodoByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	todo.ID = id
	todo.UserID = userID

	err = tc.TodoService.UpdateTodoByID(&todo)
	if err != nil {
		if err == services.ErrTodoNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update todo", http.StatusInternalServerError)
			fmt.Printf("Error: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (tc *TodoController) DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	todoID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		fmt.Printf("Error: %v", err)
		return
	}

	err = tc.TodoService.DeleteTodoByID(id, userID)
	if err != nil {
		if err == services.ErrTodoNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
			fmt.Printf("Error: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
