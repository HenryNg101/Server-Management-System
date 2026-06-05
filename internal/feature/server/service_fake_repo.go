package server

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/HenryNg101/server-management-system/internal/model"
)

type fakeRepo struct {
	servers map[uint]*model.Server
	nextID  uint
	dbExist bool
}

func NewFakeRepo() *fakeRepo {
	return &fakeRepo{
		servers: make(map[uint]*model.Server),
		nextID:  0,
		dbExist: true,
	}
}

// --------------------
// FIND ALL (basic filtering + pagination)
// --------------------
func (f *fakeRepo) FindAll(ctx context.Context, q GetServersQuery) ([]model.Server, int64, error) {
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
func (f *fakeRepo) Create(ctx context.Context, s *model.Server) (*model.Server, error) {
	if !f.dbExist {
		return nil, errors.New("db error")
	}
	f.nextID++
	s.ID = f.nextID
	f.servers[s.ID] = s
	return s, nil
}

// --------------------
// FIND BY ID
// --------------------
func (f *fakeRepo) FindByID(ctx context.Context, id uint, out *model.Server) (*model.Server, error) {
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
func (f *fakeRepo) Update(ctx context.Context, server *model.Server) (*model.Server, error) {
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
func (f *fakeRepo) ExistsByID(ctx context.Context, id uint) (bool, error) {
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
func (f *fakeRepo) Delete(ctx context.Context, id uint) error {
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
func (f *fakeRepo) BulkUpdateStatus(ctx context.Context, results []*model.Server) error {
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
func (f *fakeRepo) GetStats(ctx context.Context) (total, up, down int64, err error) {
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

func (f *fakeRepo) Seed(servers ...*model.Server) error {
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

func (f *fakeRepo) SetDbExists(v bool) {
	f.dbExist = v
}
