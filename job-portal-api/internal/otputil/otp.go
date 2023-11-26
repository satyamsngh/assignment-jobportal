package otputil

import (
	"errors"
)

type Otp struct {
	Rd int
}

//go:generate mockgen -source otp.go -destination mock_otputil.go -package otputil
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
