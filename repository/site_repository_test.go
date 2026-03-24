package repository

import (
	"testing"

	"fifu.fun/cat-dataserver/model"
)

func TestSiteRepositoryCreate(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}

	err := repo.Create(site)
	if err != nil {
		t.Errorf("Failed to create site: %v", err)
	}

	if site.ID == 0 {
		t.Error("Expected non-zero ID after creation")
	}
}

func TestSiteRepositoryFindByID(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}
	repo.Create(site)

	foundSite, err := repo.FindByID(site.SiteID)
	if err != nil {
		t.Errorf("Failed to find site by ID: %v", err)
	}

	if foundSite.SiteName != "测试站点" {
		t.Errorf("Expected site name '测试站点', got '%s'", foundSite.SiteName)
	}
}

func TestSiteRepositoryUpdate(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}
	repo.Create(site)

	site.SiteName = "更新的站点"

	err := repo.Update(site)
	if err != nil {
		t.Errorf("Failed to update site: %v", err)
	}

	if site.SiteName != "更新的站点" {
		t.Errorf("Expected site name '更新的站点', got '%s'", site.SiteName)
	}
}

func TestSiteRepositoryDelete(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}
	repo.Create(site)

	err := repo.Delete(site.SiteID)
	if err != nil {
		t.Errorf("Failed to delete site: %v", err)
	}

	_, err = repo.FindByID(site.SiteID)
	if err == nil {
		t.Error("Expected error when finding deleted site")
	}
}

func TestSiteRepositoryUpdateDisinfectTime(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}
	repo.Create(site)

	newTime := "2024-01-01 12:00:00"
	err := repo.UpdateDisinfectTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update disinfect time: %v", err)
	}

	updatedSite, _ := repo.FindByID(1)
	if updatedSite.LastDisinfectTime.IsZero() {
		t.Error("Expected non-zero disinfect time")
	}
}

func TestSiteRepositoryUpdateFeedTime(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}
	repo.Create(site)

	newTime := "2024-01-01 12:00:00"
	err := repo.UpdateFeedTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update feed time: %v", err)
	}

	updatedSite, _ := repo.FindByID(1)
	if updatedSite.LastFeedTime.IsZero() {
		t.Error("Expected non-zero feed time")
	}
}

func TestSiteRepositoryUpdateGiveWaterTime(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}
	repo.Create(site)

	newTime := "2024-01-01 12:00:00"
	err := repo.UpdateGiveWaterTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update give water time: %v", err)
	}

	updatedSite, _ := repo.FindByID(1)
	if updatedSite.LastGiveWaterTime.IsZero() {
		t.Error("Expected non-zero give water time")
	}
}

func TestSiteRepositoryUpdatePlayTime(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteID:               1,
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13900139000",
		LastDisinfectTime:    model.TimeNow(),
		LastFeedTime:         model.TimeNow(),
		LastGiveWaterTime:    model.TimeNow(),
		LastPlayTime:         model.TimeNow(),
	}
	repo.Create(site)

	newTime := "2024-01-01 12:00:00"
	err := repo.UpdatePlayTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update play time: %v", err)
	}

	updatedSite, _ := repo.FindByID(1)
	if updatedSite.LastPlayTime.IsZero() {
		t.Error("Expected non-zero play time")
	}
}

func TestSiteRepositoryFindPage(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	for i := 1; i <= 25; i++ {
		site := &model.Site{
			SiteID:               uint(i),
			SiteName:             "测试站点",
			SiteAddress:          "测试地址",
			SiteAdminPhoneNumber: "13900139000",
			LastDisinfectTime:    model.TimeNow(),
			LastFeedTime:         model.TimeNow(),
			LastGiveWaterTime:    model.TimeNow(),
			LastPlayTime:         model.TimeNow(),
		}
		repo.Create(site)
	}

	sites, total, err := repo.FindPage(1, 10)
	if err != nil {
		t.Errorf("Failed to find page: %v", err)
	}

	if total != 25 {
		t.Errorf("Expected total 25, got %d", total)
	}

	if len(sites) != 10 {
		t.Errorf("Expected 10 sites, got %d", len(sites))
	}
}

func TestNewSiteRepository(t *testing.T) {
	repo := NewSiteRepository()
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}
