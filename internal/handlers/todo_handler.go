package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

// GetTodos godoc
// @Summary Get all todos
// @Description Get a list of all todo items
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Router /todos [get]
func GetTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, completed, created_at FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item with a title
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.CreateTodoInput true "New Todo"
// @Success 201 {object} models.Todo
// @Router /todos [post]
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("INSERT INTO todos (title) VALUES (?)", input.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	newTodo := models.Todo{
		ID:        int(id),
		Title:     input.Title,
		Completed: false,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

// GetTodoByID godoc
// @Summary Get a todo
// @Description Get a specific todo by ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 404 {string} string "Not Found"
// @Router /todos/{id} [get]
func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	var t models.Todo
	err := database.DB.QueryRow("SELECT id, title, completed, created_at FROM todos WHERE id = ?", id).
		Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// UpdateTodo godoc
// @Summary Update a todo
// @Description Update completion status
// @Tags todos
// @Accept json
// @Param id path int true "Todo ID"
// @Param todo body models.UpdateTodoInput true "Update Todo"
// @Success 200
// @Router /todos/{id} [put]
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	var input models.UpdateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Completed != nil {
		_, err := database.DB.Exec("UPDATE todos SET completed = ? WHERE id = ?", *input.Completed, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Remove a todo item by ID
// @Tags todos
// @Param id path int true "Todo ID"
// @Success 204
// @Router /todos/{id} [delete]
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	_, err := database.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}