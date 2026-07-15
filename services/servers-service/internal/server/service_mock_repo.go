package server

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/HenryNg101/server-service/internal/model"
)

type mockRepo struct {
	servers     map[uint]*model.Server
	agents      map[uint]*model.Agent
	nextID      uint
	nextAgentID uint
	dbExist     bool
}

func NewFakeRepo() *mockRepo {
	return &mockRepo{
		servers:     make(map[uint]*model.Server),
		agents:      make(map[uint]*model.Agent),
		nextID:      0,
		nextAgentID: 0,
		dbExist:     true,
	}
}

// --------------------
// FIND ALL (basic filtering + pagination)
// --------------------
func (f *mockRepo) FindAll(ctx context.Context, q GetServersQuery) ([]model.Server, int64, error) {
	var result []model.Server
	if !f.dbExist {
		return result, 0, errors.New("db error")
	}

	for _, s := range f.servers {
		if q.Status != nil && s.Status != *q.Status {
			continue
		}
		if q.Protocol != nil && s.Protocol != *q.Protocol {
			continue
		}
		if q.Name != nil && !strings.Contains(strings.ToLower(s.Name), strings.ToLower(*q.Name)) {
			continue
		}
		result = append(result, *s)
	}

	// sorting (very basic)
	if q.SortBy == "name" {
		sort.Slice(result, func(i, j int) bool {
			if q.Order == "desc" {
				return result[i].Name > result[j].Name
			}
			return result[i].Name < result[j].Name
		})
	}

	total := int64(len(result))

	// pagination
	if q.Page != nil && q.PageSize != nil {
		start := (*q.Page - 1) * (*q.PageSize)
		end := start + *q.PageSize

		if start > len(result) {
			return []model.Server{}, total, nil
		}
		if end > len(result) {
			end = len(result)
		}

		result = result[start:end]
	}

	return result, total, nil
}

// --------------------
// CREATE
// --------------------
func (f *mockRepo) Create(ctx context.Context, s *model.Server) (*model.Server, error) {
	if !f.dbExist {
		return nil, errors.New("db error")
	}
	f.nextID++
	s.ID = f.nextID
	f.servers[s.ID] = s
	return s, nil
}

func (f *mockRepo) BulkUpsert(ctx context.Context, servers []*model.Server) error {
	for _, s := range servers {
		if !f.dbExist {
			return errors.New("db error")
		}
		f.nextID++
		s.ID = f.nextID
		f.servers[s.ID] = s
	}
	return nil
}

// --------------------
// CREATE AGENT
// --------------------
func (f *mockRepo) CreateAgent(ctx context.Context, agent *model.Agent) (*model.Agent, error) {
	if !f.dbExist {
		return nil, errors.New("db error")
	}
	f.nextAgentID++
	agent.ID = f.nextAgentID
	f.agents[agent.ID] = agent
	return agent, nil
}

// --------------------
// FIND BY ID
// --------------------
func (f *mockRepo) FindByID(ctx context.Context, id uint, out *model.Server) (*model.Server, error) {
	if !f.dbExist {
		return nil, errors.New("db error")
	}
	s, ok := f.servers[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

// --------------------
// UPDATE
// --------------------
func (f *mockRepo) Update(ctx context.Context, server *model.Server) (*model.Server, error) {
	if !f.dbExist {
		return nil, errors.New("db error")
	}
	_, ok := f.servers[server.ID]
	if !ok {
		return nil, errors.New("not found")
	}
	f.servers[server.ID] = server
	return server, nil
}

// --------------------
// EXISTS
// --------------------
func (f *mockRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
	if !f.dbExist {
		return false, errors.New("db error")
	}
	_, ok := f.servers[id]
	if !ok {
		return false, errors.New("User does not exist")
	}
	return ok, nil
}

// --------------------
// DELETE
// --------------------
func (f *mockRepo) Delete(ctx context.Context, id uint) error {
	if !f.dbExist {
		return errors.New("db error")
	}
	_, ok := f.servers[id]
	if !ok {
		return errors.New("not found")
	}
	delete(f.servers, id)
	return nil
}

// --------------------
// BULK UPDATE STATUS
// --------------------
func (f *mockRepo) BulkUpdateStatus(ctx context.Context, results []*model.Server) error {
	if !f.dbExist {
		return errors.New("db error")
	}
	for _, r := range results {
		s, ok := f.servers[r.ID]
		if !ok {
			continue
		}
		s.Status = r.Status
	}
	return nil
}

// --------------------
// STATS
// --------------------
func (f *mockRepo) GetStats(ctx context.Context) (total, up, down int64, err error) {
	if !f.dbExist {
		return 0, 0, 0, errors.New("db error")
	}
	for _, s := range f.servers {
		total++
		if s.Status {
			up++
		} else {
			down++
		}
	}
	return
}

func (f *mockRepo) Seed(servers ...*model.Server) error {
	if !f.dbExist {
		return errors.New("db error")
	}
	for _, s := range servers {
		f.nextID++
		s.ID = f.nextID
		f.servers[s.ID] = s
	}
	return nil
}

func (f *mockRepo) SetDbExists(v bool) {
	f.dbExist = v
}
