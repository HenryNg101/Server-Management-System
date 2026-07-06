package server

import (
	"context"

	"github.com/HenryNg101/server-management-system/internal/model"
)

type MockServerService struct {
	getServersFn  func(context.Context, GetServersQuery) (*PaginatedServers, error)
	createFn      func(context.Context, CreateServerRequest) (*model.Server, error)
	createAgentFn func(ctx context.Context, serverID uint) (*model.Agent, string, error)
	getFn         func(context.Context, uint, *model.Server) (*model.Server, error)
	updateFn      func(context.Context, uint, UpdateServerRequest) (*model.Server, error)
	deleteFn      func(context.Context, uint) error
	// importFn      func(context.Context, io.Reader) (*ImportServersResponse, error)
	// reportFn      func(time.Time, time.Time, int, *[]string, context.Context) (*Report, error)
}

func (m *MockServerService) GetServers(ctx context.Context, q GetServersQuery) (*PaginatedServers, error) {
	return m.getServersFn(ctx, q)
}

func (m *MockServerService) CreateServer(ctx context.Context, req CreateServerRequest) (*model.Server, error) {
	return m.createFn(ctx, req)
}

func (m *MockServerService) CreateAgent(ctx context.Context, serverID uint) (*model.Agent, string, error) {
	return m.createAgentFn(ctx, serverID)
}

func (m *MockServerService) GetServer(ctx context.Context, id uint, server *model.Server) (*model.Server, error) {
	return m.getFn(ctx, id, server)
}

func (m *MockServerService) UpdateServer(ctx context.Context, id uint, req UpdateServerRequest) (*model.Server, error) {
	return m.updateFn(ctx, id, req)
}

func (m *MockServerService) DeleteServer(ctx context.Context, id uint) error {
	return m.deleteFn(ctx, id)
}

// func (m *MockServerService) ImportServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error) {
// 	return m.importFn(ctx, r)
// }

func (m *MockServerService) BulkUpdateServersStatuses(ctx context.Context, servers []*model.Server) error {
	return nil
}

func (m *MockServerService) ElasticBulkInsert(ctx context.Context, serversResults []*model.Server) error {
	return nil
}

// func (m *MockServerService) SendReports(startTime time.Time, endTime time.Time, topN int, emailsList *[]string, ctx context.Context) (*Report, error) {
// 	return m.reportFn(startTime, endTime, topN, emailsList, ctx)
// }
