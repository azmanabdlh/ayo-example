package model

import (
	"time"

	teamManagement "github.com/azmanabdlh/ayo-example/internal/team-management/model"
	"gorm.io/gorm"
)

type Match struct {
	gorm.Model

	ID int64 `json:"id,omitempty" gorm:"primaryKey"`

	MatchDate time.Time `json:"match_date,omitempty"`
	// KickOffAt time.Time optional

	Title string `json:"title,omitempty"`

	HomeTeamID        int64               `json:"home_team_id,omitempty" gorm:"foreignKey:HomeTeamID"`
	HomeTeam          teamManagement.Team `json:"home_team,omitempty"`
	HomeTeamFormation string              `json:"home_team_formation,omitempty"`

	AwayTeamID        int64               `json:"away_team_id,omitempty" gorm:"foreignKey:AwayTeamID"`
	AwayTeam          teamManagement.Team `json:"away_team,omitempty"`
	AwayTeamFormation string              `json:"away_team_formation,omitempty"`

	HomeScore int `json:"home_score,omitempty"`
	AwayScore int `json:"away_score,omitempty"`

	Phase int `json:"phase,omitempty"` // phase: 1 [active], 2 [cancelled], 3: [finished]

	Goal []Goal `json:"goals,omitempty"`

	PlayerLineup []MatchPlayerLineup `json:"player_lineup,omitempty"`

	VenueID   int64  `json:"venue_id,omitempty" gorm:"foreignKey:VenueID"`
	VenueName string `json:"venue_name,omitempty"`

	Venue Venue `json:"venue,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
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

	// TODO:
	// 1. validate valid match teamID
	// 2. validate valid player on this match
	return m.ValidTeam(teamID) && m.HasPlayer(playerID)

}

func (m *Match) HasPlayer(playerID int64) bool {
	for _, player := range m.PlayerLineup {
		isPlaying := player.ID == playerID && player.IsStarter
		if isPlaying {
			return true
		}
	}

	return false
}

func (m *Match) ValidTeam(teamID int64) bool {
	return teamID == m.HomeTeamID || teamID == m.AwayTeamID
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

	highlight.Scored = make(map[string]int)
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

			player.SubstitutedForPlayerID = *lineup.SubstitutedForPlayerID
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

	ID int64 `json:"id,omitempty" gorm:"primaryKey"`

	Match   Match `json:"match,omitempty"`
	MatchID int64 `json:"match_id,omitempty" gorm:"foreignKey:MatchID;uniqueIndex:idx_match_team_player"`

	Team   teamManagement.Team `json:"team,omitempty"`
	TeamID int64               `json:"team_id,omitempty" gorm:"foreignKey:TeamID;uniqueIndex:idx_match_team_player" `

	// ["ST", "GK", "CB-L"]
	PositionSlot string `json:"position_slot,omitempty"`

	IsStarter bool `json:"is_starter,omitempty"`

	MinuteIn  int `json:"minute_in,omitempty"`
	MinuteOut int `json:"minute_out,omitempty"`

	Player   teamManagement.Player `json:"player,omitempty"`
	PlayerID int64                 `json:"player_id,omitempty" gorm:"foreignKey:PlayerID;uniqueIndex:idx_match_team_player"`

	SubstitutedForPlayer   teamManagement.Player `json:"substituted_for_player,omitempty"`
	SubstitutedForPlayerID *int64                `json:"substituted_for_player_id,omitempty" gorm:"foreignKey:SubstitutedForPlayerID"`

	Reason string `json:"reason,omitempty"`
}

func (MatchPlayerLineup) TableName() string {
	return "match_players"
}

type Goal struct {
	gorm.Model

	ID int64 `json:"id,omitempty" gorm:"primaryKey"`

	MatchID int64 `json:"match_id,omitempty"`
	Match   Match `json:"match,omitempty"`

	PlayerID int64                 `json:"player_id,omitempty" gorm:"foreignKey:PlayerID"`
	Player   teamManagement.Player `json:"player,omitempty"`

	TeamID int64               `json:"team_id,omitempty" gorm:"foreignKey:TeamID"`
	Team   teamManagement.Team `json:"team,omitempty"`

	ScoredAtMinute int `json:"scored_at_minute,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type MatchHighlight struct {
	MatchID int64                    `json:"match_id,omitempty"`
	Team    map[string]TeamHighlight `json:"team,omitempty"`
	Scored  map[string]int           `json:"scored,omitempty"`
	Goal    []GoalHighlight          `json:"goals,omitempty"`
	Phase   int                      `json:"phase,omitempty"`
	Venue   Venue                    `json:"venue,omitempty"`
}

type TeamHighlight struct {
	TeamID    int64  `json:"team_id,omitempty"`
	Name      string `json:"name,omitempty"`
	LogoURL   string `json:"logo_url,omitempty"`
	Formation string `json:"formation,omitempty"`

	Player []PlayerHighlight `json:"players,omitempty"`

	PlayerLineupPosition []LineupPosition `json:"player_lineup_positions,omitempty"`
}

type GoalHighlight struct {
	PlayerID       int64  `json:"player_id,omitempty"`
	PlayerName     string `json:"player_name,omitempty"`
	ScoredAtMinute int    `json:"scored_at_minute,omitempty"`
}

type PlayerHighlight struct {
	PlayerID   int64  `json:"player_id,omitempty"`
	PlayerName string `json:"player_name,omitempty"`
	Position   string `json:"position,omitempty"`
	BackNumber int    `json:"back_number,omitempty"`

	IsStarter bool `json:"is_starter,omitempty"`

	MinuteIn  int `json:"minute_in,omitempty"`
	MinuteOut int `json:"minute_out,omitempty"`

	SubstitutedForPlayerID   int64  `json:"substituted_for_player_id,omitempty"`
	SubstitutedForPlayerName string `json:"substituted_for_player_name,omitempty"`

	Reason string `json:"reason,omitempty"`
}

type LineupPosition struct {
	PlayerID     int64  `json:"player_id,omitempty"`
	PositionSlot string `json:"position_slot,omitempty"`

	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}
