package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/HenryNg101/server-service/internal/model"
	"github.com/go-openapi/testify/v2/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=postgres password=secret dbname=test_db sslmode=disable",
		host,
		port.Port(),
	)

	// Retry until DB is ready (important!)
	var db *gorm.DB
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Server{})
	require.NoError(t, err)

	return db
}

func seedServers(db *gorm.DB) []model.Server {
	servers := []model.Server{
		{Name: "s1", IPv4: "127.0.0.1"},
		{Name: "s2", IPv4: "127.0.0.2"},
	}

	for i := range servers {
		db.Create(&servers[i])
	}
	return servers
}

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	server := &model.Server{
		Name: "unique",
		IPv4: "127.0.0.1",
	}

	res, err := repo.Create(context.Background(), server)
	require.NoError(t, err)
	require.NotNil(t, res)

	// duplicate
	_, err = repo.Create(context.Background(), server)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already existed")
}

func TestRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	seedServers(db)

	// status := true
	page := 1
	size := 1
	// protocol := "http"
	name := "s1"

	q := GetServersQuery{
		// Status:   &status,
		Page:     &page,
		PageSize: &size,
		SortBy:   "id",
		Order:    "DESC",
		// Protocol: &protocol,
		Name: &name,
	}

	servers, total, err := repo.FindAll(context.Background(), q)

	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, servers, 1)
	// require.True(t, servers[0].Status)
}

func TestRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	servers := seedServers(db)

	// var s *model.Server
	res, err := repo.FindByID(context.Background(), servers[0].ID)

	require.NoError(t, err)
	require.Equal(t, "s1", res.Name)
}

func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	servers := seedServers(db)

	s := servers[0]
	s.Name = "updated"

	updated, err := repo.Update(context.Background(), &s)

	require.NoError(t, err)
	require.Equal(t, "updated", updated.Name)
}

func TestRepository_ExistsByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	servers := seedServers(db)

	exists, err := repo.ExistsByID(context.Background(), servers[0].ID)
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = repo.ExistsByID(context.Background(), 999)
	require.NoError(t, err)
	require.False(t, exists)
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	servers := seedServers(db)

	err := repo.Delete(context.Background(), servers[0].ID)
	require.NoError(t, err)

	exists, _ := repo.ExistsByID(context.Background(), servers[0].ID)
	require.False(t, exists)
}

// func TestRepository_BulkUpdateStatus(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := NewRepository(db)

// 	servers := seedServers(db)

// 	now := time.Now()

// 	updates := []*model.Server{
// 		{ID: servers[0].ID, Status: false, LastUpdated: now},
// 		{ID: servers[1].ID, Status: true, LastUpdated: now},
// 	}

// 	err := repo.BulkUpdateStatus(context.Background(), updates)
// 	require.NoError(t, err)

// 	var updated []model.Server
// 	db.Find(&updated)

// 	require.False(t, updated[0].Status)
// 	require.True(t, updated[1].Status)
// }

// func TestRepository_GetStats(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := NewRepository(db)

// 	seedServers(db)

// 	total, up, down, err := repo.GetStats(context.Background())

// 	require.NoError(t, err)
// 	require.Equal(t, int64(2), total)
// 	require.Equal(t, int64(1), up)
// 	require.Equal(t, int64(1), down)
// }
