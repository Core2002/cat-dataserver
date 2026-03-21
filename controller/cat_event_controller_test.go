package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
)

func setupCatEventController() *CatEventController {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")
	repo := repository.NewCatEventRepository()
	return NewCatEventController(repo)
}

func TestCreateCatEvent(t *testing.T) {
	ctrl := setupCatEventController()

	newEvent := model.CatEvent{
		EventID:   1,
		EventType: model.CatSick,
		SiteID:    1,
		CatID:     1,
		Detail:    "测试事件详情",
	}

	body, _ := json.Marshal(newEvent)
	req, _ := http.NewRequest("POST", "/cat-events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatEvent(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.CatEvent
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.EventType != model.CatSick {
		t.Errorf("Expected event type '%s', got '%s'", model.CatSick, response.EventType)
	}
}

func TestGetCatEvents(t *testing.T) {
	ctrl := setupCatEventController()

	req, _ := http.NewRequest("GET", "/cat-events", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetCatEvents(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []model.CatEvent
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
	}
}

func TestGetCatEventsPage(t *testing.T) {
	ctrl := setupCatEventController()

	req, _ := http.NewRequest("GET", "/cat-events/page?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetCatEventsPage(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response model.PaginationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Page != 1 {
		t.Errorf("Expected page 1, got %d", response.Page)
	}
}

func TestGetCatEvent(t *testing.T) {
	ctrl := setupCatEventController()

	req, _ := http.NewRequest("GET", "/cat-events/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.GetCatEvent(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestGetCatEventsByCatID(t *testing.T) {
	ctrl := setupCatEventController()

	req, _ := http.NewRequest("GET", "/cat-events/cat/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

	ctrl.GetCatEventsByCatID(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetCatEventsBySiteID(t *testing.T) {
	ctrl := setupCatEventController()

	req, _ := http.NewRequest("GET", "/cat-events/site/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "site_id", Value: "1"}}

	ctrl.GetCatEventsBySiteID(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateCatEvent(t *testing.T) {
	ctrl := setupCatEventController()

	updates := model.CatEvent{
		Detail: "更新的事件详情",
	}

	body, _ := json.Marshal(updates)
	req, _ := http.NewRequest("PUT", "/cat-events/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.UpdateCatEvent(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestDeleteCatEvent(t *testing.T) {
	ctrl := setupCatEventController()

	req, _ := http.NewRequest("DELETE", "/cat-events/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.DeleteCatEvent(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "CatEvent deleted successfully" {
		t.Errorf("Expected success message, got '%s'", response["message"])
	}
}

func TestNewCatEventController(t *testing.T) {
	repo := repository.NewCatEventRepository()
	ctrl := NewCatEventController(repo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
}
