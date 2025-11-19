package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql" // Import driver anonymously
)

var DB *sql.DB

func ConnectDB() {
    // format: username:password@tcp(host:port)/dbname
    dsn := "root:Erfan.judo8@tcp(127.0.0.1:3306)/todo_db?parseTime=true"
    
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error opening database: ", err)
    }

    // Ping to verify connection
    if err = DB.Ping(); err != nil {
        log.Fatal("Error connecting to database: ", err)
    }

    // Connection pool settings
    DB.SetMaxOpenConns(10)
    DB.SetMaxIdleConns(5)
    DB.SetConnMaxLifetime(5 * time.Minute)

    fmt.Println("✅ Connected to MySQL database!")
}