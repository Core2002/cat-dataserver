package repository

import (
	"testing"

	"fifu.fun/cat-dataserver/model"
)

func TestSiteRepositoryCreate(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13800138000",
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
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := repo.Create(site)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	foundSite, err := repo.FindByID(site.ID)
	if err != nil {
		t.Errorf("Failed to find site by ID: %v", err)
	}

	if foundSite.SiteName != "测试站点" {
		t.Errorf("Expected site name '测试站点', got '%s'", foundSite.SiteName)
	}
}

func TestSiteRepositoryDelete(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	site := &model.Site{
		SiteName:             "测试站点",
		SiteAddress:          "测试地址",
		SiteAdminPhoneNumber: "13800138000",
	}
	err := repo.Create(site)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	err = repo.Delete(site.ID)
	if err != nil {
		t.Errorf("Failed to delete site: %v", err)
	}

	_, err = repo.FindByID(site.ID)
	if err == nil {
		t.Error("Expected error when finding deleted site")
	}
}

func TestSiteRepositoryFindPage(t *testing.T) {
	setupTestDB(t)
	repo := NewSiteRepository()

	// Create multiple sites with unique names to avoid conflicts
	for i := 1; i <= 25; i++ {
		site := &model.Site{
			SiteName:             "测试站点",
			SiteAddress:          "测试地址",
			SiteAdminPhoneNumber: "13800138000",
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
