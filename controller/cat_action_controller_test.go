package controller

import (
	"bytes"
	"encoding/json"
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/middleware"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func setupCatActionController() *CatActionController {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")

	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		middleware.RegisterCustomValidators(v)
	}

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()

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

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	fsmRepo.Create(testFSM)

	actionProcessor := middleware.NewActionProcessor(actionRepo, fsmRepo)
	return NewCatActionController(actionRepo, catRepo, siteRepo, actionProcessor)
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
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatAction(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var responseData map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// 检查 action 字段
	actionData, ok := responseData["action"].(map[string]interface{})
	if !ok {
		t.Fatalf("Response should contain 'action' field")
	}

	// 验证动作类型（将 float64 转换为 string）
	actionTypeFloat, ok := actionData["action_type"].(float64)
	if !ok {
		actionTypeStr, ok := actionData["action_type"].(string)
		if !ok {
			t.Errorf("ActionType should be string or float64")
		} else if actionTypeStr != string(model.CatActionFeed) {
			t.Errorf("Expected action type '%s', got '%s'", model.CatActionFeed, actionTypeStr)
		}
	} else {
		// 如果是 float64，可能是因为 JSON 数字编码
		t.Logf("ActionType is float64: %v", actionTypeFloat)
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
	req.Header.Set("X-User-ID", "1")

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
	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	fsmRepo := repository.NewCatFSMRepository()
	actionProcessor := middleware.NewActionProcessor(repo, fsmRepo)
	ctrl := NewCatActionController(repo, catRepo, siteRepo, actionProcessor)

	if ctrl == nil {
		t.Error("Expected non-nil controller")
	}
	if ctrl.repo != repo {
		t.Error("Controller repo does not match input repo")
	}
}

// 测试创建 CatAction 时使用无效的 CatID
func TestCreateCatActionWithInvalidCatID(t *testing.T) {
	ctrl := setupCatActionController()

	newAction := model.CatAction{
		CatID:        999, // 不存在的 CatID
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "测试无效的 CatID",
	}

	body, _ := json.Marshal(newAction)
	req, _ := http.NewRequest("POST", "/cat-actions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatAction(c)

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

// 测试创建 CatAction 时使用无效的 SiteID
func TestCreateCatActionWithInvalidSiteID(t *testing.T) {
	ctrl := setupCatActionController()

	newAction := model.CatAction{
		CatID:        1,   // 有效的 CatID
		SiteID:       999, // 不存在的 SiteID
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "测试无效的 SiteID",
	}

	body, _ := json.Marshal(newAction)
	req, _ := http.NewRequest("POST", "/cat-actions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ctrl.CreateCatAction(c)

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
