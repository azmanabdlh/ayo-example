package model

import (
	"time"

	teamManagement "github.com/azmanabdlh/ayo-example/internal/team-management/model"
	"gorm.io/gorm"
)

type Match struct {
	gorm.Model

	ID int64 `gorm:"primaryKey"`

	MatchDate time.Time
	// KickOffAt time.Time optional

	Title string

	HomeTeamID        int64 `gorm:"foreignKey:HomeTeamID"`
	HomeTeam          teamManagement.Team
	HomeTeamFormation string

	AwayTeamID        int64 `gorm:"foreignKey:AwayTeamID"`
	AwayTeam          teamManagement.Team
	AwayTeamFormation string

	HomeScore int
	AwayScore int

	Phase int // phase: 1 [active], 2 [cancelled], 3: [finished]

	Goal []Goal

	PlayerLineup []MatchPlayerLineup

	VenueID   int64 `gorm:"foreignKey:VenueID"`
	VenueName string

	Venue Venue

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Match) IsUpcoming() bool {
	return time.Now().Before(m.MatchDate) && m.Phase == 1
}

func (m *Match) IsOngoing() bool {
	return time.Now().After(m.MatchDate) && m.Phase == 1
}

func (m *Match) IsFinished() bool {
	return time.Now().After(m.MatchDate) && m.Phase == 3
}

func (m *Match) IsCancelled() bool {
	return m.Phase == 2
}

func (m *Match) CanRecordGoal(playerID, teamID int64) bool {

	if !m.IsOngoing() {
		return false
	}

	// validTeamID && validPlayerID

	return m.IsParticipating(teamID) &&
		(m.IsHomePlayer(playerID) || m.IsAwayPlayer(playerID))

}

func (m *Match) IsParticipating(teamID int64) bool {
	return teamID == m.HomeTeamID || teamID == m.AwayTeamID
}

func (m *Match) IsHomePlayer(playerID int64) bool {
	return m.HomeTeam.HasPlayerID(playerID)
}

func (m *Match) IsAwayPlayer(playerID int64) bool {
	return m.AwayTeam.HasPlayerID(playerID)
}

func (m *Match) ToMatchHighlight() *MatchHighlight {

	highlight := new(MatchHighlight)

	highlight.MatchID = m.ID
	highlight.Phase = m.Phase
	highlight.Venue = m.Venue

	for _, n := range m.Goal {
		highlight.Goal = append(highlight.Goal, GoalHighlight{
			PlayerID:       n.ID,
			PlayerName:     n.Player.Name,
			ScoredAtMinute: n.ScoredAtMinute,
		})
	}

	playerLineup := make(map[int64]MatchPlayerLineup)

	for _, n := range m.PlayerLineup {
		playerLineup[n.ID] = n
	}

	highlight.Team = make(map[string]TeamHighlight)
	highlight.Team["home"] = newTeamHighlight(
		m.HomeTeam,
		playerLineup,
		m.HomeTeamFormation,
		0,
	)
	highlight.Team["away"] = newTeamHighlight(
		m.AwayTeam,
		playerLineup,
		m.AwayTeamFormation,
		100,
	)

	highlight.Scored["home"] = m.HomeScore
	highlight.Scored["away"] = m.AwayScore

	return highlight
}

func newTeamHighlight(
	team teamManagement.Team,
	playerLineup map[int64]MatchPlayerLineup,
	formation string,
	z int,
) TeamHighlight {
	var (
		players        []PlayerHighlight
		lineupPosition []LineupPosition
	)

	for _, n := range team.Player {

		player := PlayerHighlight{
			PlayerID:   n.ID,
			PlayerName: n.Name,
			Position:   n.Position,
			BackNumber: n.BackNumber,
		}

		lineup, ok := playerLineup[n.ID]
		if ok {
			player.MinuteIn = lineup.MinuteIn
			player.MinuteOut = lineup.MinuteOut
			player.IsStarter = lineup.IsStarter

			player.SubstitutedForPlayerID = lineup.SubstitutedForPlayerID
			player.SubstitutedForPlayerName = lineup.SubstitutedForPlayer.Name

			player.Reason = lineup.Reason

			// set player lineup

			position := slotCoordinates[lineup.PositionSlot]

			x := position.X
			// mirror position for away
			if z > 0 {
				x = z - x
			}

			lineupPosition = append(lineupPosition, LineupPosition{
				PlayerID:     n.ID,
				PositionSlot: lineup.PositionSlot,

				X: x,
				Y: position.Y,
			})
		}

		players = append(players, player)
	}

	return TeamHighlight{
		TeamID:  team.ID,
		Name:    team.Name,
		LogoURL: team.LogoURL,

		Formation:            formation,
		Player:               players,
		PlayerLineupPosition: lineupPosition,
	}
}

type MatchPlayerLineup struct {
	gorm.Model

	ID int64 `gorm:"primaryKey"`

	Match   Match
	MatchID int64 `gorm:"foreignKey:MatchID"`

	Team   teamManagement.Team
	TeamID int64 `gorm:"foreignKey:TeamID"`

	// ["ST", "GK", "CB-L"]
	PositionSlot string `gorm:"uniqueIndex"`

	IsStarter bool `gorm:"uniqueIndex"`

	MinuteIn  int
	MinuteOut int

	Player   teamManagement.Player
	PlayerID int64 `gorm:"foreignKey:PlayerID"`

	SubstitutedForPlayer   teamManagement.Player
	SubstitutedForPlayerID int64 `gorm:"foreignKey:SubstitutedForPlayerID"`

	Reason string
}

func (MatchPlayerLineup) TableName() string {
	return "match_players"
}

type Goal struct {
	gorm.Model

	ID int64 `gorm:"primaryKey"`

	MatchID int64
	Match   Match

	PlayerID int64 `gorm:"foreignKey:PlayerID"`
	Player   teamManagement.Player

	TeamID int64 `gorm:"foreignKey:TeamID"`
	Team   teamManagement.Team

	ScoredAtMinute int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type MatchHighlight struct {
	MatchID int64
	Team    map[string]TeamHighlight
	Scored  map[string]int
	Goal    []GoalHighlight
	Phase   int
	Venue   Venue
}

type TeamHighlight struct {
	TeamID    int64
	Name      string
	LogoURL   string
	Formation string

	Player []PlayerHighlight

	PlayerLineupPosition []LineupPosition
}

type GoalHighlight struct {
	PlayerID       int64
	PlayerName     string
	ScoredAtMinute int
}

type PlayerHighlight struct {
	PlayerID   int64
	PlayerName string
	Position   string
	BackNumber int

	IsStarter bool

	MinuteIn  int
	MinuteOut int

	SubstitutedForPlayerID   int64
	SubstitutedForPlayerName string

	Reason string
}

type LineupPosition struct {
	PlayerID     int64
	PositionSlot string

	X int
	Y int
}
