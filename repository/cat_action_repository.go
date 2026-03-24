package repository

import (
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
)

// CatActionRepository CatAction 数据访问层
type CatActionRepository struct{}

// NewCatActionRepository 创建 CatActionRepository 实例
func NewCatActionRepository() *CatActionRepository {
	return &CatActionRepository{}
}

// FindPage 分页查询 CatAction
func (r *CatActionRepository) FindPage(page, pageSize int) ([]model.CatAction, int64, error) {
	var actions []model.CatAction
	var total int64

	if err := database.DB.Model(&model.CatAction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Find(&actions).Error
	return actions, total, err
}

// FindByID 根据 ID 查找 CatAction
func (r *CatActionRepository) FindByID(actionID uint) (*model.CatAction, error) {
	var action model.CatAction
	err := database.DB.First(&action, actionID).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// FindByCatID 根据 CatID 查找操作记录
func (r *CatActionRepository) FindByCatID(catID uint) ([]model.CatAction, error) {
	var actions []model.CatAction
	err := database.DB.Where("cat_id = ?", catID).Find(&actions).Error
	return actions, err
}

// FindBySiteID 根据 SiteID 查找操作记录
func (r *CatActionRepository) FindBySiteID(siteID uint) ([]model.CatAction, error) {
	var actions []model.CatAction
	err := database.DB.Where("site_id = ?", siteID).Find(&actions).Error
	return actions, err
}

// FindByUserID 根据 UserID 查找操作记录
func (r *CatActionRepository) FindByUserID(userID uint) ([]model.CatAction, error) {
	var actions []model.CatAction
	err := database.DB.Where("user_id = ?", userID).Find(&actions).Error
	return actions, err
}

// Create 创建 CatAction
func (r *CatActionRepository) Create(action *model.CatAction) error {
	return database.DB.Create(action).Error
}

// Update 更新 CatAction
func (r *CatActionRepository) Update(action *model.CatAction) error {
	return database.DB.Save(action).Error
}

// Delete 删除 CatAction
func (r *CatActionRepository) Delete(actionID uint) error {
	return database.DB.Delete(&model.CatAction{}, actionID).Error
}
