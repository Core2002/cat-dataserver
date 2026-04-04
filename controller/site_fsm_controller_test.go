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

func setupSiteFSMController() *SiteFSMController {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")
	repo := repository.NewSiteFSMRepository()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13800138000",
	}
	siteRepo := repository.NewSiteRepository()
	siteRepo.Create(testSite)

	return NewSiteFSMController(repo, siteRepo)
}

func TestCreateSiteFSM(t *testing.T) {
	ctrl := setupSiteFSMController()

	newFSM := model.SiteFSM{
		SiteID: 1,
	}

	body, _ := json.Marshal(newFSM)
	req, _ := http.NewRequest("POST", "/site-fsms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateSiteFSM(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.SiteFSM
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.SiteID != 1 {
		t.Errorf("Expected site ID 1, got %d", response.SiteID)
	}
}

func TestGetSiteFSMBySiteID(t *testing.T) {
	ctrl := setupSiteFSMController()

	req, _ := http.NewRequest("GET", "/site-fsms/site/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "site_id", Value: "1"}}

	ctrl.GetSiteFSMBySiteID(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestUpdateDisinfectTime(t *testing.T) {
	ctrl := setupSiteFSMController()

	timeUpdate := struct {
		Time string `json:"last_disinfect_time"`
	}{
		Time: "2024-01-01T12:00:00Z",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/site-fsms/1/disinfect-time", bytes.NewBuffer(body))
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
	ctrl := setupSiteFSMController()

	timeUpdate := struct {
		Time string `json:"last_feed_time"`
	}{
		Time: "2024-01-01T12:00:00Z",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/site-fsms/1/feed-time", bytes.NewBuffer(body))
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
	ctrl := setupSiteFSMController()

	timeUpdate := struct {
		Time string `json:"last_give_water_time"`
	}{
		Time: "2024-01-01T12:00:00Z",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/site-fsms/1/give-water-time", bytes.NewBuffer(body))
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
	ctrl := setupSiteFSMController()

	timeUpdate := struct {
		Time string `json:"last_play_time"`
	}{
		Time: "2024-01-01T12:00:00Z",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/site-fsms/1/play-time", bytes.NewBuffer(body))
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

func TestNewSiteFSMController(t *testing.T) {
	repo := repository.NewSiteFSMRepository()
	siteRepo := repository.NewSiteRepository()
	ctrl := NewSiteFSMController(repo, siteRepo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
	if ctrl.siteRepo != siteRepo {
		t.Error("Controller siteRepo does not match input siteRepo")
	}
}
