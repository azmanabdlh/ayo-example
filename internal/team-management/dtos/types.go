package dtos

type TeamParam struct {
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url"`
	FoundedYear int    `json:"founded_year"`
	Address     string `json:"address"`
	City        string `json:"city"`
}

type PlayerParam struct {
	Name       string `json:"name"`
	Height     int    `json:"height"`
	Weight     int    `json:"weight"`
	Position   string `json:"position"`
	BackNumber int    `json:"back_number"`
	TeamID     int64  `json:"team_id"`
}
