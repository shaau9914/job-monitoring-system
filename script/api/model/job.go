package model

type JobListResponse struct {
	ID         int     `json:"id"`
	JobName   string  `json:"job_name"`
	Status     string  `json:"status"`
	ExecutedAt string  `json:"executed_at"`
	ErrorCode  *string `json:"error_code"`
}

type RunBookResponse struct {
	Title  string `json:"title"`
	Cause  string `json:"cause"`
	Action string `json:"action"`
}

type JobDetailResponse struct {
	ID           int              `json:"id"`
	JobName      string           `json:"job_name"`
	Status       string           `json:"status"`
	ExecutedAt   string           `json:"executed_at"`
	ErrorCode    *string          `json:"error_code"`
	ErrorMessage *string          `json:"error_message"`
	RunBook      *RunBookResponse `json:"runbook"`
}