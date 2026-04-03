package middleware

import (
	"testing"
	"time"

	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
	err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
}

func TestProcessAction_TakeTemperature(t *testing.T) {
	setupTestDB(t)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteName:             "测试站点2",
		SiteAddress:          "测试地址2",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := siteRepo.Create(testSite)
	assert.NoError(t, err)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             2,
		CatName:           "测试猫2",
		CatPhotoUri:       "http://example.com/photo2.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人2",
		MasterPhoneNumber: "13800138000",
	}
	err = catRepo.Create(testCat)
	assert.NoError(t, err)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         2,
		SiteID:        testSite.SiteID,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	err = fsmRepo.Create(testFSM)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	action := &model.CatAction{
		ActionType:   model.CatActionTakeTemperature,
		CatID:        2,
		SiteID:       testSite.SiteID,
		ActionDetail: `{"temperature": 38.5}`,
	}

	fsm, err := processor.ProcessAction(action)
	assert.NoError(t, err)
	assert.NotNil(t, fsm)
	assert.InDelta(t, 38.5, float64(fsm.TemperatureC), 0.01)
}

func TestProcessAction_TrimNails(t *testing.T) {
	setupTestDB(t)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()
	now := time.Now()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteName:             "测试站点3",
		SiteAddress:          "测试地址3",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := siteRepo.Create(testSite)
	assert.NoError(t, err)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             3,
		CatName:           "测试猫3",
		CatPhotoUri:       "http://example.com/photo3.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人3",
		MasterPhoneNumber: "13800138000",
	}
	err = catRepo.Create(testCat)
	assert.NoError(t, err)

	// 创建测试用的 FSM 记录
	oldTime := time.Now().Add(-1 * time.Hour)
	testFSM := &model.CatFSM{
		CatID:         3,
		SiteID:        testSite.SiteID,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: oldTime,
	}
	err = fsmRepo.Create(testFSM)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	action := &model.CatAction{
		ActionType: model.CatActionTrimNails,
		CatID:      3,
		SiteID:     testSite.SiteID,
	}

	fsm, err := processor.ProcessAction(action)
	assert.NoError(t, err)
	assert.NotNil(t, fsm)
	assert.WithinDuration(t, now, fsm.TrimNailsTime, time.Second)
}

func TestUpdateTemperature(t *testing.T) {
	tests := []struct {
		name         string
		actionDetail string
		expectedTemp float64
		catID        uint
		siteName     string
		catName      string
	}{
		{"Valid temperature", `{"temperature": 38.5}`, 38.5, 10, "站点10", "猫10"},
		{"High temperature", `{"temperature": 40.0}`, 40.0, 11, "站点11", "猫11"},
		{"Low temperature", `{"temperature": 35.5}`, 35.5, 12, "站点12", "猫12"},
		{"Invalid format", `{"temperature": 0}`, 0, 13, "站点13", "猫13"},
		{"Empty string", "", 0, 14, "站点14", "猫14"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestDB(t)

			catRepo := repository.NewCatRepository()
			siteRepo := repository.NewSiteRepository()
			actionRepo := repository.NewCatActionRepository()
			fsmRepo := repository.NewCatFSMRepository()

			// 创建测试用的 Site 记录
			testSite := &model.Site{
				SiteName:             tt.siteName,
				SiteAddress:          "测试地址",
				SiteAdminPhoneNumber: "13800138000",
			}
			err := siteRepo.Create(testSite)
			assert.NoError(t, err)

			// 创建测试用的 Cat 记录
			testCat := &model.Cat{
				CatID:             tt.catID,
				CatName:           tt.catName,
				CatPhotoUri:       "http://example.com/photo.jpg",
				CatType:           "英国短毛猫",
				CatGender:         "公",
				MasterName:        "测试主人",
				MasterPhoneNumber: "13800138000",
			}
			err = catRepo.Create(testCat)
			assert.NoError(t, err)

			// 创建测试用的 FSM 记录
			testFSM := &model.CatFSM{
				CatID:         tt.catID,
				SiteID:        testSite.SiteID,
				TemperatureC:  38.5,
				WeightKG:      4.2,
				TrimNailsTime: time.Now(),
			}
			err = fsmRepo.Create(testFSM)
			assert.NoError(t, err)

			processor := NewActionProcessor(actionRepo, fsmRepo)
			action := &model.CatAction{
				ActionType:   model.CatActionTakeTemperature,
				CatID:        tt.catID,
				SiteID:       testSite.SiteID,
				ActionDetail: tt.actionDetail,
			}

			fsm, err := processor.ProcessAction(action)
			if tt.expectedTemp == 0 {
				// For invalid cases, we expect no error but fsm might be nil
				// The action is still saved, but temperature is not updated
				if tt.actionDetail != "" {
					assert.NoError(t, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, fsm)
				assert.InDelta(t, tt.expectedTemp, float64(fsm.TemperatureC), 0.01)
			}
		})
	}
}

func TestUpdateTrimNailsTime(t *testing.T) {
	setupTestDB(t)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()
	now := time.Now()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteName:             "测试站点20",
		SiteAddress:          "测试地址20",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := siteRepo.Create(testSite)
	assert.NoError(t, err)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             20,
		CatName:           "测试猫20",
		CatPhotoUri:       "http://example.com/photo20.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人20",
		MasterPhoneNumber: "13800138000",
	}
	err = catRepo.Create(testCat)
	assert.NoError(t, err)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         20,
		SiteID:        testSite.SiteID,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	err = fsmRepo.Create(testFSM)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	action := &model.CatAction{
		ActionType: model.CatActionTrimNails,
		CatID:      20,
		SiteID:     testSite.SiteID,
	}

	fsm, err := processor.ProcessAction(action)
	assert.NoError(t, err)
	assert.NotNil(t, fsm)
	assert.WithinDuration(t, now, fsm.TrimNailsTime, time.Second)
}

func TestProcessAction_HealthCheck(t *testing.T) {
	setupTestDB(t)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteName:             "测试站点HealthCheck",
		SiteAddress:          "测试地址HealthCheck",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := siteRepo.Create(testSite)
	assert.NoError(t, err)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             30,
		CatName:           "测试猫HealthCheck",
		CatPhotoUri:       "http://example.com/photo30.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人HealthCheck",
		MasterPhoneNumber: "13800138000",
	}
	err = catRepo.Create(testCat)
	assert.NoError(t, err)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         30,
		SiteID:        testSite.SiteID,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	err = fsmRepo.Create(testFSM)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	action := &model.CatAction{
		ActionType:   model.CatActionHealthCheck,
		CatID:        30,
		SiteID:       testSite.SiteID,
		ActionDetail: `{"temperature": 39.0, "weight": 4.8, "notes": "体检正常"}`,
	}

	fsm, err := processor.ProcessAction(action)
	assert.NoError(t, err)
	assert.NotNil(t, fsm)
	assert.InDelta(t, 39.0, float64(fsm.TemperatureC), 0.01)
	assert.InDelta(t, 4.8, float64(fsm.WeightKG), 0.01)
}

func TestProcessAction_Deworm(t *testing.T) {
	setupTestDB(t)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteName:             "测试站点Deworm",
		SiteAddress:          "测试地址Deworm",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := siteRepo.Create(testSite)
	assert.NoError(t, err)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             31,
		CatName:           "测试猫Deworm",
		CatPhotoUri:       "http://example.com/photo31.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人Deworm",
		MasterPhoneNumber: "13800138000",
	}
	err = catRepo.Create(testCat)
	assert.NoError(t, err)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         31,
		SiteID:        testSite.SiteID,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	err = fsmRepo.Create(testFSM)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	action := &model.CatAction{
		ActionType:   model.CatActionDeworm,
		CatID:        31,
		SiteID:       testSite.SiteID,
		ActionDetail: `{"drug_name": "福来恩", "dosage": "1ml"}`,
	}

	fsm, err := processor.ProcessAction(action)
	assert.NoError(t, err)
	assert.NotNil(t, fsm)
	// 驱虫动作不更新 FSM
}

func TestProcessAction_Vaccinate(t *testing.T) {
	setupTestDB(t)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteName:             "测试站点Vaccinate",
		SiteAddress:          "测试地址Vaccinate",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := siteRepo.Create(testSite)
	assert.NoError(t, err)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             32,
		CatName:           "测试猫Vaccinate",
		CatPhotoUri:       "http://example.com/photo32.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人Vaccinate",
		MasterPhoneNumber: "13800138000",
	}
	err = catRepo.Create(testCat)
	assert.NoError(t, err)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         32,
		SiteID:        testSite.SiteID,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	err = fsmRepo.Create(testFSM)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	action := &model.CatAction{
		ActionType:   model.CatActionVaccinate,
		CatID:        32,
		SiteID:       testSite.SiteID,
		ActionDetail: `{"vaccine_name": "猫三联", "batch_no": "B2024001"}`,
	}

	fsm, err := processor.ProcessAction(action)
	assert.NoError(t, err)
	assert.NotNil(t, fsm)
	// 疫苗动作不更新 FSM
}

func TestProcessAction_Sterilize(t *testing.T) {
	setupTestDB(t)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteName:             "测试站点Sterilize",
		SiteAddress:          "测试地址Sterilize",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := siteRepo.Create(testSite)
	assert.NoError(t, err)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             33,
		CatName:           "测试猫Sterilize",
		CatPhotoUri:       "http://example.com/photo33.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人Sterilize",
		MasterPhoneNumber: "13800138000",
	}
	err = catRepo.Create(testCat)
	assert.NoError(t, err)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         33,
		SiteID:        testSite.SiteID,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	err = fsmRepo.Create(testFSM)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	action := &model.CatAction{
		ActionType:   model.CatActionSterilize,
		CatID:        33,
		SiteID:       testSite.SiteID,
		ActionDetail: `{"notes": "手术顺利"}`,
	}

	fsm, err := processor.ProcessAction(action)
	assert.NoError(t, err)
	assert.NotNil(t, fsm)
	// 绝育动作不更新 FSM
}

func TestProcessAction_Bathing(t *testing.T) {
	setupTestDB(t)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()

	// 创建测试用的 Site 记录
	testSite := &model.Site{
		SiteName:             "测试站点Bathing",
		SiteAddress:          "测试地址Bathing",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := siteRepo.Create(testSite)
	assert.NoError(t, err)

	// 创建测试用的 Cat 记录
	testCat := &model.Cat{
		CatID:             34,
		CatName:           "测试猫Bathing",
		CatPhotoUri:       "http://example.com/photo34.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "测试主人Bathing",
		MasterPhoneNumber: "13800138000",
	}
	err = catRepo.Create(testCat)
	assert.NoError(t, err)

	// 创建测试用的 FSM 记录
	testFSM := &model.CatFSM{
		CatID:         34,
		SiteID:        testSite.SiteID,
		TemperatureC:  38.5,
		WeightKG:      4.2,
		TrimNailsTime: time.Now(),
	}
	err = fsmRepo.Create(testFSM)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	action := &model.CatAction{
		ActionType:   model.CatActionBathing,
		CatID:        34,
		SiteID:       testSite.SiteID,
		ActionDetail: `{"notes": "洗澡完成"}`,
	}

	fsm, err := processor.ProcessAction(action)
	assert.NoError(t, err)
	assert.NotNil(t, fsm)
	// 洗澡动作不更新 FSM
}
