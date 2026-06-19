package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/azmanabdlh/ayo-example/internal/team-management/handler"
	"github.com/azmanabdlh/ayo-example/internal/team-management/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_FindAllTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		query string

		mockFn func(*handler.MockService)

		expectedCode  int
		expectedPage  int
		expectedLimit int
	}{
		{
			name:  "default pagination",
			query: "",
			mockFn: func(s *handler.MockService) {
				s.On(
					"FindAll",
					mock.Anything,
					0,
					10,
				).Return([]model.Team{{ID: 1, Name: "A"}}, nil)
			},
			expectedCode:  200,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:  "page2_limit5",
			query: "?page=2&limit=5",
			mockFn: func(s *handler.MockService) {
				s.On(
					"FindAll",
					mock.Anything,
					5,
					5,
				).Return([]model.Team{{ID: 2, Name: "B"}}, nil)
			},
			expectedCode:  200,
			expectedPage:  2,
			expectedLimit: 5,
		},
		{
			name:  "invalid_limit_becomes_zero",
			query: "?limit=abc",
			mockFn: func(s *handler.MockService) {
				s.On(
					"FindAll",
					mock.Anything,
					0,
					0,
				).Return([]model.Team{}, nil)
			},
			expectedCode:  200,
			expectedPage:  1,
			expectedLimit: 0,
		},
		{
			name:  "limit_too_large_clipped",
			query: "?limit=1000",
			mockFn: func(s *handler.MockService) {
				s.On(
					"FindAll",
					mock.Anything,
					0,
					100,
				).Return([]model.Team{}, nil)
			},
			expectedCode:  200,
			expectedPage:  1,
			expectedLimit: 100,
		},
		{
			name:  "service_error",
			query: "?page=1&limit=10",
			mockFn: func(s *handler.MockService) {
				s.On(
					"FindAll",
					mock.Anything,
					0,
					10,
				).Return(nil, errors.New("db failure"))
			},
			expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			svc := handler.NewMockService(t)
			tt.mockFn(svc)

			h := handler.New(svc)

			r := gin.New()
			r.GET("/", h.FindAllTeam)

			req := httptest.NewRequest(
				http.MethodGet,
				"/"+tt.query,
				nil,
			)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedCode == 200 {
				var resp map[string]interface{}
				_ = json.Unmarshal(rec.Body.Bytes(), &resp)
				meta := resp["meta"].(map[string]interface{})
				assert.Equal(t, float64(tt.expectedPage), meta["page"].(float64))
				assert.Equal(t, float64(tt.expectedLimit), meta["limit"].(float64))
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_FindTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		id string

		mockFn func(*handler.MockService)

		expectedCode int
	}{
		{
			name: "success",
			id:   "123",
			mockFn: func(s *handler.MockService) {
				s.On("Find", mock.Anything, int64(123)).Return(model.Team{ID: 123, Name: "T"}, nil)
			},
			expectedCode: 200,
		},
		{
			name: "service_error",
			id:   "123",
			mockFn: func(s *handler.MockService) {
				s.On("Find", mock.Anything, int64(123)).Return(model.Team{}, errors.New("not found"))
			},
			expectedCode: 500,
		},
		{
			name: "invalid id becomes zero",
			id:   "abc",
			mockFn: func(s *handler.MockService) {
				s.On("Find", mock.Anything, int64(0)).Return(model.Team{ID: 0, Name: "Zero"}, nil)
			},
			expectedCode: 200,
		},
		{
			name: "negative id",
			id:   "-5",
			mockFn: func(s *handler.MockService) {
				s.On("Find", mock.Anything, int64(-5)).Return(model.Team{}, errors.New("invalid"))
			},
			expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := handler.NewMockService(t)
			tt.mockFn(svc)

			h := handler.New(svc)

			r := gin.New()
			r.GET("/:id", h.FindTeam)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_ModifyTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		body         string
		mockFn       func(*handler.MockService)
		expectedCode int
	}{
		{
			name: "success",
			id:   "10",
			body: `{"name":"New"}`,
			mockFn: func(s *handler.MockService) {
				s.On("Modify", mock.Anything, int64(10), mock.Anything).Return(nil)
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
			body: `{"name":"X"}`,
			mockFn: func(s *handler.MockService) {
				s.On("Modify", mock.Anything, int64(10), mock.Anything).Return(errors.New("modify fail"))
			},
			expectedCode: 500,
		},
		{
			name: "invalid id becomes zero",
			id:   "abc",
			body: `{"name":"Zero"}`,
			mockFn: func(s *handler.MockService) {
				s.On("Modify", mock.Anything, int64(0), mock.Anything).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name:         "invalid field type",
			id:           "10",
			body:         `{"founded_year":"notint"}`,
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
			r.PUT("/:id", h.ModifyTeam)

			req := httptest.NewRequest(http.MethodPut, "/"+tt.id, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_CreateTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		body         string
		mockFn       func(*handler.MockService)
		expectedCode int
	}{
		{
			name: "success",
			id:   "0",
			body: `{"name":"New"}`,
			mockFn: func(s *handler.MockService) {
				s.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name:         "invalid json",
			id:           "0",
			body:         `{`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "service error",
			id:   "0",
			body: `{"name":"X"}`,
			mockFn: func(s *handler.MockService) {
				s.On("Create", mock.Anything, mock.Anything).Return(errors.New("create fail"))
			},
			expectedCode: 500,
		},
		{
			name:         "invalid field type",
			id:           "0",
			body:         `{"founded_year":"bad"}`,
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
			r.POST("/:id", h.CreateTeam)

			req := httptest.NewRequest(http.MethodPost, "/"+tt.id, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_RemoveTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		mockFn       func(*handler.MockService)
		expectedCode int
	}{
		{
			name: "success",
			id:   "10",
			mockFn: func(s *handler.MockService) {
				s.On("Remove", mock.Anything, int64(10)).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name: "service error",
			id:   "10",
			mockFn: func(s *handler.MockService) {
				s.On("Remove", mock.Anything, int64(10)).Return(errors.New("remove fail"))
			},
			expectedCode: 500,
		},
		{
			name: "invalid id becomes zero",
			id:   "abc",
			mockFn: func(s *handler.MockService) {
				s.On("Remove", mock.Anything, int64(0)).Return(nil)
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
			r.DELETE("/:id", h.RemoveTeam)

			req := httptest.NewRequest(http.MethodDelete, "/"+tt.id, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_AssignPlayerToTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         string
		mockFn       func(*handler.MockService)
		expectedCode int
	}{
		{
			name: "success",
			body: `{"player_id":1,"team_id":2}`,
			mockFn: func(s *handler.MockService) {
				s.On("AssignPlayerTeam", mock.Anything, int64(1), int64(2)).Return(nil)
			},
			expectedCode: 201,
		},
		{
			name:         "invalid json",
			body:         `{`,
			mockFn:       func(s *handler.MockService) {},
			expectedCode: 400,
		},
		{
			name: "service error",
			body: `{"player_id":1,"team_id":2}`,
			mockFn: func(s *handler.MockService) {
				s.On("AssignPlayerTeam", mock.Anything, int64(1), int64(2)).Return(errors.New("assign fail"))
			},
			expectedCode: 500,
		},
		{
			name:         "invalid types",
			body:         `{"player_id":"abc","team_id":2}`,
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
			r.POST("/assign", h.AssignPlayerToTeam)

			req := httptest.NewRequest(http.MethodPost, "/assign", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedCode, rec.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_FindPlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		mockFn       func(*handler.MockService)
		expectedCode int
	}{
		{
			name: "success",
			id:   "5",
			mockFn: func(s *handler.MockService) {
				s.On("FindPlayer", mock.Anything, int64(5)).Return(model.Player{ID: 5, Name: "P"}, nil)
			},
			expectedCode: 200,
		},
		{
			name: "service error",
			id:   "5",
			mockFn: func(s *handler.MockService) {
				s.On("FindPlayer", mock.Anything, int64(5)).Return(model.Player{}, errors.New("not found"))
			},
			expectedCode: 500,
		},
		{
			name: "invalid id becomes zero",
			id:   "abc",
			mockFn: func(s *handler.MockService) {
				s.On("FindPlayer", mock.Anything, int64(0)).Return(model.Player{}, nil)
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
			r.GET("/:id", h.FindPlayer)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			svc.AssertExpectations(t)
		})
	}
}

func TestHandler_RemovePlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		id           string
		mockFn       func(*handler.MockService)
		expectedCode int
	}{
		{
			name: "success",
			id:   "7",
			mockFn: func(s *handler.MockService) {
				s.On("RemovePlayer", mock.Anything, int64(7)).Return(nil)
			},
			expectedCode: 200,
		},
		{
			name: "service error",
			id:   "7",
			mockFn: func(s *handler.MockService) {
				s.On("RemovePlayer", mock.Anything, int64(7)).Return(errors.New("remove fail"))
			},
			expectedCode: 500,
		},
		{
			name: "invalid id becomes zero",
			id:   "abc",
			mockFn: func(s *handler.MockService) {
				s.On("RemovePlayer", mock.Anything, int64(0)).Return(nil)
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
			r.DELETE("/:id", h.RemovePlayer)

			req := httptest.NewRequest(http.MethodDelete, "/"+tt.id, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			svc.AssertExpectations(t)
		})
	}
}
