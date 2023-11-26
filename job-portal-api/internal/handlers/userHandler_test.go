package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	middlewares "job-portal-api/internal/middleware"
	"job-portal-api/internal/models"
	"job-portal-api/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_RequestPasswordReset(t *testing.T) {
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "missing trace id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{
			name: "invalid request body",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"invalid`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{
			name: "user not found",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"email": "user@example.com", "date_of_birth": "1990-01-01"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().OtpService(gomock.Any()).Return(errors.New("")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"msg":"signup first"}`,
		},
		{
			name: "error generating OTP",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"email": "user@example.com", "date_of_birth": "1990-01-01"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().OtpService(gomock.Any()).Return(errors.New("")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"msg":"signup first"}`,
		},
		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"email": "user@example.com", "date_of_birth": "1990-01-01"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().OtpService(gomock.Any()).Return(nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"msg":"Successfully changed the password"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				s: ms,
			}
			h.RequestPasswordReset(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_ResetPassword(t *testing.T) {
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "missing trace id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{
			name: "invalid request body",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"invalid`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: 500,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{
			name: "error",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				resetRequest := models.ResetPasswordRequest{
					Email:           "satyam",
					Otp:             "326273",
					NewPassword:     "satyam",
					ConfirmPassword: "satyam",
				}

				// Convert ResetRequest to JSON
				requestBody, _ := json.Marshal(resetRequest)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBuffer(requestBody))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().VerifyOtpService(gomock.Any()).Return(errors.New("")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"unable to verify otp"}`,
		},
		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, *services.MockService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				resetRequest := models.ResetPasswordRequest{
					Email:           "satyam",
					Otp:             "326273",
					NewPassword:     "satyam",
					ConfirmPassword: "satyam",
				}

				// Convert ResetRequest to JSON
				requestBody, _ := json.Marshal(resetRequest)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBuffer(requestBody))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().VerifyOtpService(gomock.Any()).Return(nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"msg":"Successfully changed the password"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				s: ms,
			}
			h.ResetPassword(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
