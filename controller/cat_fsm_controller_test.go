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

func setupCatFSMController() *CatFSMController {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")
	repo := repository.NewCatFSMRepository()
	return NewCatFSMController(repo)
}

func TestCreateCatFSM(t *testing.T) {
	ctrl := setupCatFSMController()

	newFSM := model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}

	body, _ := json.Marshal(newFSM)
	req, _ := http.NewRequest("POST", "/cat-fsms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatFSM(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.CatFSM
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.CatID != 1 {
		t.Errorf("Expected cat_id 1, got %d", response.CatID)
	}
}

func TestGetCatFSMsPage(t *testing.T) {
	ctrl := setupCatFSMController()

	req, _ := http.NewRequest("GET", "/cat-fsms/page?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetCatFSMsPage(c)

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

func TestGetCatFSM(t *testing.T) {
	ctrl := setupCatFSMController()

	req, _ := http.NewRequest("GET", "/cat-fsms/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.GetCatFSM(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestGetCatFSMsBySiteID(t *testing.T) {
	ctrl := setupCatFSMController()

	req, _ := http.NewRequest("GET", "/cat-fsms/site/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "site_id", Value: "1"}}

	ctrl.GetCatFSMsBySiteID(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateCatFSM(t *testing.T) {
	ctrl := setupCatFSMController()

	updates := model.CatFSM{
		WeightKG: 5.0,
	}

	body, _ := json.Marshal(updates)
	req, _ := http.NewRequest("PUT", "/cat-fsms/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.UpdateCatFSM(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestUpdateTemperature(t *testing.T) {
	ctrl := setupCatFSMController()

	temperatureUpdate := struct {
		Temperature float32 `json:"temperature"`
	}{
		Temperature: 39.0,
	}

	body, _ := json.Marshal(temperatureUpdate)
	req, _ := http.NewRequest("PATCH", "/cat-fsms/1/temperature", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

	ctrl.UpdateTemperature(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateWeight(t *testing.T) {
	ctrl := setupCatFSMController()

	weightUpdate := struct {
		Weight float32 `json:"weight"`
	}{
		Weight: 5.0,
	}

	body, _ := json.Marshal(weightUpdate)
	req, _ := http.NewRequest("PATCH", "/cat-fsms/1/weight", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

	ctrl.UpdateWeight(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateTrimNailsTime(t *testing.T) {
	ctrl := setupCatFSMController()

	timeUpdate := struct {
		Time string `json:"time"`
	}{
		Time: "2024-01-01 12:00:00",
	}

	body, _ := json.Marshal(timeUpdate)
	req, _ := http.NewRequest("PATCH", "/cat-fsms/1/trim-nails-time", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

	ctrl.UpdateTrimNailsTime(c)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

func TestDeleteCatFSM(t *testing.T) {
	ctrl := setupCatFSMController()

	req, _ := http.NewRequest("DELETE", "/cat-fsms/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	ctrl.DeleteCatFSM(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "CatFSM deleted successfully" {
		t.Errorf("Expected success message, got '%s'", response["message"])
	}
}

func TestNewCatFSMController(t *testing.T) {
	repo := repository.NewCatFSMRepository()
	ctrl := NewCatFSMController(repo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
}
