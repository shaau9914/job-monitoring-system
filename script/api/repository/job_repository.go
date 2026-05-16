package repository

import (
	"database/sql"

	"portfolio-api/model"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) FindAllJobs() ([]model.JobListResponse, error) {
	rows, err := r.db.Query(`
		SELECT id, job_name, status, executed_at, error_code
		FROM job_result
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []model.JobListResponse{}

	for rows.Next() {
		var job model.JobListResponse

		err := rows.Scan(
			&job.ID,
			&job.JobName,
			&job.Status,
			&job.ExecutedAt,
			&job.ErrorCode,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *JobRepository) FindJobByID(id int) (*model.JobDetailResponse, error) {
	var job model.JobDetailResponse

	err := r.db.QueryRow(`
		SELECT id, job_name, status, executed_at, error_code, error_message
		FROM job_result
		WHERE id = ?
	`, id).Scan(
		&job.ID,
		&job.JobName,
		&job.Status,
		&job.ExecutedAt,
		&job.ErrorCode,
		&job.ErrorMessage,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *JobRepository) FindRunBookByErrorCode(errorCode string) (*model.RunBookResponse, error) {
	var runBook model.RunBookResponse

	err := r.db.QueryRow(`
		SELECT title, cause, action
		FROM run_books
		WHERE error_code = ?
	`, errorCode).Scan(
		&runBook.Title,
		&runBook.Cause,
		&runBook.Action,
	)

	if err != nil {
		return nil, err
	}

	return &runBook, nil
}