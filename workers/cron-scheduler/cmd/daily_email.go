package main

import (
	"context"
	"log"
	"time"

	"github.com/HenryNg101/cron-scheduler/internal/app"
)

func fetchEmails(application *app.Application) (*[]string, error) {
	users, err := application.UserService.GetUsers()
	if err != nil {
		return nil, err
	}

	userEmails := make([]string, 0, len(*users))
	for _, user := range *users {
		userEmails = append(userEmails, user.Email)
	}

	return &userEmails, nil
}

func runEmailJob(app *app.Application, ctx context.Context) {
	log.Println("Sending report emails...")

	emails, err := fetchEmails(app)
	if err != nil {
		log.Println("fetch emails error:", err)
		return
	}

	end := time.Now()
	start := end.Add(-24 * time.Hour)

	_, err = app.MonitoringService.SendReports(start, end, 10, emails, ctx)
	if err != nil {
		log.Println("email job error:", err)
	}
}
