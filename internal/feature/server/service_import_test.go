package server

import (
	"context"
	"strings"
	"testing"
)

func TestImportServers(t *testing.T) {
	svc, _ := setupService()

	csvData := `name,status,ipv4_address,port,protocol
server1,true,127.0.0.1,80,tcp
server2,false,192.168.1.1,443,tcp
`

	res, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.SuccessCount != 2 {
		t.Fatalf("expected 2 successes")
	}
}

func TestImportServers_InvalidRows(t *testing.T) {
	svc, _ := setupService()

	csvData := `name,status,ipv4_address,port,protocol
server1,invalid,127.0.0.1,80,tcp
server2,true,invalid_ip,443,tcp
`

	res, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.FailedCount != 2 {
		t.Fatalf("expected 2 failures")
	}
}

func TestImportServers_EmptyCSV(t *testing.T) {
	svc, _ := setupService()

	csvData := `name,status,ipv4_address,port,protocol`

	_, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
	if err == nil {
		t.Fatalf("expected error for empty csv")
	}
}

func TestImportServers_PartialSuccess(t *testing.T) {
	svc, _ := setupService()

	csvData := `name,status,ipv4_address,port,protocol
ok,true,127.0.0.1,80,tcp
bad,true,invalid_ip,80,tcp
`

	res, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if res.SuccessCount != 1 || res.FailedCount != 1 {
		t.Fatalf("expected partial success")
	}
}

func TestImportServers_RepoFailure(t *testing.T) {
	repo := NewFakeRepo()
	repo.SetDbExists(false)

	svc := NewService(repo, nil, nil)

	csvData := `name,status,ipv4_address,port,protocol
server1,true,127.0.0.1,80,tcp
`

	res, err := svc.ImportServers(context.Background(), strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if res.FailedCount != 1 {
		t.Fatalf("expected failure due to repo error")
	}
}
