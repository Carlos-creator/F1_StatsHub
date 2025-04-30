package endpoints

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"f1statshub/models"
	"sort"
	"strings"
)

func getResumen(c *gin.Context) {
	// Top 3 ganadores (victorias)
	victorias := make(map[string]int)
	rows, _ := dbConn.Query(`
		SELECT d.first_name, d.last_name, COUNT(*) 
		FROM positions p 
		JOIN drivers d ON p.driver_number = d.driver_number 
		JOIN sessions s ON p.session_key = s.session_key 
		WHERE s.session_type = 'Race' AND p.position = 1 
		GROUP BY p.driver_number`)
	for rows.Next() {
		var first, last string
		var count int
		rows.Scan(&first, &last, &count)
		nombre := strings.TrimSpace(first + " " + last)
		victorias[nombre] = count
	}
	rows.Close()

	// Top 3 vueltas rápidas
	vueltasRapidas := make(map[string]int)
	rows, _ = dbConn.Query(`
		SELECT d.first_name, d.last_name, COUNT(*) 
		FROM laps l 
		JOIN drivers d ON l.driver_number = d.driver_number 
		JOIN sessions s ON l.session_key = s.session_key 
		WHERE s.session_type = 'Race' 
		AND l.lap_duration = (
			SELECT MIN(l2.lap_duration) 
			FROM laps l2 
			WHERE l2.session_key = s.session_key
		) 
		GROUP BY l.driver_number`)
	for rows.Next() {
		var first, last string
		var count int
		rows.Scan(&first, &last, &count)
		nombre := strings.TrimSpace(first + " " + last)
		vueltasRapidas[nombre] = count
	}
	rows.Close()

	// Top 3 poles
	poles := make(map[string]int)
	rows, _ = dbConn.Query(`
		SELECT d.first_name, d.last_name, COUNT(*) 
		FROM positions p 
		JOIN drivers d ON p.driver_number = d.driver_number 
		JOIN sessions s ON p.session_key = s.session_key 
		WHERE s.session_name = 'Qualifying' AND p.position = 1 
		GROUP BY p.driver_number`)
	for rows.Next() {
		var first, last string
		var count int
		rows.Scan(&first, &last, &count)
		nombre := strings.TrimSpace(first + " " + last)
		poles[nombre] = count
	}
	rows.Close()

	// Función para convertir map en []Stat ordenado
	getTop3 := func(m map[string]int) []models.Stat {
		stats := []models.Stat{}
		for name, val := range m {
			parts := strings.Split(name, " ")
			first := parts[0]
			last := strings.Join(parts[1:], " ")
			stats = append(stats, models.Stat{
				FirstName: first,
				LastName:  last,
				Value:     float64(val),
			})
		}
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].Value > stats[j].Value
		})
		if len(stats) > 3 {
			stats = stats[:3]
		}
		for i := range stats {
			stats[i].Position = i + 1
		}
		return stats
	}

	// Crear el resumen con solo los campos válidos
	summary := models.Summary{
		Top3Winners:       getTop3(victorias),
		Top3FastestLaps:   getTop3(vueltasRapidas),
		Top3PolePositions: getTop3(poles),
	}

	c.JSON(http.StatusOK, summary)
}
