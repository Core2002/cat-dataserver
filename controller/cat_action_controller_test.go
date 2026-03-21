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

func setupCatActionController() *CatActionController {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")
	repo := repository.NewCatActionRepository()
	return NewCatActionController(repo)
}

func TestCreateCatAction(t *testing.T) {
	ctrl := setupCatActionController()

	newAction := model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}

	body, _ := json.Marshal(newAction)
	req, _ := http.NewRequest("POST", "/cat-actions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatAction(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.CatAction
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.ActionType != model.CatActionFeed {
		t.Errorf("Expected action type '%s', got '%s'", model.CatActionFeed, response.ActionType)
	}
}

func TestGetCatActions(t *testing.T) {
	ctrl := setupCatActionController()

	req, _ := http.NewRequest("GET", "/cat-actions", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetCatActions(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []model.CatAction
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
	}
}

func TestGetCatActionsPage(t *testing.T) {
	ctrl := setupCatActionController()

	req, _ := http.NewRequest("GET", "/cat-actions/page?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetCatActionsPage(c)

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

func TestGetCatAction(t *testing.T) {
	ctrl := setupCatActionController()

	req, _ := http.NewRequest("GET", "/cat-actions/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.GetCatAction(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestGetCatActionsByCatID(t *testing.T) {
	ctrl := setupCatActionController()

	req, _ := http.NewRequest("GET", "/cat-actions/cat/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

	ctrl.GetCatActionsByCatID(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetCatActionsBySiteID(t *testing.T) {
	ctrl := setupCatActionController()

	req, _ := http.NewRequest("GET", "/cat-actions/site/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "site_id", Value: "1"}}

	ctrl.GetCatActionsBySiteID(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetCatActionsByUserID(t *testing.T) {
	ctrl := setupCatActionController()

	req, _ := http.NewRequest("GET", "/cat-actions/user/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "user_id", Value: "1"}}

	ctrl.GetCatActionsByUserID(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateCatAction(t *testing.T) {
	ctrl := setupCatActionController()

	updates := model.CatAction{
		ActionDetail: "更新的操作详情",
	}

	body, _ := json.Marshal(updates)
	req, _ := http.NewRequest("PUT", "/cat-actions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.UpdateCatAction(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestDeleteCatAction(t *testing.T) {
	ctrl := setupCatActionController()

	req, _ := http.NewRequest("DELETE", "/cat-actions/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.DeleteCatAction(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "CatAction deleted successfully" {
		t.Errorf("Expected success message, got '%s'", response["message"])
	}
}

func TestNewCatActionController(t *testing.T) {
	repo := repository.NewCatActionRepository()
	ctrl := NewCatActionController(repo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
}
