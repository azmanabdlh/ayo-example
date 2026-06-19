package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/azmanabdlh/ayo-example/internal/league-management/handler"
	"github.com/azmanabdlh/ayo-example/internal/league-management/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_SubstitutePlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		id string

		body string

		mockFn func(*handler.MockService)

		expectedCode int
	}{
		{
			name: "success",
			id:   "10",
			body: `{
                "player_out_id":1,
                "player_in_id":2
            }`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"SubstitutePlayer",
					mock.Anything,
					int64(10),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name:         "invalid json",
			id:           "10",
			body:         `{`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "service error",
			id:   "10",
			body: `{
                "player_out_id":1,
                "player_in_id":2
            }`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"SubstitutePlayer",
					mock.Anything,
					int64(10),
					mock.Anything,
				).Return(errors.New("database error"))
			},
			expectedCode: 500,
		},
		{
			name: "invalid id become zero",
			id:   "abc",
			body: `{
                "player_out_id":1,
                "player_in_id":2
            }`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"SubstitutePlayer",
					mock.Anything,
					int64(0),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := handler.NewMockService(t)
			tt.mockFn(svc)

			h := handler.New(svc)

			r := gin.New()
			r.POST("/:id/substitute", h.SubstitutePlayer)

			req := httptest.NewRequest(
				http.MethodPost,
				"/"+tt.id+"/substitute",
				strings.NewReader(tt.body),
			)

			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(
				t,
				tt.expectedCode,
				rec.Code,
			)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_Finish(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		id string

		mockFn func(*handler.MockService)

		expectedCode int
	}{
		{
			name: "success",
			id:   "10",
			mockFn: func(s *handler.MockService) {
				s.On(
					"Finish",
					mock.Anything,
					int64(10),
				).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name: "service error",
			id:   "10",
			mockFn: func(s *handler.MockService) {
				s.On(
					"Finish",
					mock.Anything,
					int64(10),
				).Return(errors.New("match already finished"))
			},
			expectedCode: 401,
		},
		{
			name: "invalid id become zero",
			id:   "abc",
			mockFn: func(s *handler.MockService) {
				s.On(
					"Finish",
					mock.Anything,
					int64(0),
				).Return(nil)
			},
			expectedCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := handler.NewMockService(t)
			tt.mockFn(svc)

			h := handler.New(svc)

			r := gin.New()
			r.POST("/:id/finish", h.Finish)

			req := httptest.NewRequest(
				http.MethodPost,
				"/"+tt.id+"/finish",
				nil,
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_AssignMatchPlayerLineup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		id string

		body string

		mockFn func(*handler.MockService)

		expectedCode int
	}{
		{
			name: "success",
			id:   "10",
			body: `{
				"team_id":1,
				"lineup":[
					{
						"position_slot":"GK",
						"is_starter":true,
						"player_id":1
					},
					{
						"position_slot":"CB-L",
						"is_starter":true,
						"player_id":2
					},
					{
						"position_slot":"ST",
						"is_starter":true,
						"player_id":9
					}
				]
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AssignMatchPlayerLineup",
					mock.Anything,
					int64(10),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name:         "invalid json",
			id:           "10",
			body:         `{`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "service error",
			id:   "10",
			body: `{
				"team_id":1,
				"lineup":[
					{
						"position_slot":"GK",
						"is_starter":true,
						"player_id":1
					}
				]
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AssignMatchPlayerLineup",
					mock.Anything,
					int64(10),
					mock.Anything,
				).Return(errors.New("database error"))
			},
			expectedCode: 500,
		},
		{
			name: "invalid id become zero",
			id:   "abc",
			body: `{
				"team_id":1,
				"lineup":[
					{
						"position_slot":"GK",
						"is_starter":true,
						"player_id":1
					}
				]
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AssignMatchPlayerLineup",
					mock.Anything,
					int64(0),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name: "empty body",
			id:   "10",
			body: `{}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AssignMatchPlayerLineup",
					mock.Anything,
					int64(10),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name: "invalid player_id type",
			id:   "10",
			body: `{
				"team_id":1,
				"lineup":[
					{
						"position_slot":"GK",
						"is_starter":true,
						"player_id":"abc"
					}
				]
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "invalid is_starter type",
			id:   "10",
			body: `{
				"team_id":1,
				"lineup":[
					{
						"position_slot":"GK",
						"is_starter":"true",
						"player_id":1
					}
				]
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "unknown field ignored",
			id:   "10",
			body: `{
				"team_id":1,
				"lineup":[
					{
						"position_slot":"GK",
						"is_starter":true,
						"player_id":1,
						"name":"courtois"
					}
				]
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AssignMatchPlayerLineup",
					mock.Anything,
					int64(10),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := handler.NewMockService(t)
			tt.mockFn(svc)

			h := handler.New(svc)

			r := gin.New()
			r.POST("/:id/lineup", h.AssignMatchPlayerLineup)

			req := httptest.NewRequest(
				http.MethodPost,
				"/"+tt.id+"/lineup",
				strings.NewReader(tt.body),
			)

			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(
				t,
				tt.expectedCode,
				rec.Code,
			)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_AddRecordGoal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		id string

		body string

		mockFn func(*handler.MockService)

		expectedCode int
	}{
		{
			name: "success",
			id:   "1",
			body: `{
				"player_id":10,
				"team_id":2,
				"scored_at_minute":90
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AddRecordGoal",
					mock.Anything,
					int64(1),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name:         "invalid json",
			id:           "1",
			body:         `{`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "service error",
			id:   "1",
			body: `{
				"player_id":10,
				"team_id":2,
				"scored_at_minute":90
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AddRecordGoal",
					mock.Anything,
					int64(1),
					mock.Anything,
				).Return(errors.New("database error"))
			},
			expectedCode: 500,
		},
		{
			name: "invalid id become zero",
			id:   "abc",
			body: `{
				"player_id":10,
				"team_id":2,
				"scored_at_minute":90
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AddRecordGoal",
					mock.Anything,
					int64(0),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name: "empty body",
			id:   "1",
			body: `{}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AddRecordGoal",
					mock.Anything,
					int64(1),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name: "unknown field ignored",
			id:   "1",
			body: `{
				"player_id":10,
				"team_id":2,
				"scored_at_minute":90,
				"name":"ronaldo"
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AddRecordGoal",
					mock.Anything,
					int64(1),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name: "invalid player_id type",
			id:   "1",
			body: `{
				"player_id":"10",
				"team_id":2,
				"scored_at_minute":90
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "invalid team_id type",
			id:   "1",
			body: `{
				"player_id":10,
				"team_id":"2",
				"scored_at_minute":90
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "invalid scored_at_minute type",
			id:   "1",
			body: `{
				"player_id":10,
				"team_id":2,
				"scored_at_minute":"90"
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "null values",
			id:   "1",
			body: `{
				"player_id":null,
				"team_id":null,
				"scored_at_minute":null
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AddRecordGoal",
					mock.Anything,
					int64(1),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name: "negative minute",
			id:   "1",
			body: `{
				"player_id":10,
				"team_id":2,
				"scored_at_minute":-1
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AddRecordGoal",
					mock.Anything,
					int64(1),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name: "extra large minute",
			id:   "1",
			body: `{
				"player_id":10,
				"team_id":2,
				"scored_at_minute":999
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"AddRecordGoal",
					mock.Anything,
					int64(1),
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := handler.NewMockService(t)
			tt.mockFn(svc)

			h := handler.New(svc)

			r := gin.New()
			r.POST("/:id/goal", h.AddRecordGoal)

			req := httptest.NewRequest(
				http.MethodPost,
				"/"+tt.id+"/goal",
				strings.NewReader(tt.body),
			)

			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(
				t,
				tt.expectedCode,
				rec.Code,
			)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_CreateMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		body string

		mockFn func(*handler.MockService)

		expectedCode int
	}{
		{
			name: "success",
			body: `{
				"match_date":"2026-06-19T19:00:00Z",
				"title":"Final Liga",
				"home_team_id":1,
				"away_team_id":2,
				"home_score":0,
				"away_score":0,
				"phase":1,
				"venue_id":1,
				"venue_name":"Gelora Bung Karno"
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"CreateMatch",
					mock.Anything,
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name:         "invalid json",
			body:         `{`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "service error",
			body: `{
				"match_date":"2026-06-19T19:00:00Z",
				"title":"Final Liga",
				"home_team_id":1,
				"away_team_id":2,
				"home_score":0,
				"away_score":0,
				"phase":1,
				"venue_id":1,
				"venue_name":"Gelora Bung Karno"
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"CreateMatch",
					mock.Anything,
					mock.Anything,
				).Return(errors.New("database error"))
			},
			expectedCode: 500,
		},
		{
			name: "empty body",
			body: `{}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"CreateMatch",
					mock.Anything,
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name: "unknown field ignored",
			body: `{
				"match_date":"2026-06-19T19:00:00Z",
				"title":"Final Liga",
				"home_team_id":1,
				"away_team_id":2,
				"home_score":0,
				"away_score":0,
				"phase":1,
				"venue_id":1,
				"venue_name":"Gelora Bung Karno",
				"stadium":"GBK"
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"CreateMatch",
					mock.Anything,
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name: "invalid match_date type",
			body: `{
				"match_date":123,
				"title":"Final Liga"
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "invalid home_team_id type",
			body: `{
				"match_date":"2026-06-19T19:00:00Z",
				"title":"Final Liga",
				"home_team_id":"1"
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "invalid away_team_id type",
			body: `{
				"match_date":"2026-06-19T19:00:00Z",
				"title":"Final Liga",
				"away_team_id":"2"
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "invalid home_score type",
			body: `{
				"home_score":"0"
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "invalid away_score type",
			body: `{
				"away_score":"0"
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "invalid phase type",
			body: `{
				"phase":"1"
			}`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "null values",
			body: `{
				"match_date":null,
				"home_team_id":null,
				"away_team_id":null,
				"phase":null
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"CreateMatch",
					mock.Anything,
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name: "negative score",
			body: `{
				"match_date":"2026-06-19T19:00:00Z",
				"title":"Final Liga",
				"home_team_id":1,
				"away_team_id":2,
				"home_score":-1,
				"away_score":0,
				"phase":1,
				"venue_id":1,
				"venue_name":"Gelora Bung Karno"
			}`,
			mockFn: func(s *handler.MockService) {
				s.On(
					"CreateMatch",
					mock.Anything,
					mock.Anything,
				).Return(nil)
			},
			expectedCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := handler.NewMockService(t)
			tt.mockFn(svc)

			h := handler.New(svc)

			r := gin.New()
			r.POST("/", h.CreateMatch)

			req := httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(tt.body),
			)

			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(
				t,
				tt.expectedCode,
				rec.Code,
			)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_FindMatchHighlight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	matchHighlight := model.MatchHighlight{
		MatchID: 1,
		Phase:   1,
		Scored: map[string]int{
			"home": 2,
			"away": 1,
		},
		Goal: []model.GoalHighlight{
			{
				PlayerID:       10,
				PlayerName:     "Ronaldo",
				ScoredAtMinute: 90,
			},
		},
	}

	tests := []struct {
		name string

		id string

		mockFn func(*handler.MockService)

		expectedCode int
	}{
		{
			name: "success",
			id:   "1",
			mockFn: func(s *handler.MockService) {
				s.On(
					"FindMatchHighlight",
					mock.Anything,
					int64(1),
				).Return(matchHighlight, nil)
			},
			expectedCode: 200,
		},
		{
			name: "service error",
			id:   "1",
			mockFn: func(s *handler.MockService) {
				s.On(
					"FindMatchHighlight",
					mock.Anything,
					int64(1),
				).Return(model.MatchHighlight{}, errors.New("match not found"))
			},
			expectedCode: 500,
		},
		{
			name:         "invalid id",
			id:           "abc",
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name:         "zero id",
			id:           "0",
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name:         "negative id",
			id:           "-1",
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := handler.NewMockService(t)
			tt.mockFn(svc)

			h := handler.New(svc)

			r := gin.New()
			r.GET("/:id/highlight", h.FindMatchHighlight)

			req := httptest.NewRequest(
				http.MethodGet,
				"/"+tt.id+"/highlight",
				nil,
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(
				t,
				tt.expectedCode,
				rec.Code,
			)

			svc.AssertExpectations(t)
		})
	}
}
