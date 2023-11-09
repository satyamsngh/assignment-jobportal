package handlers

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"

	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go.uber.org/mock/gomock"
	"job-portal-api/internal/auth"
	middlewares "job-portal-api/internal/middleware"
	"job-portal-api/internal/models"
	"job-portal-api/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_CreateJob(t *testing.T) {
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, services.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "missing trace id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
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
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"invalid`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "123"})

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid request body"}`,
		},
		{
			name: "error while creating job posting",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "123"})
				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)

				ms.EXPECT().CreateJob(c.Request.Context(), gomock.Any(), gomock.Any()).Return(models.Job{}, errors.New("test service error")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Failed to create job"}`,
		},
		{
			name: "error compnay id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "1s"})
				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)

				ms.EXPECT().CreateJob(c.Request.Context(), gomock.Any(), gomock.Any()).Return(models.Job{}, errors.New("")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: 400,
			expectedResponse:   `{"error":"Invalid company ID"}`,
		},
		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "companyID", Value: "123"})
				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)

				ms.EXPECT().CreateJob(c.Request.Context(), gomock.Any(), gomock.Any()).Return(models.Job{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: 201,
			expectedResponse:   `{"Id":0}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				s: ms,
			}
			h.CreateJob(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_AddApplicants(t *testing.T) {
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, services.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "missing trace id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "invalid request body",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`{"invalid`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid request payload"}`,
		},
		{
			name: "error while processing applicants",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`[]`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().CriteriaMeets(ctx, gomock.Any()).Return(nil, errors.New("test error")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Applicant creation failed"}`,
		},
		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.Service) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodPost, "http://test.com:8080", bytes.NewBufferString(`[]`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := services.NewMockService(mc)
				ms.EXPECT().CriteriaMeets(ctx, gomock.Any()).Return([]models.Application{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()

			h := &handler{
				s: ms,
			}
			h.AddApplicants(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
