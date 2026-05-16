package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"portfolio-api/handler"
	"portfolio-api/repository"
	"portfolio-api/service"

	_ "modernc.org/sqlite"
)

func main() {
		dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = `D:\portfolio\DB\20260429_portfolioDB.db`
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jobRepo := repository.NewJobRepository(db)
	jobService := service.NewJobService(jobRepo)
	jobHandler := handler.NewJobHandler(jobService)

	http.HandleFunc("/health", jobHandler.Health)
	http.HandleFunc("/jobs", jobHandler.GetJobs)
	http.HandleFunc("/jobs/", jobHandler.GetJobDetail)

	log.Println("server start: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}