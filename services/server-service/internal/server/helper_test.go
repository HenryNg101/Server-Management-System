package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
)

func newTestContext(query string) *gin.Context {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodGet, "/?"+query, nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return c
}

func TestParseServersQuery_Defaults(t *testing.T) {
	c := newTestContext("")

	q, err := ParseServersQuery(c)
	require.NoError(t, err)

	require.Nil(t, q.Status)
	require.Nil(t, q.Protocol)
	require.Nil(t, q.Name)

	require.Nil(t, q.Page)
	require.Nil(t, q.PageSize)

	require.Equal(t, "id", q.SortBy)
	require.Equal(t, "asc", q.Order)
}

func TestParseServersQuery_Filters(t *testing.T) {
	c := newTestContext("status=true&protocol=http&name=server1")

	q, err := ParseServersQuery(c)
	require.NoError(t, err)

	require.NotNil(t, q.Status)
	require.True(t, *q.Status)

	require.NotNil(t, q.Protocol)
	require.Equal(t, "http", *q.Protocol)

	require.NotNil(t, q.Name)
	require.Equal(t, "server1", *q.Name)
}

func TestParseServersQuery_InvalidStatus(t *testing.T) {
	c := newTestContext("status=notbool")

	_, err := ParseServersQuery(c)
	require.Error(t, err)
	require.Equal(t, "Invalid status", err.Error())
}

func TestParseServersQuery_PaginationValid(t *testing.T) {
	c := newTestContext("page=2&page_size=20")

	q, err := ParseServersQuery(c)
	require.NoError(t, err)

	require.NotNil(t, q.Page)
	require.NotNil(t, q.PageSize)

	require.Equal(t, 2, *q.Page)
	require.Equal(t, 20, *q.PageSize)
}

func TestParseServersQuery_InvalidPage(t *testing.T) {
	c := newTestContext("page=abc")

	_, err := ParseServersQuery(c)
	require.Error(t, err)
	require.Equal(t, "Invalid page", err.Error())
}

func TestParseServersQuery_InvalidPageSize(t *testing.T) {
	c := newTestContext("page_size=abc")

	_, err := ParseServersQuery(c)
	require.Error(t, err)
	require.Equal(t, "Invalid page size", err.Error())
}

func TestParseServersQuery_PaginationNormalization(t *testing.T) {
	c := newTestContext("page=0&page_size=999")

	q, err := ParseServersQuery(c)
	require.NoError(t, err)

	require.NotNil(t, q.Page)
	require.NotNil(t, q.PageSize)

	require.Equal(t, 1, *q.Page)
	require.Equal(t, 10, *q.PageSize)
}

func TestParseServersQuery_SortingDefaults(t *testing.T) {
	c := newTestContext("")

	q, err := ParseServersQuery(c)
	require.NoError(t, err)

	require.Equal(t, "id", q.SortBy)
	require.Equal(t, "asc", q.Order)
}

func TestParseServersQuery_CustomSorting(t *testing.T) {
	c := newTestContext("sort_by=name&order=desc")

	q, err := ParseServersQuery(c)
	require.NoError(t, err)

	require.Equal(t, "name", q.SortBy)
	require.Equal(t, "desc", q.Order)
}

// func TestParseServer_Valid(t *testing.T) {
// 	tests := []map[string]string{
// 		{"name": "srv", "status": "true", "ipv4_address": "127.0.0.1", "port": "80", "protocol": "tcp"},
// 		{"name": "server2", "status": "false", "ipv4_address": "8.8.8.8", "port": "80"},
// 	}

// 	for _, tc := range tests {
// 		_, err := parseServer(tc)
// 		if err != nil {
// 			t.Fatalf("unexpected error: %v", err)
// 		}
// 	}
// }

// func TestParseServer_InvalidCases(t *testing.T) {
// 	tests := []map[string]string{
// 		{"status": "true", "ipv4_address": "127.0.0.1", "port": "80"}, // missing name
// 		{"name": "a", "status": "bad", "ipv4_address": "127.0.0.1", "port": "80"},
// 		{"name": "a", "status": "true", "ipv4_address": "bad_ip", "port": "80"},
// 		{"name": "a", "status": "true", "ipv4_address": "127.0.0.1", "port": "99999"},
// 	}

// 	for _, tc := range tests {
// 		_, err := parseServer(tc)
// 		if err == nil {
// 			t.Fatalf("expected error for %+v", tc)
// 		}
// 	}
// }

// func TestMapRow(t *testing.T) {
// 	headers := []string{"a", "b"}
// 	row := []string{"1", "2"}

// 	m := mapRow(headers, row)

// 	if m["a"] != "1" || m["b"] != "2" {
// 		t.Fatalf("unexpected mapping")
// 	}
// }

// TODO: Move these tests to monitoring domain/feature when write unit tests there
// func TestBuildReportHTML(t *testing.T) {
// 	report := &Report{
// 		TotalServers: 2,
// 		ServersUp:    1,
// 		ServersDown:  1,
// 		Stats: map[uint]*ServerPullStats{
// 			1: {
// 				Uptime: 0.95,
// 			},
// 		},
// 	}

// 	html := buildReportHTML(report, time.Now(), time.Now())

// 	if len(html) == 0 {
// 		t.Fatalf("expected html output")
// 	}

// 	if !strings.Contains(html, "Server Report") {
// 		t.Fatalf("missing content")
// 	}
// }
