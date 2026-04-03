package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/middleware"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func setupCatEventController() *CatEventController {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")

	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		middleware.RegisterCustomValidators(v)
	}

	repo := repository.NewCatEventRepository()
	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13800138000",
	}
	siteRepo.Create(testSite)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             1,
		CatName:           "测试猫",
		CatPhotoUri:       "http://example.com/photo.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人",
		MasterPhoneNumber: "13900139000",
	}
	catRepo.Create(testCat)

	return NewCatEventController(repo, catRepo, siteRepo)
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
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatEvent(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d, response: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var responseData map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// 检查 event 字段
	eventData, ok := responseData["event"].(map[string]interface{})
	if !ok {
		t.Fatalf("Response should contain 'event' field")
	}

	// 验证事件类型（将 float64 或 string 转换）
	var eventType string
	if eventTypeFloat, ok := eventData["event_type"].(float64); ok {
		eventType = string(rune(int(eventTypeFloat)))
	} else if eventTypeStr, ok := eventData["event_type"].(string); ok {
		eventType = eventTypeStr
	}

	if eventType != string(model.CatSick) {
		t.Errorf("Expected event type '%s', got '%s'", model.CatSick, eventType)
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
	c.Params = []gin.Param{{Key: "event_id", Value: "1"}}

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
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

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
	c.Params = []gin.Param{{Key: "cat_id", Value: "1"}}

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
	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	ctrl := NewCatEventController(repo, catRepo, siteRepo)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
}

// 测试创建 CatEvent 时使用无效的 CatID
func TestCreateCatEventWithInvalidCatID(t *testing.T) {
	ctrl := setupCatEventController()

	newEvent := model.CatEvent{
		EventID:   1,
		CatID:     999, // 不存在的 CatID
		SiteID:    1,
		EventType: model.CatSick,
		Detail:    "测试无效的 CatID",
	}

	body, _ := json.Marshal(newEvent)
	req, _ := http.NewRequest("POST", "/cat-events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatEvent(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	errors, ok := response["errors"].([]interface{})
	if !ok || len(errors) != 1 {
		t.Fatalf("Expected errors array with 1 element, got: %v", response["errors"])
	}

	if errors[0] != "CatID does not exist" {
		t.Errorf("Expected error message 'CatID does not exist', got '%s'", errors[0])
	}
}

// 测试创建 CatEvent 时使用无效的 SiteID
func TestCreateCatEventWithInvalidSiteID(t *testing.T) {
	ctrl := setupCatEventController()

	newEvent := model.CatEvent{
		EventID:   1,
		CatID:     1,   // 有效的 CatID
		SiteID:    999, // 不存在的 SiteID
		EventType: model.CatSick,
		Detail:    "测试无效的 SiteID",
	}

	body, _ := json.Marshal(newEvent)
	req, _ := http.NewRequest("POST", "/cat-events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatEvent(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	errors, ok := response["errors"].([]interface{})
	if !ok || len(errors) != 1 {
		t.Fatalf("Expected errors array with 1 element, got: %v", response["errors"])
	}

	if errors[0] != "SiteID does not exist" {
		t.Errorf("Expected error message 'SiteID does not exist', got '%s'", errors[0])
	}
}
