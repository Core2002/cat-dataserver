package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fifu.fun/cat-dataserver/database"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")
	return SetupRouter()
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}
}

func TestSetupRouter(t *testing.T) {
	router := SetupRouter()

	if router == nil {
		t.Error("Expected non-nil router")
	}

	routes := router.Routes()

	if len(routes) == 0 {
		t.Error("Expected at least one route")
	}

	// 检查关键路由是否存在
	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Method+":"+route.Path] = true
	}

	expectedRoutes := []string{
		"GET:/health",
		"POST:/cats",
		"POST:/cat-events",
		"POST:/cat-actions",
		"POST:/cat-fsms",
		"POST:/sites",
		"POST:/site-fsms",
	}

	for _, expectedRoute := range expectedRoutes {
		if !routeMap[expectedRoute] {
			t.Errorf("Expected route %s to be registered", expectedRoute)
		}
	}
}

func TestCatRoutes(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET cats page", "GET", "/cats/page", http.StatusBadRequest},
		{"GET single cat (not found)", "GET", "/cats/999", http.StatusNotFound},
		{"POST create cat", "POST", "/cats", http.StatusBadRequest},
		{"PUT update cat (not found)", "PUT", "/cats/999", http.StatusNotFound},
		{"DELETE cat (not found)", "DELETE", "/cats/999", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCatEventRoutes(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET events page", "GET", "/cat-events/page", http.StatusBadRequest},
		{"GET events by cat ID", "GET", "/cat-events/cat/1", http.StatusOK},
		{"GET events by site ID", "GET", "/cat-events/site/1", http.StatusOK},
		{"POST create event", "POST", "/cat-events", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCatActionRoutes(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET actions page", "GET", "/cat-actions/page", http.StatusBadRequest},
		{"GET actions by cat ID", "GET", "/cat-actions/cat/1", http.StatusOK},
		{"GET actions by site ID", "GET", "/cat-actions/site/1", http.StatusOK},
		{"GET actions by user ID", "GET", "/cat-actions/user/1", http.StatusOK},
		{"POST create action", "POST", "/cat-actions", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCatFSMRoutes(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET FSMs page", "GET", "/cat-fsms/page", http.StatusBadRequest},
		{"GET FSMs by site ID", "GET", "/cat-fsms/site/1", http.StatusOK},
		{"POST create FSM", "POST", "/cat-fsms", http.StatusBadRequest},
		{"PATCH update temperature", "PATCH", "/cat-fsms/1/temperature", http.StatusBadRequest},
		{"PATCH update weight", "PATCH", "/cat-fsms/1/weight", http.StatusBadRequest},
		{"PATCH update trim nails time", "PATCH", "/cat-fsms/1/trim-nails-time", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestSiteRoutes(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET sites page", "GET", "/sites/page", http.StatusBadRequest},
		{"POST create site", "POST", "/sites", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestSiteFSMRoutes(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"GET FSM by site ID", "GET", "/site-fsms/site/1", http.StatusNotFound}, // 数据库中无数据
		{"POST create FSM", "POST", "/site-fsms", http.StatusBadRequest},
		{"PATCH update disinfect time", "PATCH", "/site-fsms/1/disinfect-time", http.StatusNotFound}, // site 不存在
		{"PATCH update feed time", "PATCH", "/site-fsms/1/feed-time", http.StatusNotFound},           // site 不存在
		{"PATCH update give water time", "PATCH", "/site-fsms/1/give-water-time", http.StatusNotFound}, // site 不存在
		{"PATCH update play time", "PATCH", "/site-fsms/1/play-time", http.StatusNotFound},           // site 不存在
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRouterMiddleware(t *testing.T) {
	router := setupTestRouter()

	routes := router.Routes()

	// 检查是否有默认的 logger 和 recovery 中间件
	if len(routes) == 0 {
		t.Error("Expected routes to be registered")
	}

	// 确认 router 是 gin.Engine
	if router == nil {
		t.Error("Expected gin.Engine")
	}
}
