package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"todo-api/internal/database"
	"todo-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)



func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserID(r *http.Request) int {
	idStr := r.Header.Get("X-User-ID")
	id, _ := strconv.Atoi(idStr)
	return id
}

func isAdmin(userID int) bool {
	var role string
	err := database.DB.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	return err == nil && role == "admin"
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var input models.AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPwd, _ := hashPassword(input.Password)

	role := "user"
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		role = "admin"
	}

	res, err := database.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", 
		input.Username, hashedPwd, role)
	
	if err != nil {
		http.Error(w, "Username taken or error", http.StatusConflict)
		return
	}

	id, _ := res.LastInsertId()
	json.NewEncoder(w).Encode(map[string]any{"id": id, "role": role, "message": "Signup successful"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User
	var storedHash string

	err := database.DB.QueryRow("SELECT id, username, role, password FROM users WHERE username = ?", input.Username).
		Scan(&user.ID, &user.Username, &user.Role, &storedHash)

	if err != nil || !checkPasswordHash(input.Password, storedHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// In a real app, send a JWT token here. 
	// For this demo, we send the ID back to be used in headers.
	json.NewEncoder(w).Encode(user)
}

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

func GetTodos(w http.ResponseWriter, r *http.Request) {
	currentUserID := getUserID(r)
	
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
		var assignedTo sql.NullInt64
		var assignedName sql.NullString
		rows.Scan(&t.ID, &t.Title, &t.Completed, &assignedTo, &assignedName, &t.CreatedAt)
		t.AssignedTo = int(assignedTo.Int64)
		t.AssignedName = assignedName.String
		todos = append(todos, t)
	}
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	if !isAdmin(getUserID(r)) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}
	var input models.CreateTodoInput
	json.NewDecoder(r.Body).Decode(&input)
	database.DB.Exec("INSERT INTO todos (title, assigned_to) VALUES (?, ?)", input.Title, input.AssignedTo)
	w.WriteHeader(http.StatusCreated)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	var input struct { Completed bool `json:"completed"` }
	json.NewDecoder(r.Body).Decode(&input)
	database.DB.Exec("UPDATE todos SET completed = ? WHERE id = ?", input.Completed, idStr)
	w.WriteHeader(http.StatusOK)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)
	database.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	w.WriteHeader(http.StatusNoContent)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}


func GetComments(w http.ResponseWriter, r *http.Request) {
	currentUserID := getUserID(r)

	if !isAdmin(currentUserID) {
		http.Error(w, "Only Admins can view comments", http.StatusForbidden)
		return
	}

	todoID := r.PathValue("id")
	rows, _ := database.DB.Query(`
		SELECT c.id, c.todo_id, c.user_id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.todo_id = ? ORDER BY c.created_at ASC`, todoID)
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
	userID := getUserID(r)
	var input models.CreateCommentInput
	json.NewDecoder(r.Body).Decode(&input)
	database.DB.Exec("INSERT INTO comments (todo_id, user_id, content) VALUES (?, ?, ?)", todoID, userID, input.Content)
	w.WriteHeader(http.StatusCreated)
}