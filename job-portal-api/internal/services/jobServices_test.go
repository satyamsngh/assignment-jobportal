package services

import (
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"
	"reflect"
	"testing"
)

func TestStore_CreateJob(t *testing.T) {
	type args struct {
		ctx    context.Context
		jobs   models.NewJob
		userID string
	}
	tests := []struct {
		name        string
		args        args
		want        models.Job
		wantErr     bool
		mockNewRepo func() (models.Job, error)
	}{
		{
			name: "Error",
			args: args{
				ctx:  context.Background(),
				jobs: models.NewJob{
					// Fill in the fields of the NewJob struct as needed for this test case
				},
				userID: "1",
			},
			want:    models.Job{}, // Define the expected result for this case
			wantErr: true,
			mockNewRepo: func() (models.Job, error) {
				return models.Job{}, errors.New("database error")
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				jobs: models.NewJob{
					MaxExp: 1,
					Shift:  "Day",
				},
				userID: "1",
			},
			want: models.Job{
				MaxExp: 1,
				Shift:  "Day",
			},
			wantErr: false,
			mockNewRepo: func() (models.Job, error) {
				return models.Job{
					MaxExp: 1,
					Shift:  "Day",
				}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mock)
			if tt.mockNewRepo != nil {
				mockRepo.EXPECT().CreateJob(tt.args.ctx, tt.args.jobs).Return(tt.mockNewRepo()).AnyTimes()
			}
			s := &Store{
				UserRepo: mockRepo,
			}

			got, err := s.CreateJob(tt.args.ctx, tt.args.jobs, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateJob() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_CriteriaMeets(t *testing.T) {
	type args struct {
		ctx       context.Context
		applicant []models.Application
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Application
		wantErr bool
		setup   func(mockRepo *repository.MockUserRepo)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				applicant: []models.Application{
					{
						JobID:            1,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     15,
						Budget:           400,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1, 2},
						Exp:              8,
						QualificationIDs: []uint{},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            2,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     15,
						Budget:           400,
						LocationIDs:      []uint{},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1, 2},
						Exp:              8,
						QualificationIDs: []uint{1},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            3,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     15,
						Budget:           4000000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1, 2},
						Exp:              8,
						QualificationIDs: []uint{1, 3},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            4,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     15,
						Budget:           4000000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1, 2},
						Exp:              8,
						QualificationIDs: []uint{1, 3},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            5,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     15,
						Budget:           400,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1, 2},
						Exp:              1566,
						QualificationIDs: []uint{1, 3},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            6,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     7,
						Budget:           40000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2, 3},
						WorkModeIDs:      []uint{1, 2},
						Exp:              1,
						QualificationIDs: []uint{1, 3},
						Shift:            "D",
						JobType:          "Remote",
					},
					{
						JobID:            7,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     7,
						Budget:           40000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2, 3},
						WorkModeIDs:      []uint{1, 2},
						Exp:              1,
						QualificationIDs: []uint{1, 3},
						Shift:            "Day",
						JobType:          "Re",
					},
					{
						JobID:            8,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     7,
						Budget:           40000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2, 3},
						WorkModeIDs:      []uint{},
						Exp:              1,
						QualificationIDs: []uint{1, 3},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            9,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     7,
						Budget:           40000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{},
						WorkModeIDs:      []uint{1, 2},
						Exp:              1,
						QualificationIDs: []uint{1, 3},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            10,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     1,
						Budget:           40000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1, 2},
						Exp:              1,
						QualificationIDs: []uint{1, 3},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            11,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     3333,
						Budget:           40000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1, 2},
						Exp:              1,
						QualificationIDs: []uint{1, 3},
						Shift:            "Day",
						JobType:          "Remote",
					},
				},
			},
			want:    []models.Application{},
			wantErr: false,
			setup: func(mockRepo *repository.MockUserRepo) {
				//mockRepo.EXPECT(3).GetJobById(gomock.Any(), uint(0)).Return(models.Job{}, errors.New("test error")).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(1)).Return(models.Job{
					Model:           gorm.Model{ID: 1},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(2)).Return(models.Job{
					Model:           gorm.Model{ID: 2},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(3)).Return(models.Job{
					Model:           gorm.Model{ID: 2},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(4)).Return(models.Job{
					Model:           gorm.Model{ID: 4},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(5)).Return(models.Job{
					Model:           gorm.Model{ID: 5},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(6)).Return(models.Job{
					Model:           gorm.Model{ID: 6},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(7)).Return(models.Job{
					Model:           gorm.Model{ID: 7},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(8)).Return(models.Job{
					Model:           gorm.Model{ID: 8},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(9)).Return(models.Job{
					Model:           gorm.Model{ID: 9},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(10)).Return(models.Job{
					Model:           gorm.Model{ID: 10},
					CompanyID:       1,
					MinNoticePeriod: 4,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(11)).Return(models.Job{
					Model:           gorm.Model{ID: 11},
					CompanyID:       1,
					MinNoticePeriod: 4,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				applicant: []models.Application{
					{
						JobID:            0,
						Name:             "John",
						Email:            "john@gmail.com",
						Phone:            "1234567890",
						Resume:           "",
						NoticePeriod:     15,
						Budget:           400000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1},
						Exp:              2,
						QualificationIDs: []uint{1, 2},
						Shift:            "Day",
						JobType:          "Remote",
					},
					{
						JobID:            1,
						Name:             "Monika",
						Email:            "monika@gmail.com",
						Phone:            "3434343434",
						Resume:           "",
						NoticePeriod:     15,
						Budget:           400000,
						LocationIDs:      []uint{1},
						TechnologyIDs:    []uint{1, 2},
						WorkModeIDs:      []uint{1},
						Exp:              2,
						QualificationIDs: []uint{1, 2},
						Shift:            "Day",
						JobType:          "Remote",
					},
				},
			},
			want: []models.Application{
				{
					JobID:            1,
					Name:             "Monika",
					Email:            "monika@gmail.com",
					Phone:            "3434343434",
					Resume:           "",
					NoticePeriod:     15,
					Budget:           400000,
					LocationIDs:      []uint{1},
					TechnologyIDs:    []uint{1, 2},
					WorkModeIDs:      []uint{1},
					Exp:              2,
					QualificationIDs: []uint{1, 2},
					Shift:            "Day",
					JobType:          "Remote",
				},
			},
			wantErr: false,
			setup: func(mockRepo *repository.MockUserRepo) {
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(0)).Return(models.Job{}, errors.New("test error")).Times(1)
				mockRepo.EXPECT().GetJobById(gomock.Any(), uint(1)).Return(models.Job{
					Model:           gorm.Model{ID: 1},
					CompanyID:       1,
					MinNoticePeriod: 0,
					MaxNoticePeriod: 60,
					Budget:          600000,
					JobLocations: []models.JobLocation{
						{Model: gorm.Model{ID: 1}},
					},
					Technology: []models.Technologies{
						{Model: gorm.Model{ID: 1}},
					},
					WorkMode: []models.WorkModes{
						{Model: gorm.Model{ID: 1}},
					},
					MaxExp: 3,
					Qualification: []models.Qualifications{
						{Model: gorm.Model{ID: 1}},
					},
					Shift:   "Day",
					JobType: "Remote",
				}, nil).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			tt.setup(mockRepo)
			s := &Store{
				UserRepo: mockRepo,
			}
			got, err := s.CriteriaMeets(tt.args.ctx, tt.args.applicant)
			if (err != nil) != tt.wantErr {
				t.Errorf("CriteriaMeets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CriteriaMeets() got = %v, want %v", got, tt.want)
			}
		})
	}
}
