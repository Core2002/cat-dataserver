package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/middleware"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"github.com/gin-gonic/gin"
)

func setupHospitalizationController(t *testing.T) (*HospitalizationController, *repository.CatActionRepository, *repository.CatEventRepository, *repository.CatFSMRepository, uint, uint) {
	gin.SetMode(gin.TestMode)
	err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("init db failed: %v", err)
	}

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	catFSMRepo := repository.NewCatFSMRepository()
	catEventRepo := repository.NewCatEventRepository()
	catActionRepo := repository.NewCatActionRepository()
	actionProcessor := middleware.NewActionProcessor(catActionRepo, catFSMRepo)

	site := &model.Site{
		SiteName:             "医院B",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13800138000",
	}
	if err := siteRepo.Create(site); err != nil {
		t.Fatalf("create site failed: %v", err)
	}

	cat := &model.Cat{
		CatName:           "流程测试猫",
		CatPhotoUri:       "http://example.com/cat.jpg",
		CatType:           "英短",
		CatGender:         "母",
		MasterName:        "李四",
		MasterPhoneNumber: "13900139000",
	}
	if err := catRepo.Create(cat); err != nil {
		t.Fatalf("create cat failed: %v", err)
	}

	fsm := &model.CatFSM{
		CatID:         cat.CatID,
		SiteID:        0,
		TemperatureC:  37.5,
		WeightKG:      4.0,
		TrimNailsTime: time.Now(),
	}
	if err := catFSMRepo.Create(fsm); err != nil {
		t.Fatalf("create fsm failed: %v", err)
	}

	ctrl := NewHospitalizationController(catRepo, siteRepo, catFSMRepo, catEventRepo, actionProcessor)
	return ctrl, catActionRepo, catEventRepo, catFSMRepo, cat.CatID, site.SiteID
}

func TestHospitalizationFlow_WritesActionAndEventAndUpdatesFSM(t *testing.T) {
	ctrl, actionRepo, eventRepo, fsmRepo, catID, siteID := setupHospitalizationController(t)

	admitBody := map[string]interface{}{
		"cat_id":                catID,
		"site_id":               siteID,
		"user_id":               1,
		"admission_reason":      "术后观察",
		"admission_note":        "需监测",
		"initial_temperature_c": 38.7,
		"initial_weight_kg":     4.5,
	}
	body, _ := json.Marshal(admitBody)
	req, _ := http.NewRequest("POST", "/hospitalizations/admit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	ctrl.AdmitCat(c)
	if w.Code != http.StatusCreated {
		t.Fatalf("admit status expected %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	dischargeBody := map[string]interface{}{
		"cat_id":               catID,
		"user_id":              1,
		"discharge_reason":     "恢复良好",
		"discharge_note":       "继续随访",
		"final_temperature_c":  38.1,
		"final_weight_kg":      4.7,
	}
	body, _ = json.Marshal(dischargeBody)
	req, _ = http.NewRequest("POST", "/hospitalizations/discharge", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	ctrl.DischargeCat(c)
	if w.Code != http.StatusOK {
		t.Fatalf("discharge status expected %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	actions, err := actionRepo.FindByCatID(catID)
	if err != nil {
		t.Fatalf("query actions failed: %v", err)
	}
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(actions))
	}

	events, err := eventRepo.FindByCatID(catID)
	if err != nil {
		t.Fatalf("query events failed: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	fsm, err := fsmRepo.FindByID(catID)
	if err != nil {
		t.Fatalf("query fsm failed: %v", err)
	}
	if fsm.SiteID != 0 {
		t.Fatalf("expected cat site_id=0 after discharge, got %d", fsm.SiteID)
	}
}

func TestHospitalizationFlow_RejectDuplicateAdmit(t *testing.T) {
	ctrl, _, _, _, catID, siteID := setupHospitalizationController(t)

	admitBody := map[string]interface{}{
		"cat_id":                catID,
		"site_id":               siteID,
		"user_id":               1,
		"admission_reason":      "首次入院",
		"admission_note":        "观察",
		"initial_temperature_c": 38.6,
		"initial_weight_kg":     4.4,
	}
	body, _ := json.Marshal(admitBody)
	req, _ := http.NewRequest("POST", "/hospitalizations/admit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	ctrl.AdmitCat(c)
	if w.Code != http.StatusCreated {
		t.Fatalf("first admit expected %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	// 重复入院应被拒绝
	body, _ = json.Marshal(admitBody)
	req, _ = http.NewRequest("POST", "/hospitalizations/admit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	ctrl.AdmitCat(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("duplicate admit expected %d, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestHospitalizationFlow_RejectDischargeWhenNotAdmitted(t *testing.T) {
	ctrl, _, _, _, catID, _ := setupHospitalizationController(t)

	dischargeBody := map[string]interface{}{
		"cat_id":               catID,
		"user_id":              1,
		"discharge_reason":     "误操作测试",
		"discharge_note":       "未入院出院",
		"final_temperature_c":  38.0,
		"final_weight_kg":      4.3,
	}
	body, _ := json.Marshal(dischargeBody)
	req, _ := http.NewRequest("POST", "/hospitalizations/discharge", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	ctrl.DischargeCat(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("discharge without admit expected %d, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}
