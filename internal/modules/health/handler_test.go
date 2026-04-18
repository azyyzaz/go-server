package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-server/internal/middleware"
	"go-server/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	rdb "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestReadyReturnsDisabledWhenDependenciesAreNotConfigured(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.RequestID())
	NewHandler(nil, nil).Register(router.Group("/api/v1/health"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health/ready", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var resp response.Body
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	data, ok := resp.Data.(map[string]any)
	require.True(t, ok)
	require.Equal(t, dependencyHealthy, data["status"])

	dependencies, ok := data["dependencies"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, dependencySkipped, dependencies["database"])
	require.Equal(t, dependencySkipped, dependencies["redis"])
}

func TestReadyReturnsUnavailableWhenDependenciesFail(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.RequestID())
	NewHandler(newClosedGORM(t), newFailingRedisClient()).Register(router.Group("/api/v1/health"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health/ready", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)

	var resp response.Body
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	data, ok := resp.Data.(map[string]any)
	require.True(t, ok)
	require.Equal(t, dependencyFailed, data["status"])

	dependencies, ok := data["dependencies"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, dependencyFailed, dependencies["database"])
	require.Equal(t, dependencyFailed, dependencies["redis"])
}

func newClosedGORM(t *testing.T) *gorm.DB {
	t.Helper()

	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := gormDB.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
	return gormDB
}

func newFailingRedisClient() *rdb.Client {
	return rdb.NewClient(&rdb.Options{
		Addr:         "127.0.0.1:1",
		DialTimeout:  100 * time.Millisecond,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
	})
}
