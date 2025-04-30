package main

import (
    "database/sql"
    "log"

    "f1statshub/db"
    "f1statshub/endpoints"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Open (or create) the SQLite database file
    conn, err := sql.Open("sqlite3", "proxy.db")
    if err != nil {
        log.Fatal("Error opening database:", err)
    }
    defer conn.Close()

    // Create necessary tables and load initial data
    db.CreateTables(conn)
    db.CargarDatosDesdeOpenF1(conn)

    // Set up the Gin router with all endpoints
    router := endpoints.SetupRouter(conn)
    log.Println("🚀 Server running on http://localhost:8080")
    // Run the server on port 8080
    router.Run(":8080")
}
