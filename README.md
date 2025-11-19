#  Go Multi-User Task Manager

This project is a robust, multi-user task management API and front-end built entirely using the **Go Standard Library** (`net/http`) and **MySQL**. It features user authentication (Signup/Login), role-based access control (Admin/User), and commenting on tasks.

##  Features

* **Go Standard Library:** Minimal external dependencies for core routing and HTTP handling.
* **MySQL & Raw SQL:** Direct interaction with the database using `database/sql`.
* **Role-Based Access Control (RBAC):**
    * **Admin:** Can see all tasks, assign tasks, and view all comments.
    * **User:** Can only see tasks assigned to them, complete their tasks, and add comments (but not view previous ones).
* **Secure Authentication:** Passwords are hashed using **BCrypt**.
* **Swagger/OpenAPI:** Automatic documentation generation for the API endpoints.
* **Simple Frontend:** Vanilla HTML, CSS, and JavaScript for testing the workflow.

##  Prerequisites

Ensure you have the following installed on your system:

* **Go:** Version **1.22** or higher (required for the enhanced `net/http` routing).
* **MySQL:** Version 8.0+ or a compatible MariaDB instance.
* **Git:** For cloning and version control.
* **Swag CLI:** The command-line tool for generating Swagger documentation.

```bash
# Install Swag CLI globally 
go install [github.com/swaggo/swag/cmd/swag@latest](https://github.com/swaggo/swag/cmd/swag@latest)
```

## Project Setup and Installation
### Step 1: Clone and Initialize
```
# Clone the repository
git clone [https://github.com/](https://github.com/)<YOUR-USERNAME>/go-todo-api.git
cd go-todo-api

# Download Go dependencies (packages like MySQL driver, bcrypt)
go mod tidy
```

### Step 2: Database Configuration (MySQL)
#### 1- Start your MySQL server.

#### 2- Edit internal/database/db.go and update the Data Source Name (DSN) string with your credentials:
```
// internal/database/db.go
// CHANGE THIS TO YOUR CREDENTIALS
dsn := "root:password@tcp(127.0.0.1:3306)/todo_db?parseTime=true"
```
#### 3- Run the Database Migration Script Execute the following SQL commands to create the tables, define roles, and add password columns:
```
CREATE DATABASE IF NOT EXISTS todo_db;
USE todo_db;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'user') DEFAULT 'user'
);

CREATE TABLE IF NOT EXISTS todos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    assigned_to INT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (assigned_to) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    todo_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```
### Step 3: Generate Swagger Documentation
#### This command parses the annotations in your handler files and generates the necessary docs folder.
```
swag init -g cmd/main.go
```

### Step 4: Run the Application
```
go run cmd/main.go
```
## Usage and Workflow
#### The system enforces a strict workflow defined by user roles.
Resource,Admin Access,User Access,Notes
View Tasks,All tasks,Only assigned tasks,Filtered by X-User-ID header.
Create Task,Yes (Assigns to any user),No (Forbidden),Only admins can assign.
Complete Task,Yes,Yes (Only on assigned task),
Add Comment,Yes,Yes,
View Comments,Yes,No (Forbidden),Requirement enforced.
