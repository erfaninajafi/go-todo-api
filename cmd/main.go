package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/internal/database"
	"todo-api/internal/handlers"

	_ "todo-api/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Go Raw SQL Todo API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	database.ConnectDB()

	router := http.NewServeMux()

	// --- API Routes ---
	router.HandleFunc("GET /todos", handlers.GetTodos)
	router.HandleFunc("POST /todos", handlers.CreateTodo)
	router.HandleFunc("GET /todos/{id}", handlers.GetTodoByID)
	router.HandleFunc("PUT /todos/{id}", handlers.UpdateTodo) // New Handler!
	router.HandleFunc("DELETE /todos/{id}", handlers.DeleteTodo)

	// --- Swagger ---
	router.Handle("/swagger/", httpSwagger.WrapHandler)

	// --- Frontend (Static Files) ---
	// This tells Go to serve files from the "web" folder
	fs := http.FileServer(http.Dir("./web"))
    
    // We mount the file server at "/"
    // Note: Specific routes (like /todos) take precedence over / in Go 1.22
	router.Handle("/", fs)

	fmt.Println("🚀 Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}