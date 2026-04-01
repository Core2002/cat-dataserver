package repository

import (
	"testing"

	"fifu.fun/cat-dataserver/model"
)

func TestCatFSMRepositoryCreate(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	fsm := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}

	err := repo.Create(fsm)
	if err != nil {
		t.Errorf("Failed to create FSM: %v", err)
	}

	if fsm.CatID == 0 {
		t.Error("Expected non-zero ID after creation")
	}
}

func TestCatFSMRepositoryFindByID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	fsm := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}
	repo.Create(fsm)

	foundFSM, err := repo.FindByID(fsm.CatID)
	if err != nil {
		t.Errorf("Failed to find FSM by ID: %v", err)
	}

	if foundFSM.CatID != 1 {
		t.Errorf("Expected CatID 1, got %d", foundFSM.CatID)
	}
}

func TestCatFSMRepositoryFindBySiteID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	fsm := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}
	repo.Create(fsm)

	fsms, err := repo.FindBySiteID(1)
	if err != nil {
		t.Errorf("Failed to find FSMs by SiteID: %v", err)
	}

	if len(fsms) == 0 {
		t.Error("Expected at least one FSM")
	}

	if fsms[0].SiteID != 1 {
		t.Error("Expected SiteID 1")
	}
}

func TestCatFSMRepositoryUpdate(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	fsm := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}
	repo.Create(fsm)

	fsm.WeightKG = 5.0

	err := repo.Update(fsm)
	if err != nil {
		t.Errorf("Failed to update FSM: %v", err)
	}

	if fsm.WeightKG != 5.0 {
		t.Errorf("Expected weight 5.0, got %f", fsm.WeightKG)
	}
}

func TestCatFSMRepositoryDelete(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	fsm := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}
	repo.Create(fsm)

	err := repo.Delete(fsm.CatID)
	if err != nil {
		t.Errorf("Failed to delete FSM: %v", err)
	}

	_, err = repo.FindByID(fsm.CatID)
	if err == nil {
		t.Error("Expected error when finding deleted FSM")
	}
}

func TestCatFSMRepositoryUpdateTemperature(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	fsm := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}
	repo.Create(fsm)

	err := repo.UpdateTemperature(1, 39.0)
	if err != nil {
		t.Errorf("Failed to update temperature: %v", err)
	}

	updatedFSM, _ := repo.FindByID(1)
	if updatedFSM.TemperatureC != 39.0 {
		t.Errorf("Expected temperature 39.0, got %f", updatedFSM.TemperatureC)
	}
}

func TestCatFSMRepositoryUpdateWeight(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	fsm := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}
	repo.Create(fsm)

	err := repo.UpdateWeight(1, 5.0)
	if err != nil {
		t.Errorf("Failed to update weight: %v", err)
	}

	updatedFSM, _ := repo.FindByID(1)
	if updatedFSM.WeightKG != 5.0 {
		t.Errorf("Expected weight 5.0, got %f", updatedFSM.WeightKG)
	}
}

func TestCatFSMRepositoryUpdateTrimNailsTime(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	fsm := &model.CatFSM{
		CatID:         1,
		SiteID:        1,
		TemperatureC:  38.5,
		WeightKG:      4.5,
		TrimNailsTime: model.TimeNow(),
	}
	repo.Create(fsm)

	newTime := "2024-01-01 12:00:00"
	err := repo.UpdateTrimNailsTime(1, newTime)
	if err != nil {
		t.Errorf("Failed to update trim nails time: %v", err)
	}

	updatedFSM, _ := repo.FindByID(1)
	if updatedFSM.TrimNailsTime.IsZero() {
		t.Error("Expected non-zero trim nails time")
	}
}

func TestCatFSMRepositoryFindPage(t *testing.T) {
	setupTestDB(t)
	repo := NewCatFSMRepository()

	for i := 1; i <= 25; i++ {
		fsm := &model.CatFSM{
			CatID:         uint(i),
			SiteID:        1,
			TemperatureC:  38.5,
			WeightKG:      4.5,
			TrimNailsTime: model.TimeNow(),
		}
		repo.Create(fsm)
	}

	fsms, total, err := repo.FindPage(1, 10)
	if err != nil {
		t.Errorf("Failed to find page: %v", err)
	}

	if total != 25 {
		t.Errorf("Expected total 25, got %d", total)
	}

	if len(fsms) != 10 {
		t.Errorf("Expected 10 FSMs, got %d", len(fsms))
	}
}

func TestNewCatFSMRepository(t *testing.T) {
	repo := NewCatFSMRepository()
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}
