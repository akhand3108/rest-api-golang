package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {

	todoController := &TodoController{
		todos: make([]Todo, 0),
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is up and running"))
	})

	r.Get("/todos", todoController.getTodosHandler)
	r.Post("/todos", todoController.createTodoHandler)
	r.Get("/todos/{id}", todoController.getTodoHandler)
	r.Put("/todos/{id}", todoController.updateTodoHandler)
	r.Delete("/todos/{id}", todoController.deleteTodoHandler)

	fmt.Println("Server is listening on 8080")
	http.ListenAndServe(":8080", r)
}
