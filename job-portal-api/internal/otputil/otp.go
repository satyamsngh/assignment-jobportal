package otputil

import (
	"errors"
)

type Otp struct {
	Rd int
}

type UserOtp interface {
	GenerateOtp(email string) (string, error)
}

func NewOtp(rd int) (UserOtp, error) {
	if rd == 0 {
		return nil, errors.New("db cannot be null")
	}
	return Otp{
		Rd: rd,
	}, nil
}
