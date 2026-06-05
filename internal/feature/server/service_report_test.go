package server

import (
	"context"
	"testing"
	"time"

	"github.com/HenryNg101/server-management-system/internal/model"
)

type fakeElastic struct{}

func (f *fakeElastic) BulkInsertStatus(ctx context.Context, s []*model.Server) error {
	return nil
}

func (f *fakeElastic) GetDailyUptime(ctx context.Context, start, end time.Time, topN int) (map[uint]float64, error) {
	return map[uint]float64{1: 0.99}, nil
}

type fakeMailer struct {
	called bool
}

func (f *fakeMailer) Send(to []string, subject, body string) error {
	f.called = true
	return nil
}

func TestElasticBulkInsert(t *testing.T) {
	repo := NewFakeRepo()
	elastic := &fakeElastic{}
	mailer := &fakeMailer{}

	svc := NewService(repo, elastic, mailer)

	server := model.Server{Status: true}
	err := svc.ElasticBulkInsert(context.Background(), []*model.Server{
		&server,
	})

	if err != nil {
		t.Fatalf("bulk insert should have been successful")
	}
}

func TestSendReports_NoEmail(t *testing.T) {
	repo := NewFakeRepo()
	elastic := &fakeElastic{}
	mailer := &fakeMailer{}

	svc := NewService(repo, elastic, mailer)

	repo.Create(context.Background(), &model.Server{Status: true})

	report, err := svc.SendReports(time.Now(), time.Now(), 10, nil, context.Background())
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if report.TotalServers == 0 {
		t.Fatalf("expected stats")
	}

	if mailer.called {
		t.Fatalf("mailer should not be called")
	}
}

func TestSendReports_WithEmail(t *testing.T) {
	repo := NewFakeRepo()
	elastic := &fakeElastic{}
	mailer := &fakeMailer{}

	svc := NewService(repo, elastic, mailer)

	repo.Create(context.Background(), &model.Server{Status: true})

	emails := []string{"test@example.com"}

	_, err := svc.SendReports(time.Now(), time.Now(), 10, &emails, context.Background())
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if !mailer.called {
		t.Fatalf("expected mailer to be called")
	}
}
