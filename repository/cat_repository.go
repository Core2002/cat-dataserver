package repository

import (
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
)

// CatRepository Cat 数据访问层
type CatRepository struct{}

// NewCatRepository 创建 CatRepository 实例
func NewCatRepository() *CatRepository {
	return &CatRepository{}
}

// FindAll 查找所有 Cat
func (r *CatRepository) FindAll() ([]model.Cat, error) {
	var cats []model.Cat
	err := database.DB.Find(&cats).Error
	return cats, err
}

// FindByID 根据 ID 查找 Cat
func (r *CatRepository) FindByID(id string) (*model.Cat, error) {
	var cat model.Cat
	err := database.DB.First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// Create 创建 Cat
func (r *CatRepository) Create(cat *model.Cat) error {
	return database.DB.Create(cat).Error
}

// Update 更新 Cat
func (r *CatRepository) Update(cat *model.Cat, updates *model.Cat) error {
	return database.DB.Model(cat).Updates(updates).Error
}

// Delete 删除 Cat
func (r *CatRepository) Delete(id string) error {
	return database.DB.Delete(&model.Cat{}, id).Error
}
