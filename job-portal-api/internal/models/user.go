package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name         string    `json:"name"`
	Email        string    `json:"email" gorm:"unique"`
	PasswordHash string    `json:"-"`
	DateOfBirth  time.Time `json:"date_of_birth"`
}

type NewUser struct {
	Name        string    `json:"name" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
}
type ResetRequest struct {
	Email       string    `json:"email" validate:"required,email"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
}
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Otp         string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"` // Add any password validation as needed
}
