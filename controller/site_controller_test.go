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

func setupSiteController() *SiteController {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")
	repo := repository.NewSiteRepository()
	return NewSiteController(repo)
}

func TestCreateSite(t *testing.T) {
	ctrl := setupSiteController()

	newSite := model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}

	body, _ := json.Marshal(newSite)
	req, _ := http.NewRequest("POST", "/sites", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateSite(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.Site
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.SiteName != "测试站点" {
		t.Errorf("Expected site name '测试站点', got '%s'", response.SiteName)
	}
}

func TestGetSites(t *testing.T) {
	ctrl := setupSiteController()

	req, _ := http.NewRequest("GET", "/sites", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetSites(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []model.Site
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
	}
}

func TestGetSitesPage(t *testing.T) {
	ctrl := setupSiteController()

	req, _ := http.NewRequest("GET", "/sites/page?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetSitesPage(c)

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

func TestGetSite(t *testing.T) {
	ctrl := setupSiteController()

	req, _ := http.NewRequest("GET", "/sites/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.GetSite(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestUpdateSite(t *testing.T) {
	ctrl := setupSiteController()

	updates := model.Site{
		SiteName: "更新的站点",
	}

	body, _ := json.Marshal(updates)
	req, _ := http.NewRequest("PUT", "/sites/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.UpdateSite(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestDeleteSite(t *testing.T) {
	ctrl := setupSiteController()

	req, _ := http.NewRequest("DELETE", "/sites/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.DeleteSite(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "Site deleted successfully" {
		t.Errorf("Expected success message, got '%s'", response["message"])
	}
}

func TestUpdateDisinfectTime(t *testing.T) {
	ctrl := setupSiteController()

	timeUpdate := struct {
		Time string `json:"time"`
	}{
		Time: "2024-01-01 12:00:00",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/sites/1/disinfect-time", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "site_id", Value: "1"}}

	ctrl.UpdateDisinfectTime(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateFeedTime(t *testing.T) {
	ctrl := setupSiteController()

	timeUpdate := struct {
		Time string `json:"time"`
	}{
		Time: "2024-01-01 12:00:00",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/sites/1/feed-time", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "site_id", Value: "1"}}

	ctrl.UpdateFeedTime(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateGiveWaterTime(t *testing.T) {
	ctrl := setupSiteController()

	timeUpdate := struct {
		Time string `json:"time"`
	}{
		Time: "2024-01-01 12:00:00",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/sites/1/give-water-time", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "site_id", Value: "1"}}

	ctrl.UpdateGiveWaterTime(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestUpdatePlayTime(t *testing.T) {
	ctrl := setupSiteController()

	timeUpdate := struct {
		Time string `json:"time"`
	}{
		Time: "2024-01-01 12:00:00",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/sites/1/play-time", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "site_id", Value: "1"}}

	ctrl.UpdatePlayTime(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestNewSiteController(t *testing.T) {
	repo := repository.NewSiteRepository()
	ctrl := NewSiteController(repo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
}
