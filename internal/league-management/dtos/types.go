package dtos

import (
	"time"
)

type MatchParam struct {
	MatchDate time.Time

	Title string

	HomeTeamID int64
	AwayTeamID int64

	HomeScore int
	AwayScore int

	Phase int // phase: 1 [active], 2 [cancelled], 3: [finished]

	VenueID   int64 `gorm:"foreignKey:VenueID"`
	VenueName string
}

type GoalParam struct {
	PlayerID int64
	TeamID   int64

	ScoredAtMinute int
}

type VenueParam struct {
	Name          string
	Address       string
	City          string
	Capacity      int
	GoogleMapsURL string
}

type SubstitutePlayerParam struct {
	TeamID                 int64
	Minute                 int
	PlayerID               int64
	SubstitutedForPlayerID int64
	Reason                 string
}

type PlayerLineup struct {
	PositionSlot string // ["ST", "GK", "CB-L"]
	IsStarter    bool

	PlayerID int64
}

type MatchPlayerParam struct {
	TeamID int64

	Lineup []PlayerLineup
}
