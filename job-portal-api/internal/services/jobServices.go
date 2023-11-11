package services

import (
	"context"
	"encoding/json"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"job-portal-api/internal/models"
	"strconv"
	"sync"
	"time"
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
	rd := redis()

	for _, application := range applicant {
		wg.Add(1)
		go func(app models.Application) {
			defer wg.Done()
			key := strconv.Itoa(int(app.JobID))
			fmt.Println(key)
			job, err := CheckRedisKey(rd, key)
			if err != nil {
				jobs, err := s.UserRepo.GetJobById(ctx, app.JobID)
				fmt.Println("[[[[[[[", job, "[[[[[[[", err, "]]]]]]]]]]]]]]]]]]]]]]]]")
				if err != nil {
					return
				}
				SetRedisKey(rd, key, jobs)
				job = jobs
				fmt.Println("[[[[[[[", job)
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

func redis() *redis2.Client {
	rdb := redis2.NewClient(&redis2.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB

	})
	fmt.Println(rdb)
	return rdb
}

func CheckRedisKey(rdb *redis2.Client, key string) (models.Job, error) {
	var ctx = context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err == redis2.Nil {
		return models.Job{}, err

	}
	fmt.Println("[]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]", val)
	var job models.Job
	err = json.Unmarshal([]byte(val), &job)
	if err != nil {
		log.Err(err)
	}
	return job, nil
}
func SetRedisKey(rdb *redis2.Client, key string, value models.Job) {
	var ctx = context.Background()
	jobdata, err := json.Marshal(value)
	if err != nil {
		log.Err(err)
		return
	}
	data := string(jobdata)
	err = rdb.Set(ctx, key, data, 10*time.Minute).Err()
	if err != nil {
		log.Err(err)
		return
	}

}
func CriteriaCheck(app models.Application, job models.Job) bool {

	//for location
	count1 := 0
	for _, v := range app.LocationIDs {
		fmt.Println("application", v)
		for _, c := range job.JobLocations {
			fmt.Println("job", c.ID)
			if v == c.ID {
				count1 += 1
			}
		}

	}
	if count1 == 0 {
		return false
	}

	//for qualification

	count2 := 1
	for i, _ := range app.QualificationIDs {
		count2 += i
	}
	if count2 < 2 {
		return false
	}

	// Compare Budget
	if app.Budget > job.Budget {
		return false
	}

	if app.Exp > job.MaxExp {
		return false
	}

	//for notice period
	if app.NoticePeriod < job.MinNoticePeriod {
		return false
	}
	if app.NoticePeriod > job.MaxNoticePeriod {
		return false
	}

	//for technology
	count3 := 1

	for i, _ := range app.TechnologyIDs {
		count3 += i
	}
	if count3 == 1 {
		return false
	}
	//for job type

	if len(app.JobType) != len(job.JobType) {
		return false
	}
	//for shift type

	if len(app.Shift) != len(job.Shift) {
		return false
	}
	if len(app.WorkModeIDs) == 0 {
		return false
	}

	return true

}
