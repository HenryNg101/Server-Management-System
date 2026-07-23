package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/HenryNg101/server-service/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"gorm.io/gorm"
)

func setupRouter(s Service) *gin.Engine {
	gin.SetMode(gin.TestMode)

	h := NewHandler(s)
	r := gin.New()

	// Inject fake user context
	r.Use(func(c *gin.Context) {
		c.Set("user", &UserContext{
			UserID: 1,
			Role:   "user",
		})
	})

	r.POST("/servers", h.CreateServer)
	r.GET("/servers", h.GetServers)
	r.GET("/servers/:id", h.GetServer)
	r.PATCH("/servers/:id", h.UpdateServer)
	r.DELETE("/servers/:id", h.DeleteServer)
	// r.POST("/servers/import", h.ImportServers)
	r.GET("/servers/export", h.ExportServers)
	// r.POST("/servers/report", h.SendReports)

	return r
}

func TestHandlerCreateServer_Success(t *testing.T) {
	mock := &MockServerService{
		createFn: func(ctx context.Context, req CreateServerRequest, userID uint) (*model.Server, error) {
			return &model.Server{
				Name:   req.Name,
				IPv4:   req.IPv4,
				UserID: userID,
			}, nil
		},
	}
	r := setupRouter(mock)

	body := `{
		"name":"test",
		"ipv4": "192.168.1.1",
	}`

	req := httptest.NewRequest(http.MethodPost, "/servers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, 201, w.Code)
	require.Contains(t, w.Body.String(), "test")
}

func TestHandlerCreateServer_MissingField(t *testing.T) {
	mock := &MockServerService{
		createFn: func(ctx context.Context, req CreateServerRequest, userID uint) (*model.Server, error) {
			return &model.Server{
				Name:   req.Name,
				IPv4:   req.IPv4,
				UserID: userID,
			}, nil
		},
	}
	r := setupRouter(mock)

	body := `{
		"name":"test",
		"status":true,
	}`

	req := httptest.NewRequest(http.MethodPost, "/servers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
}

func TestHandlerGetServers_Success(t *testing.T) {
	mock := &MockServerService{
		getServersFn: func(ctx context.Context, userCtx *UserContext, q GetServersQuery) (*PaginatedServers, error) {
			return &PaginatedServers{
				Servers: []model.Server{
					{Name: "s1"},
					{Name: "s2"},
				},
			}, nil
		},
	}
	r := setupRouter(mock)

	req := httptest.NewRequest(http.MethodGet, "/servers?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "s1")
	require.Contains(t, w.Body.String(), "s2")
}

func TestHandlerGetServers_InvalidQuery(t *testing.T) {
	mock := &MockServerService{
		getServersFn: func(ctx context.Context, userCtx *UserContext, q GetServersQuery) (*PaginatedServers, error) {
			return &PaginatedServers{}, nil
		},
	}
	r := setupRouter(mock)

	req := httptest.NewRequest("GET", "/servers?status=invalid_bool", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
}

func TestGetServers_ServiceError(t *testing.T) {
	mock := &MockServerService{
		getServersFn: func(ctx context.Context, userCtx *UserContext, q GetServersQuery) (*PaginatedServers, error) {
			return nil, errors.New("db fail")
		},
	}
	r := setupRouter(mock)

	req := httptest.NewRequest("GET", "/servers", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 500, w.Code)
}

func TestHandlerGetServer_Success(t *testing.T) {
	mock := &MockServerService{
		getFn: func(ctx context.Context, id uint, userCtx *UserContext, s *model.Server) (*model.Server, error) {
			return &model.Server{Name: "server1"}, nil
		},
	}
	r := setupRouter(mock)

	req := httptest.NewRequest("GET", "/servers/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "server1")
}

func TestHandlerGetServer_InvalidID(t *testing.T) {
	mock := &MockServerService{}
	r := setupRouter(mock)

	req := httptest.NewRequest("GET", "/servers/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
}

func TestHandlerGetServer_NotFound(t *testing.T) {
	mock := &MockServerService{
		getFn: func(ctx context.Context, id uint, userCtx *UserContext, s *model.Server) (*model.Server, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	r := setupRouter(mock)

	req := httptest.NewRequest("GET", "/servers/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 404, w.Code)
}

func TestHandlerGetServer_InternalError(t *testing.T) {
	mock := &MockServerService{
		getFn: func(ctx context.Context, id uint, userCtx *UserContext, s *model.Server) (*model.Server, error) {
			return nil, errors.New("db error")
		},
	}
	r := setupRouter(mock)

	req := httptest.NewRequest("GET", "/servers/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 500, w.Code)
}

func TestHandlerUpdateServer_InvalidJSON(t *testing.T) {
	mock := &MockServerService{}
	r := setupRouter(mock)

	req := httptest.NewRequest("PATCH", "/servers/1", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
}

func TestHandlerUpdateServer_InvalidID(t *testing.T) {
	mock := &MockServerService{}

	r := setupRouter(mock)

	req := httptest.NewRequest("PATCH", "/servers/abc", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
}

func TestHandlerUpdateServer_ServiceError(t *testing.T) {
	mock := &MockServerService{
		updateFn: func(ctx context.Context, id uint, userCtx *UserContext, req UpdateServerRequest) (*model.Server, error) {
			return nil, errors.New("update fail")
		},
	}

	r := setupRouter(mock)

	req := httptest.NewRequest("PATCH", "/servers/1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, 500, w.Code)
}

func TestHandlerUpdateServer_Success(t *testing.T) {
	mock := &MockServerService{
		updateFn: func(ctx context.Context, id uint, userCtx *UserContext, req UpdateServerRequest) (*model.Server, error) {
			return &model.Server{Name: "updated"}, nil
		},
	}

	r := setupRouter(mock)

	req := httptest.NewRequest("PATCH", "/servers/1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
}

func TestHandlerDeleteServer_NotFound(t *testing.T) {
	mock := &MockServerService{
		deleteFn: func(ctx context.Context, id uint, userCtx *UserContext) error {
			return ErrNotFound
		},
	}

	r := setupRouter(mock)

	req := httptest.NewRequest("DELETE", "/servers/999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 404, w.Code)
}

func TestHandlerDeleteServer_InvalidID(t *testing.T) {
	mock := &MockServerService{}

	r := setupRouter(mock)

	req := httptest.NewRequest("DELETE", "/servers/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 400, w.Code)
}

func TestHandlerDeleteServer_InternalError(t *testing.T) {
	mock := &MockServerService{
		deleteFn: func(ctx context.Context, id uint, userCtx *UserContext) error {
			return errors.New("fail")
		},
	}

	r := setupRouter(mock)

	req := httptest.NewRequest("DELETE", "/servers/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 500, w.Code)
}

func TestHandlerDeleteServer_Success(t *testing.T) {
	mock := &MockServerService{
		deleteFn: func(ctx context.Context, id uint, userCtx *UserContext) error {
			return nil
		},
	}

	r := setupRouter(mock)

	req := httptest.NewRequest("DELETE", "/servers/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 204, w.Code)
}

// func TestHandlerImportServers_Success(t *testing.T) {
// 	csvData := `name,status,ipv4,port,protocol
// server1,true,127.0.0.1,8080,http`

// 	mock := &MockServerService{
// 		importFn: func(ctx context.Context, r io.Reader) (*ImportServersResponse, error) {
// 			return &ImportServersResponse{
// 				SuccessCount: 1,
// 			}, nil
// 		},
// 	}

// 	r := setupRouter(mock)

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	part, _ := writer.CreateFormFile("file", "test.csv")
// 	part.Write([]byte(csvData))
// 	writer.Close()

// 	req := httptest.NewRequest("POST", "/servers/import", body)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 200, w.Code)
// }

// func TestHandlerImportServers_NoFile(t *testing.T) {
// 	mock := &MockServerService{}

// 	r := setupRouter(mock)

// 	req := httptest.NewRequest("POST", "/servers/import", nil)
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 400, w.Code)
// }

// func TestHandlerImportServers_ServiceError(t *testing.T) {
// 	mock := &MockServerService{
// 		importFn: func(ctx context.Context, r io.Reader) (*ImportServersResponse, error) {
// 			return nil, errors.New("import fail")
// 		},
// 	}

// 	r := setupRouter(mock)

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)
// 	part, _ := writer.CreateFormFile("file", "test.csv")
// 	part.Write([]byte("a,b\n1,2"))
// 	writer.Close()

// 	req := httptest.NewRequest("POST", "/servers/import", body)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 500, w.Code)
// }

func TestHandlerExportServers_ServiceError(t *testing.T) {
	mock := &MockServerService{
		getServersFn: func(ctx context.Context, userCtx *UserContext, q GetServersQuery) (*PaginatedServers, error) {
			return nil, errors.New("fail")
		},
	}

	r := setupRouter(mock)

	req := httptest.NewRequest("GET", "/servers/export", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 500, w.Code)
}

func TestHandlerExportServers_Success(t *testing.T) {
	mock := &MockServerService{
		getServersFn: func(ctx context.Context, userCtx *UserContext, q GetServersQuery) (*PaginatedServers, error) {
			return &PaginatedServers{
				Servers: []model.Server{
					{Name: "s1", IPv4: "127.0.0.1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				},
			}, nil
		},
	}

	r := setupRouter(mock)

	req := httptest.NewRequest("GET", "/servers/export", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "s1")
}

// TODO: Move these tests to monitoring domain/feature when write unit tests there
// func TestHandlerSendReports_Success(t *testing.T) {
// 	mock := &MockServerService{
// 		reportFn: func(start, end time.Time, topN int, emails *[]string, ctx context.Context) (*Report, error) {
// 			return &Report{
// 				TotalServers: 10,
// 			}, nil
// 		},
// 	}

// 	r := setupRouter(mock)

// 	body := `{
// 		"count": 5,
// 		"emails": ["test@test.com"]
// 	}`

// 	req := httptest.NewRequest("POST", "/servers/report", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 200, w.Code)
// 	require.Contains(t, w.Body.String(), "10")
// }

// func TestHandlerSendReports_SuccessNoBody(t *testing.T) {
// 	mock := &MockServerService{
// 		reportFn: func(start, end time.Time, topN int, emails *[]string, ctx context.Context) (*Report, error) {
// 			return &Report{
// 				TotalServers: 10,
// 			}, nil
// 		},
// 	}

// 	r := setupRouter(mock)

// 	body := `{
// 	}`

// 	req := httptest.NewRequest("POST", "/servers/report", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 200, w.Code)
// 	require.Contains(t, w.Body.String(), "10")
// }

// func TestHandlerSendReports_InvalidJSON(t *testing.T) {
// 	mock := &MockServerService{}

// 	r := setupRouter(mock)

// 	req := httptest.NewRequest("POST", "/servers/report", strings.NewReader(`invalid`))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 400, w.Code)
// }

// func TestHandlerSendReports_InvalidStartDate(t *testing.T) {
// 	mock := &MockServerService{}

// 	r := setupRouter(mock)

// 	body := `{"start":"bad-date","count":1}`

// 	req := httptest.NewRequest("POST", "/servers/report", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 400, w.Code)
// }

// func TestHandlerSendReports_InvalidEndDate(t *testing.T) {
// 	mock := &MockServerService{}

// 	r := setupRouter(mock)

// 	body := `{"end":"bad-date","count":1}`

// 	req := httptest.NewRequest("POST", "/servers/report", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 400, w.Code)
// }

// func TestHandlerSendReports_InvalidTopN(t *testing.T) {
// 	mock := &MockServerService{}

// 	r := setupRouter(mock)

// 	body := `{"count":0}`

// 	req := httptest.NewRequest("POST", "/servers/report", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 400, w.Code)
// }

// func TestHandlerSendReports_ServiceError(t *testing.T) {
// 	mock := &MockServerService{
// 		reportFn: func(start, end time.Time, topN int, emails *[]string, ctx context.Context) (*Report, error) {
// 			return nil, errors.New("fail")
// 		},
// 	}

// 	r := setupRouter(mock)

// 	body := `{"count":1}`

// 	req := httptest.NewRequest("POST", "/servers/report", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	require.Equal(t, 500, w.Code)
// }
