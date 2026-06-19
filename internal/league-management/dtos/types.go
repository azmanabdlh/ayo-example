package dtos

import (
	"time"
)

type MatchParam struct {
	MatchDate time.Time `json:"match_date"`

	Title string `json:"title"`

	HomeTeamID int64 `json:"home_team_id"`
	AwayTeamID int64 `json:"away_team_id"`

	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`

	Phase int `json:"phase"` // phase: 1 [active], 2 [cancelled], 3: [finished]

	VenueID   int64  `json:"venue_id" gorm:"foreignKey:VenueID"`
	VenueName string `json:"venue_name"`
}

type GoalParam struct {
	PlayerID int64 `json:"player_id"`
	TeamID   int64 `json:"team_id"`

	ScoredAtMinute int `json:"scored_at_minute"`
}

type VenueParam struct {
	Name          string `json:"name"`
	City          string `json:"city"`
	Address       string `json:"address"`
	Capacity      int    `json:"capacity"`
	GoogleMapsURL string `json:"google_maps_url"`
}

type SubstitutePlayerParam struct {
	TeamID                 int64  `json:"team_id"`
	Minute                 int    `json:"minute"`
	PlayerID               int64  `json:"player_id"`
	SubstitutedForPlayerID int64  `json:"substituted_for_player_id"`
	Reason                 string `json:"reason"`
}

type PlayerLineup struct {
	PositionSlot string `json:"position_slot"` // ["ST", "GK", "CB-L"]
	IsStarter    bool   `json:"is_starter"`

	PlayerID int64 `json:"player_id"`
}

type MatchPlayerParam struct {
	TeamID int64 `json:"team_id"`

	Lineup []PlayerLineup `json:"lineup"`
}
