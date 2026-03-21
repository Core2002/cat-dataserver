package repository

import (
	"testing"

	"fifu.fun/cat-dataserver/model"
)

func TestCatActionRepositoryCreate(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	action := &model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}

	err := repo.Create(action)
	if err != nil {
		t.Errorf("Failed to create action: %v", err)
	}

	if action.ID == 0 {
		t.Error("Expected non-zero ID after creation")
	}
}

func TestCatActionRepositoryFindAll(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	action := &model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}
	repo.Create(action)

	actions, err := repo.FindAll()
	if err != nil {
		t.Errorf("Failed to find all actions: %v", err)
	}

	if len(actions) == 0 {
		t.Error("Expected at least one action")
	}
}

func TestCatActionRepositoryFindByID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	action := &model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}
	repo.Create(action)

	foundAction, err := repo.FindByID(action.ActionID)
	if err != nil {
		t.Errorf("Failed to find action by ID: %v", err)
	}

	if foundAction.ActionType != model.CatActionFeed {
		t.Errorf("Expected action type '%s', got '%s'", model.CatActionFeed, foundAction.ActionType)
	}
}

func TestCatActionRepositoryFindByCatID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	action := &model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}
	repo.Create(action)

	actions, err := repo.FindByCatID(1)
	if err != nil {
		t.Errorf("Failed to find actions by CatID: %v", err)
	}

	if len(actions) == 0 {
		t.Error("Expected at least one action")
	}

	if actions[0].CatID != 1 {
		t.Error("Expected CatID 1")
	}
}

func TestCatActionRepositoryFindBySiteID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	action := &model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}
	repo.Create(action)

	actions, err := repo.FindBySiteID(1)
	if err != nil {
		t.Errorf("Failed to find actions by SiteID: %v", err)
	}

	if len(actions) == 0 {
		t.Error("Expected at least one action")
	}

	if actions[0].SiteID != 1 {
		t.Error("Expected SiteID 1")
	}
}

func TestCatActionRepositoryFindByUserID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	action := &model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}
	repo.Create(action)

	actions, err := repo.FindByUserID(1)
	if err != nil {
		t.Errorf("Failed to find actions by UserID: %v", err)
	}

	if len(actions) == 0 {
		t.Error("Expected at least one action")
	}

	if actions[0].UserID != 1 {
		t.Error("Expected UserID 1")
	}
}

func TestCatActionRepositoryUpdate(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	action := &model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}
	repo.Create(action)

	action.ActionDetail = "更新的操作详情"

	err := repo.Update(action)
	if err != nil {
		t.Errorf("Failed to update action: %v", err)
	}

	if action.ActionDetail != "更新的操作详情" {
		t.Errorf("Expected detail '更新的操作详情', got '%s'", action.ActionDetail)
	}
}

func TestCatActionRepositoryDelete(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	action := &model.CatAction{
		ActionID:     1,
		CatID:        1,
		SiteID:       1,
		UserID:       1,
		ActionType:   model.CatActionFeed,
		ActionDetail: "喂食测试",
	}
	repo.Create(action)

	err := repo.Delete(action.ActionID)
	if err != nil {
		t.Errorf("Failed to delete action: %v", err)
	}

	_, err = repo.FindByID(action.ActionID)
	if err == nil {
		t.Error("Expected error when finding deleted action")
	}
}

func TestCatActionRepositoryFindPage(t *testing.T) {
	setupTestDB(t)
	repo := NewCatActionRepository()

	for i := 1; i <= 25; i++ {
		action := &model.CatAction{
			ActionID:     uint(i),
			CatID:        1,
			SiteID:       1,
			UserID:       1,
			ActionType:   model.CatActionFeed,
			ActionDetail: "喂食测试",
		}
		repo.Create(action)
	}

	actions, total, err := repo.FindPage(1, 10)
	if err != nil {
		t.Errorf("Failed to find page: %v", err)
	}

	if total != 25 {
		t.Errorf("Expected total 25, got %d", total)
	}

	if len(actions) != 10 {
		t.Errorf("Expected 10 actions, got %d", len(actions))
	}
}

func TestNewCatActionRepository(t *testing.T) {
	repo := NewCatActionRepository()
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}
