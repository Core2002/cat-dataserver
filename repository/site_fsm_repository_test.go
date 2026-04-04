package repository

import (
	"testing"
	"time"

	"fifu.fun/cat-dataserver/model"
)

func TestSiteFSMRepositoryCreate(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}

	err := repo.Create(fsm)
	if err != nil {
		t.Errorf("Failed to create site FSM: %v", err)
	}

	if fsm.SiteID == 0 {
		t.Error("Expected non-zero ID after creation")
	}
}

func TestSiteFSMRepositoryFindByID(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}
	repo.Create(fsm)

	foundFSM, err := repo.FindByID(fsm.SiteID)
	if err != nil {
		t.Errorf("Failed to find site FSM by ID: %v", err)
	}

	if foundFSM.SiteID != 1 {
		t.Errorf("Expected site ID 1, got %d", foundFSM.SiteID)
	}
}

func TestSiteFSMRepositoryFindBySiteID(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}
	repo.Create(fsm)

	foundFSM, err := repo.FindBySiteID(1)
	if err != nil {
		t.Errorf("Failed to find site FSM by SiteID: %v", err)
	}

	if foundFSM.SiteID != 1 {
		t.Errorf("Expected site ID 1, got %d", foundFSM.SiteID)
	}
}

func TestSiteFSMRepositoryUpdate(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}
	repo.Create(fsm)

	newTime := time.Now()
	fsm.LastDisinfectTime = &newTime
	fsm.SiteID = 2

	err := repo.Update(fsm)
	if err != nil {
		t.Errorf("Failed to update site FSM: %v", err)
	}

	if fsm.SiteID != 2 {
		t.Errorf("Expected site ID 2, got %d", fsm.SiteID)
	}
}

func TestSiteFSMRepositoryDelete(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}
	repo.Create(fsm)

	err := repo.Delete(fsm.SiteID)
	if err != nil {
		t.Errorf("Failed to delete site FSM: %v", err)
	}

	_, err = repo.FindByID(fsm.SiteID)
	if err == nil {
		t.Error("Expected error when finding deleted site FSM")
	}
}

func TestSiteFSMRepositoryUpdateDisinfectTime(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}
	repo.Create(fsm)

	newTime := "2024-01-01T12:00:00Z"
	err := repo.UpdateDisinfectTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update disinfect time: %v", err)
	}

	updatedFSM, _ := repo.FindBySiteID(1)
	if updatedFSM.LastDisinfectTime == nil {
		t.Error("Expected non-nil disinfect time")
	}
}

func TestSiteFSMRepositoryUpdateFeedTime(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}
	repo.Create(fsm)

	newTime := "2024-01-01T12:00:00Z"
	err := repo.UpdateFeedTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update feed time: %v", err)
	}

	updatedFSM, _ := repo.FindBySiteID(1)
	if updatedFSM.LastFeedTime == nil {
		t.Error("Expected non-nil feed time")
	}
}

func TestSiteFSMRepositoryUpdateGiveWaterTime(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}
	repo.Create(fsm)

	newTime := "2024-01-01T12:00:00Z"
	err := repo.UpdateGiveWaterTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update give water time: %v", err)
	}

	updatedFSM, _ := repo.FindBySiteID(1)
	if updatedFSM.LastGiveWaterTime == nil {
		t.Error("Expected non-nil give water time")
	}
}

func TestSiteFSMRepositoryUpdatePlayTime(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteFSMRepository()

	now := time.Now()
	fsm := &model.SiteFSM{
		SiteID:            1,
		LastDisinfectTime: &now,
		LastFeedTime:      &now,
		LastGiveWaterTime: &now,
		LastPlayTime:      &now,
	}
	repo.Create(fsm)

	newTime := "2024-01-01T12:00:00Z"
	err := repo.UpdatePlayTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update play time: %v", err)
	}

	updatedFSM, _ := repo.FindBySiteID(1)
	if updatedFSM.LastPlayTime == nil {
		t.Error("Expected non-nil play time")
	}
}

func TestNewSiteFSMRepository(t *testing.T) {
	repo := NewSiteFSMRepository()
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}
