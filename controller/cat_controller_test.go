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

func setupCatController() *CatController {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")
	repo := repository.NewCatRepository()
	fsmRepo := repository.NewCatFSMRepository()
	siteRepo := repository.NewSiteRepository()
	return NewCatController(repo, fsmRepo, siteRepo)
}

func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	HealthCheck(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestCreateCat(t *testing.T) {
	ctrl := setupCatController()

	// 先创建一个站点
	site := model.Site{
		SiteName:    "测试站点",
		SiteAddress: "测试地址",
	}
	ctrl.siteRepo.Create(&site)

	// 使用 CreateCatRequest 的字段名
	requestBody := map[string]interface{}{
		"cat_name":            "测试猫",
		"cat_photo_uri":       "http://example.com/cat.jpg",
		"cat_type":            "英国短毛猫",
		"cat_gender":          "公",
		"master_name":         "张三",
		"master_phone_number": "13800138000",
		"site_id":             site.SiteID,
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/cats", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCat(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d, body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	cat, ok := response["cat"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected cat in response")
	}

	if cat["cat_name"] != "测试猫" {
		t.Errorf("Expected cat name '测试猫', got '%s'", cat["cat_name"])
	}
}

func TestGetCatPage(t *testing.T) {
	ctrl := setupCatController()

	req, _ := http.NewRequest("GET", "/cats/page?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.GetCatsPage(c)

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
	if response.PageSize != 10 {
		t.Errorf("Expected page_size 10, got %d", response.PageSize)
	}
}

func TestGetCat(t *testing.T) {
	ctrl := setupCatController()

	req, _ := http.NewRequest("GET", "/cats/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

	ctrl.GetCat(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestUpdateCat(t *testing.T) {
	ctrl := setupCatController()

	updates := model.Cat{
		CatName: "更新的猫",
	}

	body, _ := json.Marshal(updates)
	req, _ := http.NewRequest("PUT", "/cats/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

	ctrl.UpdateCat(c)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNotFound, w.Code)
	}
}

func TestDeleteCat(t *testing.T) {
	ctrl := setupCatController()

	req, _ := http.NewRequest("DELETE", "/cats/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

	ctrl.DeleteCat(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "Cat deleted successfully" {
		t.Errorf("Expected success message, got '%s'", response["message"])
	}
}

func TestNewCatController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")
	repo := repository.NewCatRepository()
	fsmRepo := repository.NewCatFSMRepository()
	siteRepo := repository.NewSiteRepository()
	ctrl := NewCatController(repo, fsmRepo, siteRepo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
	if ctrl.fsmRepo != fsmRepo {
		t.Error("Controller fsmRepo does not match input fsmRepo")
	}
	if ctrl.siteRepo != siteRepo {
		t.Error("Controller siteRepo does not match input siteRepo")
	}
}
