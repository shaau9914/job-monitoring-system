package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

const dbPath = `D:\portfolio\DB\20260429_portfolioDB.db`

type JobListResponse struct {
	ID         int     `json:"id"`
	JobName    string  `json:"job_name"`
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

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/jobs", jobsHandler)
	http.HandleFunc("/jobs/", jobDetailHandler)

	log.Println("server start: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

func jobsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`
		SELECT id, job_name, status, executed_at, error_code
		FROM job_result
		ORDER BY id ASC
	`)
	if err != nil {
		http.Error(w, "failed to get job list", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	jobs := []JobListResponse{}

	for rows.Next() {
		var job JobListResponse

		err := rows.Scan(
			&job.ID,
			&job.JobName,
			&job.Status,
			&job.ExecutedAt,
			&job.ErrorCode,
		)
		if err != nil {
			http.Error(w, "failed to scan job list", http.StatusInternalServerError)
			return
		}

		jobs = append(jobs, job)
	}

	writeJSON(w, jobs)
}

func jobDetailHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := strings.TrimPrefix(r.URL.Path, "/jobs/")
	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	var job JobDetailResponse

	err = db.QueryRow(`
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

	if err == sql.ErrNoRows {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get job detail", http.StatusInternalServerError)
		return
	}

	if job.Status == "failed" && job.ErrorCode != nil {
		var runBook RunBookResponse

		err = db.QueryRow(`
			SELECT title, cause, action
			FROM run_books
			WHERE error_code = ?
		`, *job.ErrorCode).Scan(
			&runBook.Title,
			&runBook.Cause,
			&runBook.Action,
		)

		if err == nil {
			job.RunBook = &runBook
		}
	}

	writeJSON(w, job)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
