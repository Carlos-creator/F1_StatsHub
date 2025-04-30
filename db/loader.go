package db

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"


    "f1statshub/models"
)

// CargarDatosDesdeOpenF1 loads initial data from the OpenF1 API into the database.
func CargarDatosDesdeOpenF1(db *sql.DB) {
    log.Println("📡 Cargando datos desde OpenF1...")

    cargarPilotos(db)
    cargarSesiones(db)
    cargarPosiciones(db)
    cargarVueltas(db)
    cargarResultados(db)

    log.Println("✅ Datos iniciales cargados correctamente.")
}

func cargarPilotos(db *sql.DB) {
    type DriverAPI struct {
        DriverNumber int    `json:"driver_number"`
        FirstName    string `json:"first_name"`
        LastName     string `json:"last_name"`
        NameAcronym  string `json:"name_acronym"`
        TeamName     string `json:"team_name"`
        CountryCode  string `json:"country_code"`
    }
    insertStmt := `INSERT OR IGNORE INTO drivers 
        (driver_number, first_name, last_name, name_acronym, team_name, country_code) 
        VALUES (?, ?, ?, ?, ?, ?)`

    // Known session keys to retrieve drivers (covering all drivers)
    sessions := map[int][]int{
        9574: {1, 2, 3, 4, 10, 11, 14, 16, 18, 20, 22, 23, 24, 27, 31, 44, 55, 63, 77, 81},
        9636: {30, 50, 43},
    }
    for sessionKey, allowedDrivers := range sessions {
        url := fmt.Sprintf("https://api.openf1.org/v1/drivers?session_key=%d", sessionKey)
        resp, err := http.Get(url)
        if err != nil {
            log.Println("❌ Error al obtener pilotos:", err)
            continue
        }
        defer resp.Body.Close()
        var driversAPI []DriverAPI
        if err := json.NewDecoder(resp.Body).Decode(&driversAPI); err != nil {
            log.Println("❌ Error decodificando pilotos:", err)
            continue
        }
        for _, d := range driversAPI {
            if contains(allowedDrivers, d.DriverNumber) {
                _, err := db.Exec(insertStmt, d.DriverNumber, d.FirstName, d.LastName, d.NameAcronym, d.TeamName, d.CountryCode)
                if err != nil {
                    log.Println("❌ Error insertando piloto:", err)
                }
            }
        }
    }
    log.Println("✅ Pilotos cargados correctamente.")
}

func cargarSesiones(db *sql.DB) {
    url := "https://api.openf1.org/v1/sessions?session_name=Race&year=2024"
    resp, err := http.Get(url)
    if err != nil {
        log.Println("❌ Error al obtener sesiones:", err)
        return
    }
    defer resp.Body.Close()
    type SessionAPI struct {
        SessionKey       int    `json:"session_key"`
        SessionName      string `json:"session_name"`
        SessionType      string `json:"session_type"`
        Location         string `json:"location"`
        CountryName      string `json:"country_name"`
        Year             int    `json:"year"`
        CircuitShortName string `json:"circuit_short_name"`
        DateStart        string `json:"date_start"`
    }
    var sessions []SessionAPI
    if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
        log.Println("❌ Error decodificando sesiones:", err)
        return
    }
    insertStmt := `INSERT OR IGNORE INTO sessions 
        (session_key, session_name, session_type, location, country_name, year, circuit_short_name, date_start) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
    for _, s := range sessions {
        _, err := db.Exec(insertStmt, s.SessionKey, s.SessionName, s.SessionType, s.Location, s.CountryName, s.Year, s.CircuitShortName, s.DateStart)
        if err != nil {
            log.Println("❌ Error insertando sesión:", err)
        }
    }
    log.Printf("🔍 %d sesiones de carreras cargadas.\n", len(sessions))
}

func cargarPosiciones(db *sql.DB) {
    log.Println("📡 Cargando posiciones (resultados finales) de carreras...")
    // Clear any existing data
    _, _ = db.Exec("DELETE FROM positions")
    // We will use a subset of sessions (3 known race sessions) for initial data
    sessionKeys := []int{9420, 9574, 9636}
    for _, key := range sessionKeys {
        url := fmt.Sprintf("https://api.openf1.org/v1/positions?session_key=%d", key)
        resp, err := http.Get(url)
        if err != nil {
            log.Printf("❌ Error al obtener posiciones para session_key %d: %v\n", key, err)
            continue
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()
        // The positions endpoint returns an array of {driver_number, position, date} for the session
        var entries []map[string]interface{}
        if err := json.Unmarshal(body, &entries); err != nil {
            log.Println("❌ Error decodificando posiciones:", err)
            continue
        }
        count := 0
        for _, entry := range entries {
            // Extract driver_number, position, date (if present)
            dn, ok1 := entry["driver_number"].(float64)
            pos, ok2 := entry["position"].(float64)
            dateVal, ok3 := entry["date"].(string)
            if !ok1 || !ok2 {
                continue
            }
            driverNum := int(dn)
            position := int(pos)
            date := ""
            if ok3 {
                date = dateVal
            }
            _, err := db.Exec(`INSERT INTO positions (driver_number, session_key, position, date) VALUES (?, ?, ?, ?)`,
                driverNum, key, position, date)
            if err != nil {
                log.Println("❌ Error insertando posición:", err)
            } else {
                count++
            }
        }
        log.Printf("➕ Posiciones insertadas para carrera %d: %d registros.\n", key, count)
    }
    log.Println("✅ Posiciones cargadas correctamente.")
}

func cargarVueltas(db *sql.DB) {
    log.Println("📡 Cargando vueltas de carreras...")
    _, _ = db.Exec("DELETE FROM laps")
    sessionKeys := []int{9420, 9574, 9636}
    for _, key := range sessionKeys {
        url := fmt.Sprintf("https://api.openf1.org/v1/laps?session_key=%d", key)
        resp, err := http.Get(url)
        if err != nil {
            log.Printf("❌ Error al obtener vueltas para session_key %d: %v\n", key, err)
            continue
        }
        var laps []models.Lap
        if err := json.NewDecoder(resp.Body).Decode(&laps); err != nil {
            log.Println("❌ Error decodificando vueltas:", err)
            resp.Body.Close()
            continue
        }
        resp.Body.Close()
        for _, lap := range laps {
            _, err := db.Exec(`INSERT INTO laps 
                (driver_number, session_key, lap_number, lap_duration, duration_sector_1, duration_sector_2, duration_sector_3, st_speed, date_start) 
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
                lap.DriverNumber, lap.SessionKey, lap.LapNumber, lap.LapDuration, lap.DurationSector1, lap.DurationSector2, lap.DurationSector3, lap.StSpeed, lap.DateStart)
            if err != nil {
                log.Println("❌ Error insertando vuelta:", err)
            }
        }
    }
    log.Println("✅ Vueltas cargadas correctamente.")
}

func cargarResultados(db *sql.DB) {
    log.Println("📡 Cargando resultados adicionales (vueltas rápidas y velocidades)...")
    _, _ = db.Exec("DELETE FROM results")
    // Using the same session keys for results
    sessionKeys := []int{9420, 9574, 9636}
    for _, key := range sessionKeys {
        url := fmt.Sprintf("https://api.openf1.org/v1/results?session_key=%d", key)
        resp, err := http.Get(url)
        if err != nil {
            log.Printf("❌ Error al obtener resultados para session_key %d: %v\n", key, err)
            continue
        }
        defer resp.Body.Close()
        type ResultAPI struct {
            SessionKey    int     `json:"session_key"`
            DriverNumber  int     `json:"driver_number"`
            FinalPosition int     `json:"final_position"`
            Laps          int     `json:"laps"`
            TopSpeed      float64 `json:"top_speed"`
        }
        var results []ResultAPI
        if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
            log.Println("❌ Error decodificando resultados:", err)
            continue
        }
        for _, r := range results {
            _, err := db.Exec(`INSERT INTO results (session_key, driver_number, final_position, laps, top_speed) VALUES (?, ?, ?, ?, ?)`,
                r.SessionKey, r.DriverNumber, r.FinalPosition, r.Laps, r.TopSpeed)
            if err != nil {
                log.Println("❌ Error insertando resultado:", err)
            }
        }
        log.Printf("🔍 Resultados cargados para carrera %d: %d registros.\n", key, len(results))
    }
    log.Println("✅ Resultados cargados correctamente.")
}

// contains checks if a value is in a slice of ints.
func contains(slice []int, val int) bool {
    for _, v := range slice {
        if v == val {
            return true
        }
    }
    return false
}
