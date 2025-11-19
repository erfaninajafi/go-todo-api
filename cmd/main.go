package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/internal/database"
	"todo-api/internal/handlers"

	_ "todo-api/docs" // Import generated docs

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Go Raw SQL Todo API
// @version 1.0
// @description This is a sample server using Standard Lib & MySQL.
// @host localhost:8080
// @BasePath /
func main() {
	// 1. Initialize Database
	database.ConnectDB()

	// 2. Setup Routing (Go 1.22+)
	router := http.NewServeMux()

	// --- API Routes ---
	
	// Users (New)
	router.HandleFunc("GET /users", handlers.GetUsers)

	// Todos
	router.HandleFunc("GET /todos", handlers.GetTodos)
	router.HandleFunc("POST /todos", handlers.CreateTodo) // Admin only
	router.HandleFunc("GET /todos/{id}", handlers.GetTodoByID)
	router.HandleFunc("PUT /todos/{id}", handlers.UpdateTodo)
	router.HandleFunc("DELETE /todos/{id}", handlers.DeleteTodo)

	// Comments (New)
	router.HandleFunc("GET /todos/{id}/comments", handlers.GetComments)
	router.HandleFunc("POST /todos/{id}/comments", handlers.AddComment)

	// Swagger Docs Route
	router.Handle("/swagger/", httpSwagger.WrapHandler)

	// --- Frontend (Static Files) ---
	// Serve files from "web" folder
	fs := http.FileServer(http.Dir("./web"))
	router.Handle("/", fs)

	// 3. Start Server
	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("📄 Swagger docs at http://localhost:8080/swagger/index.html")
	
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}