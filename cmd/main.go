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

func main() {
	database.ConnectDB()
	router := http.NewServeMux()
	router.HandleFunc("GET /users", handlers.GetUsers)
	router.HandleFunc("POST /signup", handlers.Signup)
	router.HandleFunc("POST /login", handlers.Login)
	router.HandleFunc("GET /todos", handlers.GetTodos)
	router.HandleFunc("POST /todos", handlers.CreateTodo) // Admin only
	router.HandleFunc("GET /todos/{id}", handlers.GetTodoByID)
	router.HandleFunc("PUT /todos/{id}", handlers.UpdateTodo)
	router.HandleFunc("DELETE /todos/{id}", handlers.DeleteTodo)
	router.HandleFunc("GET /todos/{id}/comments", handlers.GetComments)
	router.HandleFunc("POST /todos/{id}/comments", handlers.AddComment)
	router.Handle("/swagger/", httpSwagger.WrapHandler)

	fs := http.FileServer(http.Dir("./web"))
	router.Handle("/", fs)

	fmt.Println("🚀 Server starting on http://localhost:8080")
	fmt.Println("📄 Swagger docs at http://localhost:8080/swagger/index.html")
	
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}