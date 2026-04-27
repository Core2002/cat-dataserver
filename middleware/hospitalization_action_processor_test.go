package middleware

import (
	"testing"
	"time"

	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"github.com/stretchr/testify/assert"
)

func TestProcessAction_AdmitAndDischarge_UpdateFSM(t *testing.T) {
	err := database.InitDB(":memory:")
	assert.NoError(t, err)

	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	actionRepo := repository.NewCatActionRepository()
	fsmRepo := repository.NewCatFSMRepository()

	site := &model.Site{
		SiteName:             "医院A",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13800138000",
	}
	err = siteRepo.Create(site)
	assert.NoError(t, err)

	cat := &model.Cat{
		CatID:             101,
		CatName:           "住院测试猫",
		CatPhotoUri:       "http://example.com/cat.jpg",
		CatType:           "英短",
		CatGender:         "公",
		MasterName:        "张三",
		MasterPhoneNumber: "13900139000",
	}
	err = catRepo.Create(cat)
	assert.NoError(t, err)

	fsm := &model.CatFSM{
		CatID:         cat.CatID,
		SiteID:        0,
		TemperatureC:  37.5,
		WeightKG:      4.0,
		TrimNailsTime: time.Now(),
	}
	err = fsmRepo.Create(fsm)
	assert.NoError(t, err)

	processor := NewActionProcessor(actionRepo, fsmRepo)

	admitAction := &model.CatAction{
		CatID:      cat.CatID,
		SiteID:     site.SiteID,
		UserID:     1,
		ActionType: model.CatActionAdmit,
		ActionDetail: `{"reason":"初次住院","notes":"观察","temperature_c":38.9,"weight_kg":4.6}`,
	}
	updatedFSM, err := processor.ProcessAction(admitAction)
	assert.NoError(t, err)
	assert.Equal(t, site.SiteID, updatedFSM.SiteID)
	assert.InDelta(t, 38.9, float64(updatedFSM.TemperatureC), 0.01)
	assert.InDelta(t, 4.6, float64(updatedFSM.WeightKG), 0.01)

	dischargeAction := &model.CatAction{
		CatID:      cat.CatID,
		SiteID:     site.SiteID,
		UserID:     1,
		ActionType: model.CatActionDischarge,
		ActionDetail: `{"reason":"恢复出院","notes":"状态稳定","temperature_c":38.2,"weight_kg":4.8}`,
	}
	updatedFSM, err = processor.ProcessAction(dischargeAction)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), updatedFSM.SiteID)
	assert.InDelta(t, 38.2, float64(updatedFSM.TemperatureC), 0.01)
	assert.InDelta(t, 4.8, float64(updatedFSM.WeightKG), 0.01)
}
