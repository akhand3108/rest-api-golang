package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	connectionString := "postgres://baloo:junglebook@localhost:5432/lenslocked?sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	todoController := &TodoController{
		DB: db,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is up and running"))
	})

	r.Get("/todos", todoController.GetAllTodos)
	r.Post("/todos", todoController.CreateTodo)
	r.Get("/todos/{id}", todoController.GetTodoByID)
	r.Put("/todos/{id}", todoController.UpdateTodoByID)
	r.Delete("/todos/{id}", todoController.DeleteTodoByID)

	fmt.Println("Server is listening on 8080")
	http.ListenAndServe(":8080", r)
}
