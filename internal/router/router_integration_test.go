package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "go-server/docs"
	"go-server/internal/config"
	appjwt "go-server/internal/jwt"
	filemodule "go-server/internal/modules/file"
	"go-server/internal/response"
	"go-server/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRouterSupportsCORSAndSwagger(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	require.NoError(t, validation.Init())

	engine := newIntegrationRouter(t, func(cfg *config.Config) {})

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, "http://localhost:3000", rec.Header().Get("Access-Control-Allow-Origin"))
	require.NotEmpty(t, rec.Header().Get("X-Trace-ID"))

	swaggerReq := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	swaggerRec := httptest.NewRecorder()
	engine.ServeHTTP(swaggerRec, swaggerReq)

	require.Equal(t, http.StatusOK, swaggerRec.Code)
}

func TestRouterProfileValidationReturnsChineseMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	require.NoError(t, validation.Init())

	engine := newIntegrationRouter(t, func(cfg *config.Config) {
		cfg.Rate.Enabled = false
	})

	jwtManager := appjwt.NewManager("test-secret", time.Hour, 24*time.Hour)
	token, err := jwtManager.GenerateAccessToken(1, "admin")
	require.NoError(t, err)

	body, err := json.Marshal(map[string]string{
		"name":  "A",
		"email": "bad-email",
		"phone": "13800138000",
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp response.Body
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "INVALID_ARGUMENT", resp.Code)
	require.Contains(t, resp.Message, "email")
	require.NotEqual(t, "参数错误", resp.Message)
	require.NotEmpty(t, resp.RequestID)
	require.Equal(t, resp.RequestID, resp.TraceID)
}

func TestRouterRateLimitProtectsAuthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	require.NoError(t, validation.Init())

	engine := newIntegrationRouter(t, func(cfg *config.Config) {
		cfg.Rate.Enabled = true
		cfg.Rate.RequestsPerSecond = 1
		cfg.Rate.MaxDelay = 50 * time.Millisecond
		cfg.Rate.ProtectedPrefixes = []string{"/api/v1/auth/captcha"}
	})

	firstReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/captcha", nil)
	firstReq.RemoteAddr = "127.0.0.1:12345"
	firstRec := httptest.NewRecorder()
	engine.ServeHTTP(firstRec, firstReq)
	require.Equal(t, http.StatusOK, firstRec.Code)

	secondReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/captcha", nil)
	secondReq.RemoteAddr = "127.0.0.1:12345"
	secondRec := httptest.NewRecorder()
	engine.ServeHTTP(secondRec, secondReq)
	require.Equal(t, http.StatusTooManyRequests, secondRec.Code)
}

func newIntegrationRouter(t *testing.T, mutate func(*config.Config)) *gin.Engine {
	t.Helper()

	cfg := config.Config{
		AppName: "go-server-test",
		Env:     "test",
		HTTP: config.HTTPConfig{
			Host: "127.0.0.1",
			Port: "8080",
		},
		JWT: config.JWTConfig{
			Secret:          "test-secret",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: 24 * time.Hour,
		},
		Audit: config.AuditConfig{
			Enabled:        false,
			RegionFallback: "test",
		},
		File: config.FileConfig{
			Storage:       "local",
			MaxSizeMB:     10,
			AllowedExts:   []string{".png", ".jpg", ".pdf"},
			AvatarMaxSize: 2,
			Local: config.LocalFileConfig{
				BaseDir:    t.TempDir(),
				BaseURL:    "/uploads",
				AvatarDir:  "avatars",
				DateLayout: "2006/01/02",
			},
		},
		CORS: config.CORSConfig{
			Enabled:       true,
			AllowOrigins:  []string{"http://localhost:3000"},
			AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:  []string{"Authorization", "Content-Type", "X-Request-ID", "X-Trace-ID"},
			ExposeHeaders: []string{"X-Request-ID", "X-Trace-ID"},
			MaxAge:        12 * time.Hour,
		},
		Rate: config.RateLimitConfig{
			Enabled:           false,
			RequestsPerSecond: 5,
			MaxDelay:          0,
			ProtectedPrefixes: []string{"/api/v1/auth/captcha"},
		},
	}

	if mutate != nil {
		mutate(&cfg)
	}

	storage, err := filemodule.NewStorage(cfg.File)
	require.NoError(t, err)

	jwtManager := appjwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	return New(cfg, nil, nil, jwtManager, nil, nil, nil, storage)
}
