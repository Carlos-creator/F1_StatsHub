package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"f1statshub/models"
)

func main() {
	for {
		fmt.Println("\n--- F1 StatsHub CLI ---")
		fmt.Println("1. Ver todos los corredores")
		fmt.Println("2. Ver detalle de un corredor")
		fmt.Println("3. Ver todas las carreras")
		fmt.Println("4. Ver resultados de una carrera")
		fmt.Println("5. Ver resumen de rendimiento")
		fmt.Println("6. Salir")
		fmt.Print("Selecciona una opción: ")

		var opcion int
		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			verCorredores()
		case 2:
			verDetalleCorredor()
		case 3:
			verCarreras()
		case 4:
			verResultadosCarrera()
		case 5:
			verResumen()
		case 6:
			fmt.Println("Saliendo...")
			os.Exit(0)
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func verCorredores() {
	resp, err := http.Get("http://localhost:8080/api/corredor")
	if err != nil {
		fmt.Println("Error al solicitar corredores:", err)
		return
	}
	defer resp.Body.Close()

	var corredores []models.Driver
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &corredores)

	fmt.Println("\n#  | Nombre             | Team             | País")
	fmt.Println("----+--------------------+------------------+-------")
	for _, c := range corredores {
		fmt.Printf("%-3d| %-18s| %-17s| %-6s\n", c.DriverNumber, c.FirstName+" "+c.LastName, c.TeamName, c.CountryCode)
	}
}

func verDetalleCorredor() {
	fmt.Print("Ingresa el número del corredor: ")
	var numero int
	fmt.Scanln(&numero)

	url := fmt.Sprintf("http://localhost:8080/api/corredor/detalle/%d", numero)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error al obtener detalle:", err)
		return
	}
	defer resp.Body.Close()

	var detalle models.DriverDetail
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &detalle)

	fmt.Println("\nNombre:", detalle.FullName)
	fmt.Println("Team:", detalle.Team)
	fmt.Println("Promedio Posición:", detalle.AvgPosition)
	fmt.Println("Velocidad Máxima:", detalle.MaxTopSpeed)
}

func verCarreras() {
	resp, err := http.Get("http://localhost:8080/api/carrera")
	if err != nil {
		fmt.Println("Error al solicitar carreras:", err)
		return
	}
	defer resp.Body.Close()

	var carreras []models.Carrera
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &carreras)

	fmt.Println("\nKey | Nombre                       | Lugar        | Año")
	fmt.Println("-----+------------------------------+--------------+------")
	for _, c := range carreras {
		fmt.Printf("%-4d| %-28s| %-13s| %-4d\n", c.SessionKey, c.SessionName, c.Location, c.Year)
	}
}

func verResultadosCarrera() {
	fmt.Print("Ingresa el SessionKey de la carrera: ")
	var key int
	fmt.Scanln(&key)

	url := fmt.Sprintf("http://localhost:8080/api/carrera/detalle/%d", key)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error al obtener resultados:", err)
		return
	}
	defer resp.Body.Close()

	var resultados []models.RaceResult
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &resultados)

	fmt.Println("\nPos | Nombre               | Vueltas | Top Speed")
	fmt.Println("-----+----------------------+---------+-----------")
	for _, r := range resultados {
		fmt.Printf("%-4d| %-21s| %-8d| %.2f km/h\n", r.FinalPosition, r.FullName, r.Laps, r.TopSpeed)
	}
}

func verResumen() {
	resp, err := http.Get("http://localhost:8080/api/resumen")
	if err != nil {
		fmt.Println("Error al solicitar resumen:", err)
		return
	}
	defer resp.Body.Close()

	var resumen []models.PerformanceSummary
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &resumen)

	fmt.Println("\n#  | Prom. Posición | Prom. Velocidad")
	fmt.Println("----+----------------+------------------")
	for _, r := range resumen {
		fmt.Printf("%-3d| %-14.2f| %.2f km/h\n", r.DriverNumber, r.AvgPosition, r.AvgSpeed)
	}
}