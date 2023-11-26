package services

import (
	"context"
	"errors"
	"fmt"
	"job-portal-api/internal/models"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var otpTime time.Time

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
		DateOfBirth:  nu.DateOfBirth,
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

func (s *Store) OtpService(details models.ResetRequest) error {

	userPresent, err := s.UserRepo.IsUserPresentByEmail(details.Email)
	if err != nil {
		log.Printf("Error checking user presence: %v", err)
		return err
	}

	if !userPresent {
		// Log that the user is not present
		log.Printf("User not present for email %s", details.Email)
		errs := errors.New("user not present in the database")
		return errs
	}
	generatedOtp, err := s.UserOtp.GenerateOtp(details.Email)
	if err != nil {
		// Log the error
		log.Printf("Error generating OTP: %v", err)
		return err
	}
	err = s.UserCache.SetRedisKeyOtp(details.Email, generatedOtp)
	if err != nil {
		// Log the error
		log.Printf("Error setting OTP in cache: %v", err)
		return err
	}
	log.Printf("OTP sent successfully for email %s", details.Email)
	otpTime = time.Now()
	return nil
}

func (s *Store) VerifyOtpService(details models.ResetPasswordRequest) error {
	//nowTime := time.Now()
	//duration := nowTime.Sub(otpTime)
	//fmt.Println(duration)
	//if duration <= time.Minute {
	//	fmt.Println("otp is valid")
	//} else {
	//	err := s.UserCache.DelRedisKey(details.Email)
	//	if err != nil {
	//		return err
	//	}
	//	errs := errors.New("otp expired regenerate otp")
	//	return errs
	//}
	storedOtp, err := s.UserCache.GetRedisKeyOtp(details.Email)
	if err != nil {
		return err
	}
	if details.Otp != storedOtp {
		errs := errors.New("otp is incorrect")
		return errs
	}
	if len(details.NewPassword) != len(details.ConfirmPassword) {
		errs := errors.New("password not match,retry")
		return errs
	}
	err = s.UserRepo.ResetPasswordByEmail(details.Email, details.NewPassword)
	if err != nil {
		return err
	}
	err = s.UserCache.DelRedisKey(details.Email)
	if err != nil {
		return err
	}
	return nil
}
