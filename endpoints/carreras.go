package endpoints

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "f1statshub/models"
)

// getCarreras handles GET /api/carrera - returns a list of all races (sessions).
func getCarreras(c *gin.Context) {
    rows, err := dbConn.Query("SELECT session_key, country_name, date_start, year, circuit_short_name FROM sessions WHERE session_type = 'Race'")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
        return
    }
    defer rows.Close()

    var carreras []models.Carrera
    for rows.Next() {
        var race models.Carrera
        rows.Scan(&race.SessionKey, &race.CountryName, &race.Date, &race.Year, &race.Circuit)
        carreras = append(carreras, race)
    }
    c.JSON(http.StatusOK, carreras)
}
