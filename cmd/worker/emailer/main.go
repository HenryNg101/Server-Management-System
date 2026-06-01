package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/HenryNg101/server-management-system/internal/app"
)

func fetchEmails(application *app.Application) (*[]string, error) {
	users, err := (*application.UserService).GetUsers()
	if err != nil {
		return nil, err
	}
	userEmails := make([]string, 0, len(*users))
	for _, user := range *users {
		userEmails = append(userEmails, user.Email)
	}
	return &userEmails, nil
}

// --- Scheduler ---
func main() {
	// Create app
	newApplication, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	// Create a ticker with a channel to tick every 5 seconds
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Sending report emails...")

			// Fetch all emails
			emailAddresses, err := fetchEmails(newApplication)
			if err != nil {
				log.Fatal("Error while fetching emails: ", err)
			}

			// Send emails
			endReportTime := time.Now()
			startReportTime := endReportTime.Add(-24 * time.Hour)
			_, err = (*newApplication.ServerService).SendReports(
				startReportTime, endReportTime, 10, emailAddresses, ctx,
			)
			if err != nil {
				log.Println("cron report error:", err)
			}

		case <-ctx.Done():
			return
		}
	}
}
