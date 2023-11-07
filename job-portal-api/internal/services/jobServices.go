package services

import (
	"context"
	"job-portal-api/internal/models"
	"strconv"
	"sync"
)

func (s *Store) CreatCompanies(ctx context.Context, nc models.NewComapanies, UserID uint) (models.Companies, error) {

	com := models.Companies{
		CompanyName: nc.CompanyName,
		FoundedYear: nc.FoundedYear,
		Location:    nc.Location,
		UserId:      UserID,
		Address:     nc.Address,
		Jobs:        nc.Jobs,
	}

	com, err := s.UserRepo.CreateCompany(ctx, com)
	if err != nil {
		return models.Companies{}, err
	}
	return com, nil
}

func (s *Store) ViewCompanies(ctx context.Context, companyID string) ([]models.Companies, error) {
	companies, err := s.UserRepo.ViewCompanies(ctx)
	if err != nil {
		return nil, err
	}
	return companies, nil

}

func (s *Store) ViewCompaniesById(ctx context.Context, companyID uint, userID string) ([]models.Companies, error) {
	company, err := s.UserRepo.ViewCompanyById(ctx, companyID)
	if err != nil {
		return []models.Companies{}, err
	}

	return company, nil
}
func (s *Store) CreateJob(ctx context.Context, jobs models.NewJob, userID string) (models.Job, error) {

	job, err := s.UserRepo.CreateJob(ctx, jobs)
	if err != nil {
		return models.Job{}, err
	}

	return job, nil
}
func (s *Store) ListJobs(ctx context.Context, companyID uint, userid string) ([]models.Job, error) {
	jobs, err := s.UserRepo.ViewJobByCompanyId(ctx, companyID)
	if err != nil {
		return jobs, err
	}

	return jobs, nil
}
func (s *Store) AllJob(ctx context.Context, userId string) ([]models.Job, error) {
	jobs, err := s.UserRepo.FindAllJobs(ctx)
	if err != nil {
		return []models.Job{}, err
	}

	return jobs, nil
}
func (s *Store) JobsByID(ctx context.Context, jobID uint64, userId string) (models.Job, error) {
	job, err := s.UserRepo.ViewJobDetailsById(ctx, jobID)
	if err != nil {
		return models.Job{}, err

	}
	return job, nil
}
func (s *Store) CriteriaMeets(ctx context.Context, applicant []models.Application) ([]models.Application, error) {
	ch := make(chan models.Application)
	var wg sync.WaitGroup

	for _, application := range applicant {
		wg.Add(1)
		go func(app models.Application) {
			defer wg.Done()
			job, err := s.UserRepo.GetJobById(ctx, app.JobID)
			if err != nil {
				return
			}
			if CriteriaCheck(app, job) {
				ch <- app
			}
		}(application)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := []models.Application{}
	for app := range ch {
		result = append(result, app)
	}

	return result, nil
}
func CriteriaCheck(app models.Application, job models.Job) bool {

	//Compare Min Notice Period
	a, _ := strconv.Atoi(job.MinNoticePeriod)
	if app.NoticePeriod < a {
		return false
	}

	// Compare Max Notice Period
	a, _ = strconv.Atoi(job.MaxNoticePeriod)
	if app.NoticePeriod > a {
		return false
	}

	// Compare Budget
	if app.Budget > job.Budget {
		return false
	}
	//
	// Compare Job Locations (Assuming both are ordered lists)
	//if app.JobLocations[0] != job.JobLocations[0] {
	//	return false
	//}
	//for _,v:=range app.JobLocations{
	//	for _,c:=range job.JobLocations{
	//		if
	//	}
	//}

	//for i, loc := range app.JobLocations {
	//	if loc.Location != job.JobLocations[i].Location {
	//		return false
	//	}
	//}
	//
	//// Compare Technologies (Assuming all three are required to match)
	//if len(app.Technology) != len(job.Technology) {
	//	return false
	//}
	//
	//// Compare Work Modes (Assuming both are ordered lists)
	//if len(app.WorkMode) != len(job.WorkMode) {
	//	return false
	//}
	//
	// Compare Max Experience
	if app.Exp != job.MaxExp {
		return false
	}
	//
	//// Compare Qualifications (Assuming all five are required to match)
	//if len(app.Qualification) != len(job.Qualification) {
	//	return false
	//}
	//
	// Compare Shift
	//jobLocations := make(map[string]bool)
	//for _, jl := range job.JobLocations {
	//	jobLocations[jl.Location] = true
	//}
	//for _, jl := range app.JobLocations {
	//	if !jobLocations[jl.Location] {
	//		return false
	//	}
	//}
	value := []string{}
	for _, v := range app.JobLocations {
		for _, aa := range v {

		}
	}

	if len(job.JobType) != len(app.JobType) {
		return false
	}
	if len(job.Shift) != len(app.Shift) {
		return false
	}

	return true
}
