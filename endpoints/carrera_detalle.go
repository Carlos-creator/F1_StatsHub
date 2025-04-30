package endpoints

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "f1statshub/models"
)

// getCarreraDetail handles GET /api/carrera/detalle/:id - returns details (results) for a specific race.
func getCarreraDetail(c *gin.Context) {
    id := c.Param("id")
    // Query race session info (name/location)
    var raceName string
    _ = dbConn.QueryRow("SELECT location, country_name FROM sessions WHERE session_key = ?", id).Scan(&raceName, new(string))
    if raceName == "" {
        // If location is empty or not found, try country_name
        dbConn.QueryRow("SELECT country_name FROM sessions WHERE session_key = ?", id).Scan(&raceName)
    }
    if raceName == "" {
        raceName = "Carrera " + id
    }

    // Query all drivers results for this race
    query := `
        SELECT d.first_name, d.last_name, d.team_name, p.position, MAX(l.st_speed), MIN(l.lap_duration)
        FROM positions p
        JOIN drivers d ON p.driver_number = d.driver_number
        JOIN laps l ON l.session_key = p.session_key AND l.driver_number = p.driver_number
        WHERE p.session_key = ?
        GROUP BY p.driver_number
        ORDER BY p.position ASC
    `
    rows, err := dbConn.Query(query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
        return
    }
    defer rows.Close()

    var results []models.DriverRaceResult
    var globalFastest float64
    // Determine the fastest lap of the race (global)
    _ = dbConn.QueryRow("SELECT MIN(lap_duration) FROM laps WHERE session_key = ?", id).Scan(&globalFastest)
    for rows.Next() {
        var res models.DriverRaceResult
        err := rows.Scan(&res.FirstName, &res.LastName, &res.TeamName, &res.Position, &res.MaxSpeed, &res.BestLapDuration)
        if err != nil {
            continue
        }
        // We need driver number for completeness
        // Fetch driver number (since it's not selected in the aggregate query)
        _ = dbConn.QueryRow("SELECT driver_number FROM drivers WHERE first_name = ? AND last_name = ? AND team_name = ?", res.FirstName, res.LastName, res.TeamName).Scan(&res.DriverNumber)
        res.FastestLap = (globalFastest > 0 && res.BestLapDuration == globalFastest)
        results = append(results, res)
    }
    detail := models.RaceDetail{
        RaceName:    "GP de " + raceName,
        RaceResults: results,
    }
    c.JSON(http.StatusOK, detail)
}
