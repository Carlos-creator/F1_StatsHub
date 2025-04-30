package endpoints

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"f1statshub/models"
	"strconv"
)

// GET /api/corredor
func getDrivers(c *gin.Context) {
	rows, err := dbConn.Query("SELECT driver_number, first_name, last_name, name_acronym, team_name, country_code FROM drivers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		rows.Scan(&d.DriverNumber, &d.FirstName, &d.LastName, &d.NameAcronym, &d.TeamName, &d.CountryCode)
		drivers = append(drivers, d)
	}
	c.JSON(http.StatusOK, drivers)
}

// GET /api/corredor/detalle/:id
func getDriverDetail(c *gin.Context) {
	id := c.Param("id")
	driverID := atoiOrZero(id)

	// Obtener info del piloto
	var driver models.Driver
	err := dbConn.QueryRow("SELECT driver_number, first_name, last_name, name_acronym, team_name, country_code FROM drivers WHERE driver_number = ?", driverID).
		Scan(&driver.DriverNumber, &driver.FirstName, &driver.LastName, &driver.NameAcronym, &driver.TeamName, &driver.CountryCode)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	query := `
		SELECT s.session_key, s.circuit_short_name, s.session_name,
			   p.position, MAX(l.st_speed), MIN(l.lap_duration)
		FROM positions p
		JOIN sessions s ON p.session_key = s.session_key
		JOIN laps l ON l.session_key = p.session_key AND l.driver_number = p.driver_number
		WHERE p.driver_number = ?
		GROUP BY s.session_key
	`
	rows, err := dbConn.Query(query, driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}
	defer rows.Close()

	var results []models.RaceResult
	wins := 0
	top3 := 0
	var maxSpeedOverall float64

	for rows.Next() {
		var r models.RaceResult
		err := rows.Scan(&r.SessionKey, &r.CircuitShortName, &r.Race, &r.Position, &r.MaxSpeed, &r.BestLapDuration)
		if err != nil {
			continue
		}
		if r.Position == 1 {
			wins++
		}
		if r.Position <= 3 {
			top3++
		}
		if r.MaxSpeed > maxSpeedOverall {
			maxSpeedOverall = r.MaxSpeed
		}
		// Verificar vuelta más rápida global de la sesión
		var globalFast float64
		err = dbConn.QueryRow("SELECT MIN(lap_duration) FROM laps WHERE session_key = ?", r.SessionKey).Scan(&globalFast)
		r.FastestLap = (err == nil && r.BestLapDuration == globalFast)
		results = append(results, r)
	}

	// Resumen de rendimiento
	perf := models.PerformanceSummary{
		Wins:         wins,
		Top3Finishes: top3,
		MaxSpeed:     maxSpeedOverall,
	}

	c.JSON(http.StatusOK, models.DriverDetail{
		Driver:             driver,
		PerformanceSummary: perf,
		RaceResults:        results,
	})
}

// convierte string a int
func atoiOrZero(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
