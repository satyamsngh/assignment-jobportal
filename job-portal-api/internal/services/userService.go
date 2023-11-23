package services

import (
	"context"
	"errors"
	"fmt"
	"job-portal-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Store) CreateUser(ctx context.Context, nu models.NewUser) (models.User, error) {

	// We hash the user's password for storage in the database.
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("generating password hash: %w", err)
	}

	u := models.User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: string(hashedPass),
	}

	user, err := s.UserRepo.CreateUser(ctx, u)
	if err != nil {
		return models.User{}, err

	}
	return user, nil
}

func (s *Store) Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims,
	error) {

	claims, err := s.UserRepo.CheckEmail(ctx, email, password)
	if err != nil {
		return jwt.RegisteredClaims{}, errors.New("not able to generate claims")
	}
	return claims, nil
}

func (s *Store) OtpService(details models.ResetRequest) (string, error) {
	status, err := s.UserRepo.IsUserPresentByEmailAndDOB(details.Email, details.DateOfBirth)
	if err != nil {
		// Handle database error
		return "", err
	}
	fmt.Println(status)
	generatedOtp := s.UserOtp.GenerateOtp(details.Email)
	s.UserCache.SetRedisKeyOtp(details.Email, generatedOtp)
	return "otp send succesfully", nil

}
func (s *Store) VerifyOtpService(details models.ResetPasswordRequest) (bool, error) {
	// Retrieve stored OTP from Redis
	storedOtp, err := s.UserCache.GetRedisKeyOtp(details.Email)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve OTP from Redis: %v", err)
	}

	// Verify if the provided OTP matches the stored OTP
	if details.Otp != storedOtp {
		return false, nil // Invalid OTP
	}

	// Reset the password in the database (replace with your actual password reset logic)
	err = s.UserRepo.ResetPasswordByEmail(details.Email, details.NewPassword)
	if err != nil {
		return false, fmt.Errorf("failed to reset password: %v", err)
	}

	// Delete the OTP from Redis after successful password reset
	//err = s.RedisClient.Del(context.
	//if err != nil {
	//	return true, fmt.Errorf("failed to delete OTP from Redis: %v", err)
	//}

	return true, nil // Password reset successfully
}
