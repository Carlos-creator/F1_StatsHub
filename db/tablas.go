package db

import (
    "database/sql"
    "log"
)

// CreateTables creates the necessary tables in the SQLite database.
func CreateTables(db *sql.DB) {
    queries := []string{
        // Table drivers
        `CREATE TABLE IF NOT EXISTS drivers (
            driver_number INTEGER PRIMARY KEY,
            first_name TEXT,
            last_name TEXT,
            name_acronym TEXT,
            team_name TEXT,
            country_code TEXT
        );`,
        // Table sessions
        `CREATE TABLE IF NOT EXISTS sessions (
            session_key INTEGER PRIMARY KEY,
            session_name TEXT,
            session_type TEXT,
            location TEXT,
            country_name TEXT,
            year INTEGER,
            circuit_short_name TEXT,
            date_start TEXT
        );`,
        // Table results
        `CREATE TABLE IF NOT EXISTS results (
            session_key INTEGER,
            driver_number INTEGER,
            final_position INTEGER,
            laps INTEGER,
            top_speed REAL,
            PRIMARY KEY(session_key, driver_number)
        );`,
        // Table positions
        `CREATE TABLE IF NOT EXISTS positions (
            driver_number INTEGER,
            session_key INTEGER,
            position INTEGER,
            date TEXT,
            PRIMARY KEY(driver_number, session_key)
        );`,
        // Table laps
        `CREATE TABLE IF NOT EXISTS laps (
            driver_number INTEGER,
            session_key INTEGER,
            lap_number INTEGER,
            lap_duration REAL,
            duration_sector_1 REAL,
            duration_sector_2 REAL,
            duration_sector_3 REAL,
            st_speed REAL,
            date_start TEXT,
            PRIMARY KEY(driver_number, session_key, lap_number)
        );`,
    }
    for _, q := range queries {
        _, err := db.Exec(q)
        if err != nil {
            log.Fatalf("❌ Error creando tabla: %v", err)
        }
    }
    log.Println("✅ Tablas creadas correctamente.")
}
