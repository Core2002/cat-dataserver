package middleware

import (
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestActionProcessor_ProcessAction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")

	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()
	processor := NewActionProcessor(actionRepo, fsmRepo)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	if err := fsmRepo.Create(testFSM); err != nil {
		t.Fatalf("Failed to create test FSM: %v", err)
	}

	// 测试测体温动作
	temperatureAction := &model.CatAction{
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionTakeTemperature,
		ActionDetail: "39.5",
	}

	updatedFSM, err := processor.ProcessAction(temperatureAction)
	if err != nil {
		t.Errorf("ProcessAction failed: %v", err)
	}

	// 验证状态机是否更新
	freshFSM, err := fsmRepo.FindByID(1)
	if err != nil {
		t.Errorf("Failed to fetch FSM: %v", err)
	}

	if freshFSM.TemperatureC != 39.5 {
		t.Errorf("Expected temperature 39.5, got %v", freshFSM.TemperatureC)
	}

	if updatedFSM != nil && updatedFSM.TemperatureC != 39.5 {
		t.Errorf("Updated FSM temperature should be 39.5, got %v", updatedFSM.TemperatureC)
	}
}

func TestActionProcessor_ProcessWeightAction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")

	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()
	processor := NewActionProcessor(actionRepo, fsmRepo)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	if err := fsmRepo.Create(testFSM); err != nil {
		t.Fatalf("Failed to create test FSM: %v", err)
	}

	// 测试体检动作（更新体重）
	healthCheckAction := &model.CatAction{
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionHealthCheck,
		ActionDetail: "5.2",
	}

	_, err := processor.ProcessAction(healthCheckAction)
	if err != nil {
		t.Errorf("ProcessAction failed: %v", err)
	}

	// 验证状态机是否更新
	freshFSM, err := fsmRepo.FindByID(1)
	if err != nil {
		t.Errorf("Failed to fetch FSM: %v", err)
	}

	if freshFSM.WeightKG != 5.2 {
		t.Errorf("Expected weight 5.2, got %v", freshFSM.WeightKG)
	}
}

func TestActionProcessor_ProcessTrimNailsAction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.InitDB(":memory:")

	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()
	processor := NewActionProcessor(actionRepo, fsmRepo)

	// 创建测试用的 FSM 记录
	oldTime := time.Now().Add(-24 * time.Hour)
	testFSM := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: oldTime,
	}
	if err := fsmRepo.Create(testFSM); err != nil {
		t.Fatalf("Failed to create test FSM: %v", err)
	}

	// 测试修剪指甲动作
	trimNailsAction := &model.CatAction{
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionTrimNails,
		ActionDetail: "修剪指甲",
	}

	_, err := processor.ProcessAction(trimNailsAction)
	if err != nil {
		t.Errorf("ProcessAction failed: %v", err)
	}

	// 验证状态机是否更新
	freshFSM, err := fsmRepo.FindByID(1)
	if err != nil {
		t.Errorf("Failed to fetch FSM: %v", err)
	}

	if !freshFSM.TrimNailsTime.After(oldTime) {
		t.Errorf("TrimNailsTime should be updated after old time")
	}
}
