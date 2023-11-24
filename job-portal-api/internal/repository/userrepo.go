package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"job-portal-api/internal/models"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func (r *Repo) CreateUser(ctx context.Context, UserDetails models.User) (models.User, error) {
	result := r.DB.Create(&UserDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.User{}, errors.New("could not create the user")
	}
	return UserDetails, nil
}
func (r *Repo) CheckEmail(ctx context.Context, email string, password string) (jwt.RegisteredClaims, error) {
	var u models.User
	tx := r.DB.Where("email = ?", email).First(&u)
	if tx.Error != nil {
		return jwt.RegisteredClaims{}, tx.Error
	}

	// We check if the provided password matches the hashed password in the database.
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}
	c := jwt.RegisteredClaims{
		Issuer:    "jobportal project",
		Subject:   strconv.FormatUint(uint64(u.ID), 10),
		Audience:  jwt.ClaimStrings{"companies"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return c, nil

}
func (r *Repo) IsUserPresentByEmail(email string) (bool, error) {
	var count int64

	result := r.DB.Where("email = ?", email).Model(&models.User{}).Count(&count)
	fmt.Print("==================---------", result)
	if result.Error != nil {
		return false, result.Error
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r *Repo) ResetPasswordByEmail(email, newPassword string) error {
	// Fetch the user by email from the database
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Error().AnErr("email id is not present signup firt", err)
		return err
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().AnErr("failed to hash password: %v", err)
		return err
	}

	// Update the user's password hash in the database
	if err := r.DB.Model(&user).Update("password_hash", string(hashedPassword)).Error; err != nil {
		log.Error().AnErr("failed to update password: %v", err)
		return err
	}

	return nil
}
