package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azmanabdlh/ayo-example/internal/authentication/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		body any

		mockFn func(*handler.MockService)

		expectedCode int
		expectedBody map[string]any
	}{
		{
			name: "success",
			body: map[string]any{
				"email":    "admin@test.com",
				"password": "secret",
			},
			mockFn: func(s *handler.MockService) {
				s.On(
					"Login",
					mock.Anything,
					"admin@test.com",
					"secret",
				).Return(
					"jwt-token",
					nil,
				).Once()
			},
			expectedCode: 200,
			expectedBody: map[string]any{
				"code": float64(200),
				"data": "jwt-token",
			},
		},
		{
			name: "service error",
			body: map[string]any{
				"email":    "admin@test.com",
				"password": "secret",
			},
			mockFn: func(s *handler.MockService) {
				s.On(
					"Login",
					mock.Anything,
					"admin@test.com",
					"secret",
				).Return(
					"",
					errors.New("database error"),
				).Once()
			},
			expectedCode: 500,
			expectedBody: map[string]any{
				"code":    float64(500),
				"message": "database error",
			},
		},
		{
			name: "invalid json",
			body: "{invalid json",
			mockFn: func(s *handler.MockService) {

			},
			expectedCode: 400,
		},
		{
			name: "empty body",
			body: map[string]any{},
			mockFn: func(s *handler.MockService) {

			},
			expectedCode: 400,
			expectedBody: map[string]any{
				"code":    float64(400),
				"message": "Key: 'Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
		{
			name: "password number type",
			body: map[string]any{
				"email":    "admin@test.com",
				"password": 123,
			},
			mockFn: func(s *handler.MockService) {

			},
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSvc := handler.NewMockService(t)
			tt.mockFn(mockSvc)

			h := handler.New(mockSvc)

			r := gin.New()
			r.POST("/login", h.Login)

			var body []byte

			switch v := tt.body.(type) {
			case string:
				body = []byte(v)

			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewBuffer(body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)

			if tt.expectedBody != nil {
				var got map[string]any

				err := json.Unmarshal(rec.Body.Bytes(), &got)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, got)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}
