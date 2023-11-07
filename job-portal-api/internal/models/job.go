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
type NewJob struct {
	CompanyID       uint             `json:"company_id" validate:"required"`
	MinNoticePeriod string           `json:"min_notice_period" validate:"required"`
	MaxNoticePeriod string           `json:"max_notice_period" validate:"required"`
	Budget          float64          `json:"budget" validate:"required"`
	JobLocations    []JobLocation    `json:"job_locations" validate:"required"`
	Technology      []Technologies   `json:"technology" validate:"required"`
	WorkMode        []WorkModes      `json:"work_mode" validate:"required"`
	MaxExp          int              `json:"max_experience" validate:"required"`
	Qualification   []Qualifications `json:"qualification" validate:"required"`
	Shift           string           `json:"shift" validate:"required"`
	JobType         string           `json:"job_type" validate:"required"`
}

type JobLocation struct {
	gorm.Model

	Location string `json:"location" gorm:"unique"`
}

type Technologies struct {
	gorm.Model
	Technology string `json:"technology" gorm:"unique"`
}

type Qualifications struct {
	gorm.Model
	Qualification string `json:"qualification" gorm:"unique"`
}

type WorkModes struct {
	gorm.Model

	WorkMode1 string `json:"work_mode" gorm:"unique"`
}

type Application struct {
	JobID         uint             `json:"job_id"`
	Name          string           `json:"name"`
	Email         string           `json:"email"`
	Phone         string           `json:"phone"`
	Resume        string           `json:"resume"`
	NoticePeriod  int              `json:"notice_period"`
	Budget        float64          `json:"budget"`
	JobLocations  []JobLocation    `json:"job_locations"`
	Technology    []Technologies   `json:"technology"`
	WorkMode      []WorkModes      `json:"work_mode"`
	Exp           int              `json:"exp"`
	Qualification []Qualifications `json:"qualification"`
	Shift         string           `json:"shift"`
	JobType       string           `json:"job_type"`
}

//{
//"job_id": 1,
//"name": "John Doe",
//"email": "john.doe@example.com",
//"phone": "1234567890",
//"resume": "path_to_resume.pdf",
//"min_notice_period": "1 month",
//"max_notice_period": "2 months",
//"budget": 50000.00,
//"job_locations": [
//{
//"id":1,
//"location": "Delhi"
//}
//],
//"technology": [
//{
//"id":1,
//"technology": "Java"
//}
//],
//"work_mode": [
//{
//"id":1,
//"work_mode": "Full time"
//}
//],
//"exp": 5,
//"qualification": [
//{
//"id":1,
//"qualification": "BE"
//}
//],
//"shift": "Day",
//"job_type": "Software Engineer"
//}
