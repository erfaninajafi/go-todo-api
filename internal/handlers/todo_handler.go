package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

// Helper to get current user from headers (Simulated Auth)
func getUserID(r *http.Request) int {
	idStr := r.Header.Get("X-User-ID")
	id, _ := strconv.Atoi(idStr)
	return id
}

// Helper to check if user is admin
func isAdmin(userID int) bool {
	var role string
	err := database.DB.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	return err == nil && role == "admin"
}

// --- USERS ---

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, _ := database.DB.Query("SELECT id, username, role FROM users")
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Role)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

// --- TODOS ---

func GetTodos(w http.ResponseWriter, r *http.Request) {
	currentUserID := getUserID(r)

	// Logic: Admin sees ALL. Users see only their ASSIGNED tasks.
	query := `
		SELECT t.id, t.title, t.completed, t.assigned_to, u.username, t.created_at 
		FROM todos t 
		LEFT JOIN users u ON t.assigned_to = u.id`
	
	var rows *sql.Rows
	var err error

	if isAdmin(currentUserID) {
		rows, err = database.DB.Query(query)
	} else {
		query += " WHERE t.assigned_to = ?"
		rows, err = database.DB.Query(query, currentUserID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		var assignedTo sql.NullInt64 // Handle NULLs safely
		var assignedName sql.NullString

		rows.Scan(&t.ID, &t.Title, &t.Completed, &assignedTo, &assignedName, &t.CreatedAt)
		
		t.AssignedTo = int(assignedTo.Int64)
		t.AssignedName = assignedName.String
		todos = append(todos, t)
	}
	
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	currentUserID := getUserID(r)
	
	// Only Admin can create/assign tasks
	if !isAdmin(currentUserID) {
		http.Error(w, "Only Admins can assign tasks", http.StatusForbidden)
		return
	}

	var input models.CreateTodoInput
	json.NewDecoder(r.Body).Decode(&input)

	res, err := database.DB.Exec("INSERT INTO todos (title, assigned_to) VALUES (?, ?)", input.Title, input.AssignedTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]any{"id": id, "status": "assigned"})
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	var t models.Todo
	var assignedTo sql.NullInt64
	var assignedName sql.NullString

	// Updated query to include assignment info
	query := `SELECT t.id, t.title, t.completed, t.assigned_to, u.username, t.created_at 
	          FROM todos t 
	          LEFT JOIN users u ON t.assigned_to = u.id 
	          WHERE t.id = ?`

	err := database.DB.QueryRow(query, id).
		Scan(&t.ID, &t.Title, &t.Completed, &assignedTo, &assignedName, &t.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.AssignedTo = int(assignedTo.Int64)
	t.AssignedName = assignedName.String

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	
	var input struct { Completed bool `json:"completed"` }
	json.NewDecoder(r.Body).Decode(&input)

	_, err := database.DB.Exec("UPDATE todos SET completed = ? WHERE id = ?", input.Completed, idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

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

// --- COMMENTS ---

func GetComments(w http.ResponseWriter, r *http.Request) {
	todoID := r.PathValue("id")
	
	query := `
		SELECT c.id, c.todo_id, c.user_id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.todo_id = ? ORDER BY c.created_at ASC
	`
	rows, err := database.DB.Query(query, todoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		rows.Scan(&c.ID, &c.TodoID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt)
		comments = append(comments, c)
	}
	json.NewEncoder(w).Encode(comments)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	todoID := r.PathValue("id")
	userID := getUserID(r) // Who is commenting?

	var input models.CreateCommentInput
	json.NewDecoder(r.Body).Decode(&input)

	_, err := database.DB.Exec("INSERT INTO comments (todo_id, user_id, content) VALUES (?, ?, ?)", 
		todoID, userID, input.Content)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}