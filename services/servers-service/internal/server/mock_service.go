package server

import (
	"context"

	"github.com/HenryNg101/server-service/internal/model"
)

type MockServerService struct {
	getServersFn func(context.Context, *UserContext, GetServersQuery) (*PaginatedServers, error)
	createFn     func(context.Context, CreateServerRequest, uint) (*model.Server, error)
	getFn        func(context.Context, uint, *UserContext, *model.Server) (*model.Server, error)
	updateFn     func(context.Context, uint, *UserContext, UpdateServerRequest) (*model.Server, error)
	deleteFn     func(context.Context, uint, *UserContext) error
	// createAgentFn func(ctx context.Context, serverID uint) (*model.Agent, string, error)
	// importFn      func(context.Context, io.Reader) (*ImportServersResponse, error)
	// reportFn      func(time.Time, time.Time, int, *[]string, context.Context) (*Report, error)
}

func (m *MockServerService) GetServers(ctx context.Context, userCtx *UserContext, q GetServersQuery) (*PaginatedServers, error) {
	if m.getServersFn != nil {
		return m.getServersFn(ctx, userCtx, q)
	}
	return nil, nil
}

func (m *MockServerService) CreateServer(ctx context.Context, req CreateServerRequest, userID uint) (*model.Server, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req, userID)
	}
	return nil, nil
}

func (m *MockServerService) GetServer(ctx context.Context, serverID uint, userCtx *UserContext, server *model.Server) (*model.Server, error) {
	if m.getFn != nil {
		return m.getFn(ctx, serverID, userCtx, server)
	}
	return nil, nil
}

func (m *MockServerService) UpdateServer(ctx context.Context, serverID uint, userCtx *UserContext, req UpdateServerRequest) (*model.Server, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, serverID, userCtx, req)
	}
	return nil, nil
}

func (m *MockServerService) DeleteServer(ctx context.Context, serverID uint, userCtx *UserContext) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, serverID, userCtx)
	}
	return nil
}

// func (m *MockServerService) CreateAgent(ctx context.Context, serverID uint) (*model.Agent, string, error) {
// 	return m.createAgentFn(ctx, serverID)
// }

// func (m *MockServerService) ImportServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error) {
// 	return m.importFn(ctx, r)
// }

// func (m *MockServerService) BulkUpdateServersStatuses(ctx context.Context, servers []*model.Server) error {
// 	return nil
// }

// func (m *MockServerService) ElasticBulkInsert(ctx context.Context, serversResults []*model.Server) error {
// 	return nil
// }

// func (m *MockServerService) SendReports(startTime time.Time, endTime time.Time, topN int, emailsList *[]string, ctx context.Context) (*Report, error) {
// 	return m.reportFn(startTime, endTime, topN, emailsList, ctx)
// }
