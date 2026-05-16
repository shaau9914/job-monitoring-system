package service

import (
	"database/sql"

	"portfolio-api/model"
	"portfolio-api/repository"
)

type JobService struct {
	repo *repository.JobRepository
}

func NewJobService(repo *repository.JobRepository) *JobService {
	return &JobService{repo: repo}
}

func (s *JobService) GetJobs() ([]model.JobListResponse, error) {
	return s.repo.FindAllJobs()
}

func (s *JobService) GetJobDetail(id int) (*model.JobDetailResponse, error) {
	job, err := s.repo.FindJobByID(id)
	if err != nil {
		return nil, err
	}

	if job.Status == "failed" && job.ErrorCode != nil {
		runBook, err := s.repo.FindRunBookByErrorCode(*job.ErrorCode)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if err == nil {
			job.RunBook = runBook
		}
	}

	return job, nil
}