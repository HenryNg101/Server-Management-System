package server

import (
	"context"
	"testing"

	"github.com/HenryNg101/servers-service/internal/model"
)

// type fakeElastic struct{}

// func (f *fakeElastic) BulkInsertStatus(ctx context.Context, s []*model.Server) error {
// 	return nil
// }

// func (f *fakeElastic) GetDailyUptime(ctx context.Context, start, end time.Time, topN int) (map[uint]*ServerPullStats, error) {
// 	return map[uint]*ServerPullStats{1: {Uptime: 0.99}}, nil
// }

// type fakeMailer struct {
// 	called bool
// }

// func (f *fakeMailer) Send(to []string, subject, body string) error {
// 	f.called = true
// 	return nil
// }

func testUserCtx() *UserContext {
	return &UserContext{
		UserID: 1,
		Role:   "user",
	}
}

func setupService() (Service, *mockRepo) {
	repo := NewFakeRepo()
	svc := NewService(repo)
	return svc, repo
}

// func TestBulkUpdateServersStatuses(t *testing.T) {
// 	svc, repo := setupService()

// 	var servers []*model.Server

// 	for i := 0; i < 1500; i++ {
// 		s, _ := repo.Create(context.Background(), &model.Server{
// 			Status: false,
// 		})
// 		servers = append(servers, &model.Server{
// 			ID:     s.ID,
// 			Status: true,
// 		})
// 	}

// 	err := svc.BulkUpdateServersStatuses(context.Background(), servers)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	for _, s := range servers {
// 		if !s.Status {
// 			t.Fatalf("expected all servers to be updated")
// 		}
// 	}
// }

// func TestBulkUpdate_EmptyList(t *testing.T) {
// 	svc, _ := setupService()

// 	err := svc.BulkUpdateServersStatuses(context.Background(), []*model.Server{})
// 	if err != nil {
// 		t.Fatalf("should not fail on empty input")
// 	}
// }

// func TestBulkUpdate_SmallChunk(t *testing.T) {
// 	svc, repo := setupService()

// 	s, _ := repo.Create(context.Background(), &model.Server{Status: false})

// 	err := svc.BulkUpdateServersStatuses(context.Background(), []*model.Server{
// 		{ID: s.ID, Status: true},
// 	})

// 	if err != nil {
// 		t.Fatalf("unexpected error")
// 	}

// 	if !repo.servers[s.ID].Status {
// 		t.Fatalf("status not updated")
// 	}
// }

func TestCreateServer(t *testing.T) {
	svc, _ := setupService()
	userCtx := testUserCtx()

	res, err := svc.CreateServer(context.Background(), CreateServerRequest{
		Name: "test",
		IPv4: "127.0.0.1",
	}, userCtx.UserID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID == 0 {
		t.Fatalf("expected ID to be set")
	}
}

func TestGetServer(t *testing.T) {
	svc, repo := setupService()
	userCtx := testUserCtx()

	created, _ := repo.Create(context.Background(), &model.Server{
		Name:   "test",
		UserID: userCtx.UserID,
	})

	res, err := svc.GetServer(context.Background(), created.ID, userCtx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID != created.ID {
		t.Fatalf("expected same ID")
	}
}

func TestGetServers(t *testing.T) {
	svc, repo := setupService()
	userCtx := testUserCtx()

	repo.Create(context.Background(), &model.Server{Name: "a", UserID: userCtx.UserID})
	repo.Create(context.Background(), &model.Server{Name: "b", UserID: userCtx.UserID})

	page := 1
	size := 10

	res, err := svc.GetServers(context.Background(), userCtx, GetServersQuery{
		Page:     &page,
		PageSize: &size,
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if res.Total != 2 {
		t.Fatalf("expected total 2")
	}
}

func TestGetServer_NotFound(t *testing.T) {
	svc, _ := setupService()
	userCtx := testUserCtx()

	_, err := svc.GetServer(context.Background(), 999, userCtx, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetServers_FilteringAndPagination(t *testing.T) {
	repo := NewFakeRepo()
	userCtx := testUserCtx()

	err := repo.Seed(
		&model.Server{Name: "alpha", UserID: 1},
		&model.Server{Name: "beta", UserID: 1},
		&model.Server{Name: "gamma", UserID: 2},
	)

	svc := NewService(repo)

	// status := true
	// protocol := "tcp"
	page := 1
	size := 1

	res, err := svc.GetServers(context.Background(), userCtx, GetServersQuery{
		// Status:   &status,
		// Protocol: &protocol,
		Page:     &page,
		PageSize: &size,
		SortBy:   "name",
		Order:    "asc",
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if res.Total != 2 {
		t.Fatalf("expected 2 filtered results")
	}

	if len(res.Servers) != 1 {
		t.Fatalf("expected pagination to limit results")
	}
}

func TestUpdateServer(t *testing.T) {
	svc, repo := setupService()
	userCtx := testUserCtx()

	created, _ := repo.Create(context.Background(), &model.Server{
		Name:   "old",
		UserID: userCtx.UserID,
	})

	newName := "new"

	res, err := svc.UpdateServer(context.Background(), created.ID, userCtx, UpdateServerRequest{
		Name: &newName,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Name != "new" {
		t.Fatalf("expected name updated")
	}
}

func TestUpdateServer_NoChanges(t *testing.T) {
	svc, repo := setupService()
	userCtx := testUserCtx()

	created, _ := repo.Create(context.Background(), &model.Server{})

	_, err := svc.UpdateServer(context.Background(), created.ID, userCtx, UpdateServerRequest{})
	if err == nil {
		t.Fatalf("expected error when nothing updated")
	}
}

func TestUpdateServer_AllFields(t *testing.T) {
	svc, repo := setupService()
	userCtx := testUserCtx()

	created, _ := repo.Create(context.Background(), &model.Server{
		Name: "old",
		IPv4: "127.0.0.1",
	})

	name := "new"
	ip := "192.168.1.1"

	res, err := svc.UpdateServer(context.Background(), created.ID, userCtx, UpdateServerRequest{
		Name: &name,
		IPv4: &ip,
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if res.Name != name {
		t.Fatalf("fields not updated correctly")
	}
}

func TestUpdateServer_NotFound(t *testing.T) {
	svc, _ := setupService()
	userCtx := testUserCtx()

	name := "new"

	_, err := svc.UpdateServer(context.Background(), 999, userCtx, UpdateServerRequest{
		Name: &name,
	})

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDeleteServer(t *testing.T) {
	svc, repo := setupService()
	userCtx := testUserCtx()

	created, _ := repo.Create(context.Background(), &model.Server{})

	err := svc.DeleteServer(context.Background(), created.ID, userCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteServer_NotFound(t *testing.T) {
	svc, _ := setupService()
	userCtx := testUserCtx()

	err := svc.DeleteServer(context.Background(), 999, userCtx)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDeleteServer_RepoError(t *testing.T) {
	repo := NewFakeRepo()
	repo.SetDbExists(false)
	userCtx := testUserCtx()

	svc := NewService(repo)

	err := svc.DeleteServer(context.Background(), 1, userCtx)
	if err == nil {
		t.Fatalf("expected error")
	}
}

// func TestImportServers(t *testing.T) {
// 	svc, _ := setupService()

// 	csvData := `name,status,ipv4_address,port,protocol
// server1,true,127.0.0.1,80,tcp
// server2,false,192.168.1.1,443,tcp
// `

// 	res, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if res.SuccessCount != 2 {
// 		t.Fatalf("expected 2 successes")
// 	}
// }

// func TestImportServers_InvalidRows(t *testing.T) {
// 	svc, _ := setupService()

// 	csvData := `name,status,ipv4_address,port,protocol
// server1,invalid,127.0.0.1,80,tcp
// server2,true,invalid_ip,443,tcp
// `

// 	res, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if res.FailedCount != 2 {
// 		t.Fatalf("expected 2 failures")
// 	}
// }

// func TestImportServers_EmptyCSV(t *testing.T) {
// 	svc, _ := setupService()

// 	csvData := `name,status,ipv4_address,port,protocol`

// 	_, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
// 	if err == nil {
// 		t.Fatalf("expected error for empty csv")
// 	}
// }

// func TestImportServers_PartialSuccess(t *testing.T) {
// 	svc, _ := setupService()

// 	csvData := `name,status,ipv4_address,port,protocol
// ok,true,127.0.0.1,80,tcp
// bad,true,invalid_ip,80,tcp
// `

// 	res, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
// 	if err != nil {
// 		t.Fatalf("unexpected error")
// 	}

// 	if res.SuccessCount != 1 || res.FailedCount != 1 {
// 		t.Fatalf("expected partial success")
// 	}
// }

// func TestImportServers_RepoFailure(t *testing.T) {
// 	repo := NewFakeRepo()
// 	repo.SetDbExists(false)

// 	svc := NewService(repo, nil)

// 	csvData := `name,status,ipv4_address,port,protocol
// server1,true,127.0.0.1,80,tcp
// `

// 	res, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
// 	if err != nil {
// 		t.Fatalf("unexpected error")
// 	}

// 	if res.FailedCount != 1 {
// 		t.Fatalf("expected failure due to repo error")
// 	}
// }

// func TestElasticBulkInsert(t *testing.T) {
// 	repo := NewFakeRepo()
// 	elastic := &fakeElastic{}

// 	svc := NewService(repo, elastic)

// 	server := model.Server{Status: true}
// 	err := svc.ElasticBulkInsert(context.Background(), []*model.Server{
// 		&server,
// 	})

// 	if err != nil {
// 		t.Fatalf("bulk insert should have been successful")
// 	}
// }

// TODO: Move these tests to monitoring domain/feature when write unit tests there
// func TestSendReports_NoEmail(t *testing.T) {
// 	repo := NewFakeRepo()
// 	elastic := &fakeElastic{}
// 	mailer := &fakeMailer{}

// 	svc := NewService(repo, elastic, mailer)

// 	repo.Create(context.Background(), &model.Server{Status: true})

// 	report, err := svc.SendReports(time.Now(), time.Now(), 10, nil, context.Background())
// 	if err != nil {
// 		t.Fatalf("unexpected error")
// 	}

// 	if report.TotalServers == 0 {
// 		t.Fatalf("expected stats")
// 	}

// 	if mailer.called {
// 		t.Fatalf("mailer should not be called")
// 	}
// }

// func TestSendReports_WithEmail(t *testing.T) {
// 	repo := NewFakeRepo()
// 	elastic := &fakeElastic{}
// 	mailer := &fakeMailer{}

// 	svc := NewService(repo, elastic, mailer)

// 	repo.Create(context.Background(), &model.Server{Status: true})

// 	emails := []string{"test@example.com"}

// 	_, err := svc.SendReports(time.Now(), time.Now(), 10, &emails, context.Background())
// 	if err != nil {
// 		t.Fatalf("unexpected error")
// 	}

// 	if !mailer.called {
// 		t.Fatalf("expected mailer to be called")
// 	}
// }
