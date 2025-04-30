package models

type Driver struct {
	DriverNumber int    `json:"driver_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	NameAcronym  string `json:"name_acronym"`
	TeamName     string `json:"team_name"`
	CountryCode  string `json:"country_code"`
}

type Carrera struct {
	SessionKey       int    `json:"session_key"`
	SessionName      string `json:"session_name"`
	SessionType      string `json:"session_type"`
	Location         string `json:"location"`
	CountryName      string `json:"country_name"`
	Year             int    `json:"year"`
	CircuitShortName string `json:"circuit_short_name"`
	Date             string `json:"date_start"`
	Circuit          string `json:"circuit"` // si usas esto explícitamente
}

type RaceResult struct {
	SessionKey       int     `json:"session_key"`
	DriverNumber     int     `json:"driver_number"`
	FinalPosition    int     `json:"final_position"`
	Laps             int     `json:"laps"`
	TopSpeed         float64 `json:"top_speed"`
	BestLapDuration  float64 `json:"best_lap_duration"`
	CircuitShortName string  `json:"circuit_short_name"`
	Race             string  `json:"race"`
	Position         int     `json:"position"`
	MaxSpeed         float64 `json:"max_speed"`
	FastestLap       bool    `json:"fastest_lap"`
}

type Lap struct {
	DriverNumber     int     `json:"driver_number"`
	SessionKey       int     `json:"session_key"`
	LapNumber        int     `json:"lap_number"`
	LapDuration      float64 `json:"lap_duration"`
	DurationSector1  float64 `json:"duration_sector_1"`
	DurationSector2  float64 `json:"duration_sector_2"`
	DurationSector3  float64 `json:"duration_sector_3"`
	StSpeed          float64 `json:"st_speed"`
	DateStart        string  `json:"date_start"`
}

type Position struct {
	DriverNumber int    `json:"driver_number"`
	SessionKey   int    `json:"session_key"`
	Position     int    `json:"position"`
	Date         string `json:"date"`
}

type DriverRaceResult struct {
	DriverNumber     int     `json:"driver_number"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	TeamName         string  `json:"team_name"`
	Position         int     `json:"position"`
	MaxSpeed         float64 `json:"max_speed"`
	BestLapDuration  float64 `json:"best_lap_duration"`
	FastestLap       bool    `json:"fastest_lap"`
}

type RaceDetail struct {
	RaceName    string             `json:"race_name"`
	RaceResults []DriverRaceResult `json:"race_results"`
}

type PerformanceSummary struct {
	MaxSpeed     float64 `json:"max_speed"`
	Top3Finishes int     `json:"top_3_finishes"`
	Wins         int     `json:"wins"`
}

type DriverSummary struct {
	DriverID           int                `json:"driver_id"`
	PerformanceSummary PerformanceSummary `json:"performance_summary"`
	RaceResults        []RaceResult       `json:"race_results"`
}

type DriverDetail struct {
	Driver            Driver            `json:"driver"`
	PerformanceSummary PerformanceSummary `json:"performance_summary"`
	RaceResults       []RaceResult      `json:"race_results"`
}


type Stat struct {
	DriverNumber int     `json:"driver_number"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Value        float64 `json:"value"`
	Position     int     `json:"position"`
}

type Summary struct {
	TopAvgPosition   []Stat `json:"top_avg_position"`
	TopSpeed         []Stat `json:"top_speed"`
	Top3Winners      []Stat `json:"top_3_winners"`
	Top3FastestLaps  []Stat `json:"top_3_fastest_laps"`
	Top3PolePositions []Stat `json:"top_3_pole_positions"`
}
