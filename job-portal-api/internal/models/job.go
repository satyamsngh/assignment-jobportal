package models

import (
	"gorm.io/gorm"
)

type Companies struct {
	gorm.Model
	CompanyName string `json:"company_name"`
	FoundedYear int    `json:"founded_year"`
	Location    string `json:"location"`
	UserId      uint   `json:"user_id"`
	Address     string `json:"address"`
	Jobs        []Job  `json:"jobs,omitempty" gorm:"foreignKey:CompanyID"`
}

type NewComapanies struct {
	CompanyName string `json:"company_name" validate:"required"`
	FoundedYear int    `json:"founded_year" validate:"required,number"`
	Location    string `json:"location" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Jobs        []Job  `json:"jobs"`
}

type Job struct {
	gorm.Model
	CompanyID       uint             `json:"company_id"`
	MinNoticePeriod string           `json:"min_notice_period"`
	MaxNoticePeriod string           `json:"max_notice_period"`
	Budget          float64          `json:"budget"`
	JobLocations    []JobLocation    `gorm:"many2many:job_location_relations;" json:"job_locations"`
	Technology      []Technologies   `gorm:"many2many:job_technologies_relations;" json:"job_technologies"`
	WorkMode        []WorkModes      `gorm:"many2many:job_workmodes_relations;" json:"job_workmodes"`
	MaxExp          int              `json:"max_exp"`
	Qualification   []Qualifications `gorm:"many2many:job_qualification_relations;" json:"job_qualification"`
	Shift           string           `json:"shift"`
	JobType         string           `json:"job_type"`
}
type JobLocation struct {
	gorm.Model

	Location string `json:"location"`
}
type Technologies struct {
	gorm.Model
	Technology string `json:"technology_1"`
}

type Qualifications struct {
	gorm.Model
	Qualifications string `json:"qualifications_1"`
}

type WorkModes struct {
	gorm.Model

	WorkMode1 string `json:"work_mode_1"`
}

type Application struct {
	JobID           uint             `json:"job_id"`
	Name            string           `json:"name"`
	Email           string           `json:"email"`
	Phone           string           `json:"phone"`
	Resume          string           `json:"resume"`
	MinNoticePeriod string           `json:"min_notice_period"`
	MaxNoticePeriod string           `json:"max_notice_period"`
	Budget          float64          `json:"budget"`
	JobLocations    []JobLocation    `json:"job_locations"`
	Technology      []Technologies   `json:"technology"`
	WorkMode        []WorkModes      `json:"work_mode"`
	Exp             int              `json:"exp"`
	Qualification   []Qualifications `json:"qualification"`
	Shift           string           `json:"shift"`
	JobType         string           `json:"job_type"`
}
