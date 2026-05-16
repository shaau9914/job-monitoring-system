package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"portfolio-api/service"
)

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(service *service.JobService) *JobHandler {
	return &JobHandler{service: service}
}

func (h *JobHandler) Health(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobs, err := h.service.GetJobs()
	if err != nil {
		http.Error(w, "failed to get job list", http.StatusInternalServerError)
		return
	}

	writeJSON(w, jobs)
}

func (h *JobHandler) GetJobDetail(w http.ResponseWriter, r *http.Request) {
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

	job, err := h.service.GetJobDetail(id)
	if err == sql.ErrNoRows {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get job detail", http.StatusInternalServerError)
		return
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