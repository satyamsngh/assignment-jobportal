package services

import (
	"context"
	"errors"
	"job-portal-api/internal/cache"
	"job-portal-api/internal/models"
	"job-portal-api/internal/otputil"
	"job-portal-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source services.go -destination mock_service.go -package services

type Service interface {
	CreatCompanies(ctx context.Context, nc models.NewComapanies, UserId uint) (models.Companies, error)
	ViewCompanies(ctx context.Context, companyId string) ([]models.Companies, error)
	ViewCompaniesById(ctx context.Context, companybyid uint, userId string) ([]models.Companies, error)
	CreateUser(ctx context.Context, nu models.NewUser) (models.User, error)
	CreateJob(ctx context.Context, newJob models.NewJob, userId string) (models.Job, error)
	AllJob(ctx context.Context, userId string) ([]models.Job, error)
	ListJobs(ctx context.Context, companyId uint, userId string) ([]models.Job, error)
	Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims,
		error)
	JobsByID(ctx context.Context, jobID uint64, userId string) (models.Job, error)
	CriteriaMeets(ctx context.Context, applicant []models.Application) ([]models.Application, error)
	VerifyOtpService(details models.ResetPasswordRequest) (bool, error)
	OtpService(details models.ResetRequest) (string, error)
}

type Store struct {
	UserRepo  repository.UserRepo
	UserCache cache.UserCache
	UserOtp   otputil.UserOtp
}

func NewStore(userRepo repository.UserRepo, userCache cache.UserCache, userOtp otputil.UserOtp) (Service, error) {
	if userRepo == nil {
		return nil, errors.New("interface cannot be null")
	}
	if userCache == nil {
		return nil, errors.New("cache interface cannot be null")
	}

	if userOtp == nil {
		return nil, errors.New("otp interface cannot be null")
	}
	return &Store{
		UserRepo:  userRepo,
		UserCache: userCache,
		UserOtp:   userOtp,
	}, nil
}
