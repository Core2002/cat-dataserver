package repository

import (
	"testing"

	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
)

func setupTestDB(t *testing.T) {
	err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
}

func TestCatRepositoryCreate(t *testing.T) {
	setupTestDB(t)
	repo := NewCatRepository()

	cat := &model.Cat{
		CatID:             1,
		CatName:           "测试猫",
		CatPhotoUri:       "http://example.com/cat.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "张三",
		MasterPhoneNumber: "13800138000",
	}

	err := repo.Create(cat)
	if err != nil {
		t.Errorf("Failed to create cat: %v", err)
	}

	if cat.CatID == 0 {
		t.Error("Expected non-zero ID after creation")
	}
}

func TestCatRepositoryFindByID(t *testing.T) {
	setupTestDB(t)
	repo := NewCatRepository()

	cat := &model.Cat{
		CatID:             1,
		CatName:           "测试猫",
		CatPhotoUri:       "http://example.com/cat.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "张三",
		MasterPhoneNumber: "13800138000",
	}
	repo.Create(cat)

	foundCat, err := repo.FindByID(cat.CatID)
	if err != nil {
		t.Errorf("Failed to find cat by ID: %v", err)
	}

	if foundCat.CatName != "测试猫" {
		t.Errorf("Expected cat name '测试猫', got '%s'", foundCat.CatName)
	}
}

func TestCatRepositoryUpdate(t *testing.T) {
	setupTestDB(t)
	repo := NewCatRepository()

	cat := &model.Cat{
		CatID:             1,
		CatName:           "测试猫",
		CatPhotoUri:       "http://example.com/cat.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "张三",
		MasterPhoneNumber: "13800138000",
	}
	repo.Create(cat)

	updates := &model.Cat{
		CatName: "更新的猫",
	}

	err := repo.Update(cat, updates)
	if err != nil {
		t.Errorf("Failed to update cat: %v", err)
	}

	if cat.CatName != "更新的猫" {
		t.Errorf("Expected cat name '更新的猫', got '%s'", cat.CatName)
	}
}

func TestCatRepositoryDelete(t *testing.T) {
	setupTestDB(t)
	repo := NewCatRepository()

	cat := &model.Cat{
		CatID:             1,
		CatName:           "测试猫",
		CatPhotoUri:       "http://example.com/cat.jpg",
		CatType:           "英国短毛猫",
		CatGender:         "公",
		MasterName:        "张三",
		MasterPhoneNumber: "13800138000",
	}
	repo.Create(cat)

	err := repo.Delete(cat.CatID)
	if err != nil {
		t.Errorf("Failed to delete cat: %v", err)
	}

	_, err = repo.FindByID(cat.CatID)
	if err == nil {
		t.Error("Expected error when finding deleted cat")
	}
}

func TestCatRepositoryFindPage(t *testing.T) {
	setupTestDB(t)
	repo := NewCatRepository()

	for i := 1; i <= 25; i++ {
		cat := &model.Cat{
			CatID:             uint(i),
			CatName:           "测试猫",
			CatPhotoUri:       "http://example.com/cat.jpg",
			CatType:           "英国短毛猫",
			CatGender:         "公",
			MasterName:        "张三",
			MasterPhoneNumber: "13800138000",
		}
		repo.Create(cat)
	}

	cats, total, err := repo.FindPage(1, 10)
	if err != nil {
		t.Errorf("Failed to find page: %v", err)
	}

	if total != 25 {
		t.Errorf("Expected total 25, got %d", total)
	}

	if len(cats) != 10 {
		t.Errorf("Expected 10 cats, got %d", len(cats))
	}
}

func TestNewCatRepository(t *testing.T) {
	repo := NewCatRepository()
	if repo == nil {
		t.Error("Expected non-nil repository")
	}
}
