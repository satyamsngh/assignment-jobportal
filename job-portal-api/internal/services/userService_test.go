package services

import (
	"errors"
	"go.uber.org/mock/gomock"
	"job-portal-api/internal/cache"
	"job-portal-api/internal/models"
	"job-portal-api/internal/otputil"
	"job-portal-api/internal/repository"
	"testing"
)

func TestStore_VerifyOtpService(t *testing.T) {
	type args struct {
		details models.ResetPasswordRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp)
	}{
		{
			name: "success",
			args: args{
				models.ResetPasswordRequest{
					Email:           "satyam18577@gmail.com",
					Otp:             "123456",
					NewPassword:     "satyam",
					ConfirmPassword: "satyam",
				},
			},
			wantErr: false,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockcache.EXPECT().GetRedisKeyOtp(gomock.Any()).Return("123456", nil)
				mockrepo.EXPECT().ResetPasswordByEmail("satyam18577@gmail.com", "satyam").Return(nil)
				mockcache.EXPECT().DelRedisKey("satyam18577@gmail.com").Return(nil)
			},
		},
		{
			name: "error",
			args: args{
				models.ResetPasswordRequest{
					Email:           "xyz@gmail.com",
					Otp:             "123456",
					NewPassword:     "satyam",
					ConfirmPassword: "satyam",
				},
			},
			wantErr: true,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockcache.EXPECT().GetRedisKeyOtp(gomock.Any()).Return("", errors.New(""))
			},
		},
		{
			name: "error",
			args: args{
				models.ResetPasswordRequest{
					Email:           "xyz@gmail.com",
					Otp:             "123456",
					NewPassword:     "satyam",
					ConfirmPassword: "satyam",
				},
			},
			wantErr: true,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockcache.EXPECT().GetRedisKeyOtp(gomock.Any()).Return("", errors.New(""))
			},
		},
		{
			name: "error",
			args: args{
				models.ResetPasswordRequest{
					Email:           "xyz@gmail.com",
					Otp:             "123456",
					NewPassword:     "satyam",
					ConfirmPassword: "satyam",
				},
			},
			wantErr: true,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockcache.EXPECT().GetRedisKeyOtp(gomock.Any()).Return("123", nil)
			},
		},
		{
			name: "error",
			args: args{
				models.ResetPasswordRequest{
					Email:           "xyz@gmail.com",
					Otp:             "123456",
					NewPassword:     "satyam",
					ConfirmPassword: "saty",
				},
			},
			wantErr: true,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockcache.EXPECT().GetRedisKeyOtp(gomock.Any()).Return("123456", nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			repomock := repository.NewMockUserRepo(mc)
			cachemock := cache.NewMockUserCache(mc)
			otputilmock := otputil.NewMockUserOtp(mc)
			tt.setup(repomock, cachemock, otputilmock)

			s := &Store{
				UserRepo:  repomock,
				UserCache: cachemock,
				UserOtp:   otputilmock,
			}
			err := s.VerifyOtpService(tt.args.details)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyOtpService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_OtpService(t *testing.T) {
	type args struct {
		details models.ResetRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp)
	}{
		{
			name: "success",
			args: args{
				models.ResetRequest{
					Email:       "xyz@gmail.com",
					DateOfBirth: "23-12-1999",
				},
			},
			wantErr: false,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockrepo.EXPECT().IsUserPresentByEmail(gomock.Any()).Return(true, nil).AnyTimes()
				mockotp.EXPECT().GenerateOtp(gomock.Any()).Return("123456", nil)
				mockcache.EXPECT().SetRedisKeyOtp(gomock.Any(), "123456").Return(nil)
			},
		},
		{
			name: "error",
			args: args{
				models.ResetRequest{
					Email:       "xyz@gmail.com",
					DateOfBirth: "23-12-1999",
				},
			},
			wantErr: true,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockrepo.EXPECT().IsUserPresentByEmail(gomock.Any()).Return(false, errors.New("")).AnyTimes()

			},
		},
		{
			name: "error",
			args: args{
				models.ResetRequest{
					Email:       "xyz@gmail.com",
					DateOfBirth: "23-12-1999",
				},
			},
			wantErr: true,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockrepo.EXPECT().IsUserPresentByEmail(gomock.Any()).Return(false, nil).AnyTimes()

			},
		},
		{
			name: "error",
			args: args{
				models.ResetRequest{
					Email:       "xyz@gmail.com",
					DateOfBirth: "23-12-1999",
				},
			},
			wantErr: true,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockrepo.EXPECT().IsUserPresentByEmail(gomock.Any()).Return(true, nil).AnyTimes()
				mockotp.EXPECT().GenerateOtp(gomock.Any()).Return("", errors.New(""))

			},
		},
		{
			name: "errors",
			args: args{
				models.ResetRequest{
					Email:       "xyz@gmail.com",
					DateOfBirth: "23-12-1999",
				},
			},
			wantErr: true,
			setup: func(mockrepo *repository.MockUserRepo, mockcache *cache.MockUserCache, mockotp *otputil.MockUserOtp) {
				mockrepo.EXPECT().IsUserPresentByEmail(gomock.Any()).Return(true, nil).AnyTimes()
				mockotp.EXPECT().GenerateOtp(gomock.Any()).Return("123456", nil)
				mockcache.EXPECT().SetRedisKeyOtp(gomock.Any(), "123456").Return(errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockrepo := repository.NewMockUserRepo(mc)
			mockcache := cache.NewMockUserCache(mc)
			mockotp := otputil.NewMockUserOtp(mc)
			tt.setup(mockrepo, mockcache, mockotp)

			s := &Store{
				UserRepo:  mockrepo,
				UserCache: mockcache,
				UserOtp:   mockotp,
			}
			if err := s.OtpService(tt.args.details); (err != nil) != tt.wantErr {
				t.Errorf("OtpService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
