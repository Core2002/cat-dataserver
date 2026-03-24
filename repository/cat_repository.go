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


// FindPage 分页查询 Cat
func (r *CatRepository) FindPage(page, pageSize int) ([]model.Cat, int64, error) {
	var cats []model.Cat
	var total int64

	// 查询总数
	if err := database.DB.Model(&model.Cat{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Find(&cats).Error
	return cats, total, err
}

// FindByID 根据 ID 查找 Cat
func (r *CatRepository) FindByID(catID uint) (*model.Cat, error) {
	var cat model.Cat
	err := database.DB.First(&cat, catID).Error
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
func (r *CatRepository) Delete(catID uint) error {
	return database.DB.Delete(&model.Cat{}, catID).Error
}
