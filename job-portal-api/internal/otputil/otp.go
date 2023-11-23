package otputil

import (
	"errors"
)

type Otp struct {
	Rd *string
}

type UserOtp interface {
	GenerateOtp(email string) string
}

func NewOtp(rd *string) (UserOtp, error) {
	if rd == nil {
		return nil, errors.New("db cannot be null")
	}
	return &Otp{
		Rd: rd,
	}, nil
}
