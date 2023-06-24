package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	controllers "github.com/akhand3108/restgo/Controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	db := setupDB()

	tokenSecretKey := os.Getenv("TOKEN_SECRET_KEY")
	if tokenSecretKey == "" {
		panic("TOKEN_SECRET_KEY environment variable not set")
	}

	todoController := controllers.NewTodoController(db)
	authController := controllers.NewAuthController(db, tokenSecretKey)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is up and running"))
	})

	r.Post("/signup", authController.Signup)
	r.Post("/signin", authController.Signin)

	r.Group(func(api chi.Router) {
		api.Use(authController.AuthService.Middleware)

		api.Get("/todos", todoController.GetAllTodos)
		api.Post("/todos", todoController.CreateTodo)
		api.Get("/todos/{id}", todoController.GetTodoByID)
		api.Put("/todos/{id}", todoController.UpdateTodoByID)
		api.Delete("/todos/{id}", todoController.DeleteTodoByID)
	})

	fmt.Println("Server is listening on 8080")
	http.ListenAndServe(":8080", r)
}

func setupDB() *sql.DB {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		panic("DB_CONNECTION_STRING environment variable not set")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		passwordhash VARCHAR(255) NOT NULL
	)
`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		done BOOLEAN DEFAULT false,
		user_id INT REFERENCES users(id)
	)
`)
	if err != nil {
		panic(err)
	}

	return db
}
