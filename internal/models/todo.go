package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Comment struct {
	ID        int    `json:"id"`
	TodoID    int    `json:"todo_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"` // Enriched field for display
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Todo struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Completed    bool      `json:"completed"`
	AssignedTo   int       `json:"assigned_to"`
	AssignedName string    `json:"assigned_name"` // Enriched field
	CreatedAt    string    `json:"created_at"`
	Comments     []Comment `json:"comments,omitempty"` // Nested comments
}

type CreateTodoInput struct {
	Title      string `json:"title"`
	AssignedTo int    `json:"assigned_to"` // Admin assigns this
}

type CreateCommentInput struct {
	Content string `json:"content"`
}