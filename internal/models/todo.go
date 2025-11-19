package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Password string `json:"-"` 
}

type AuthInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Comment struct {
	ID        int    `json:"id"`
	TodoID    int    `json:"todo_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Todo struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Completed    bool      `json:"completed"`
	AssignedTo   int       `json:"assigned_to"`
	AssignedName string    `json:"assigned_name"`
	CreatedAt    string    `json:"created_at"`
}

type CreateTodoInput struct {
	Title      string `json:"title"`
	AssignedTo int    `json:"assigned_to"`
}

type CreateCommentInput struct {
	Content string `json:"content"`
}