package repository

import (
	"testing"

	"fifu.fun/cat-dataserver/model"
)

func TestCatEventRepositoryCreate(t *testing.T) {
	setupTestDB(t)
	repo := NewCatEventRepository()

	event := &model.CatEvent{
		EventID:   1,
		EventType: model.CatSick,
		SiteID:    1,
		CatID:     1,
		Detail:    "测试事件详情",
	}

	err := repo.Create(event)
	if err != nil {
		t.Errorf("Failed to create event: %v", err)
	}

	if event.ID == 0 {
		t.Error("Expected non-zero ID after creation")
	}
}

func TestCatEventRepositoryFindByID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatEventRepository()

	event := &model.CatEvent{
		EventID:   1,
		EventType: model.CatSick,
		SiteID:    1,
		CatID:     1,
		Detail:    "测试事件详情",
	}
	repo.Create(event)

	foundEvent, err := repo.FindByID(event.EventID)
	if err != nil {
		t.Errorf("Failed to find event by ID: %v", err)
	}

	if foundEvent.EventType != model.CatSick {
		t.Errorf("Expected event type '%s', got '%s'", model.CatSick, foundEvent.EventType)
	}
}

func TestCatEventRepositoryFindByCatID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatEventRepository()

	event := &model.CatEvent{
		EventID:   1,
		EventType: model.CatSick,
		SiteID:    1,
		CatID:     1,
		Detail:    "测试事件详情",
	}
	repo.Create(event)

	events, err := repo.FindByCatID(1)
	if err != nil {
		t.Errorf("Failed to find events by CatID: %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected at least one event")
	}

	if events[0].CatID != 1 {
		t.Error("Expected CatID 1")
	}
}

func TestCatEventRepositoryFindBySiteID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatEventRepository()

	event := &model.CatEvent{
		EventID:   1,
		EventType: model.CatSick,
		SiteID:    1,
		CatID:     1,
		Detail:    "测试事件详情",
	}
	repo.Create(event)

	events, err := repo.FindBySiteID(1)
	if err != nil {
		t.Errorf("Failed to find events by SiteID: %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected at least one event")
	}

	if events[0].SiteID != 1 {
		t.Error("Expected SiteID 1")
	}
}

func TestCatEventRepositoryUpdate(t *testing.T) {
	setupTestDB(t)
	repo := NewCatEventRepository()

	event := &model.CatEvent{
		EventID:   1,
		EventType: model.CatSick,
		SiteID:    1,
		CatID:     1,
		Detail:    "测试事件详情",
	}
	repo.Create(event)

	event.Detail = "更新的事件详情"

	err := repo.Update(event)
	if err != nil {
		t.Errorf("Failed to update event: %v", err)
	}

	if event.Detail != "更新的事件详情" {
		t.Errorf("Expected detail '更新的事件详情', got '%s'", event.Detail)
	}
}

func TestCatEventRepositoryDelete(t *testing.T) {
	setupTestDB(t)
	repo := NewCatEventRepository()

	event := &model.CatEvent{
		EventID:   1,
		EventType: model.CatSick,
		SiteID:    1,
		CatID:     1,
		Detail:    "测试事件详情",
	}
	repo.Create(event)

	err := repo.Delete(event.EventID)
	if err != nil {
		t.Errorf("Failed to delete event: %v", err)
	}

	_, err = repo.FindByID(event.EventID)
	if err == nil {
		t.Error("Expected error when finding deleted event")
	}
}

func TestCatEventRepositoryFindPage(t *testing.T) {
	setupTestDB(t)
	repo := NewCatEventRepository()

	for i := 1; i <= 25; i++ {
		event := &model.CatEvent{
			EventID:   uint(i),
			EventType: model.CatSick,
			SiteID:    1,
			CatID:     1,
			Detail:    "测试事件详情",
		}
		repo.Create(event)
	}

	events, total, err := repo.FindPage(1, 10)
	if err != nil {
		t.Errorf("Failed to find page: %v", err)
	}

	if total != 25 {
		t.Errorf("Expected total 25, got %d", total)
	}

	if len(events) != 10 {
		t.Errorf("Expected 10 events, got %d", len(events))
	}
}

func TestNewCatEventRepository(t *testing.T) {
	repo := NewCatEventRepository()
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}
