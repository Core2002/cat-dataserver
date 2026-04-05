package repository

import (
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
)

// SiteActionRepository SiteAction 数据访问层
type SiteActionRepository struct{}

// NewSiteActionRepository 创建 SiteActionRepository 实例
func NewSiteActionRepository() *SiteActionRepository {
	return &SiteActionRepository{}
}

// FindPage 分页查询 SiteAction
func (r *SiteActionRepository) FindPage(page, pageSize int) ([]model.SiteAction, int64, error) {
	var actions []model.SiteAction
	var total int64

	if err := database.DB.Model(&model.SiteAction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Find(&actions).Error
	return actions, total, err
}

// FindByID 根据 ID 查找 SiteAction
func (r *SiteActionRepository) FindByID(actionID uint) (*model.SiteAction, error) {
	var action model.SiteAction
	err := database.DB.First(&action, actionID).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// FindBySiteID 根据 SiteID 查找操作记录
func (r *SiteActionRepository) FindBySiteID(siteID uint) ([]model.SiteAction, error) {
	var actions []model.SiteAction
	err := database.DB.Where("site_id = ?", siteID).Find(&actions).Error
	return actions, err
}

// FindByUserID 根据 UserID 查找操作记录
func (r *SiteActionRepository) FindByUserID(userID uint) ([]model.SiteAction, error) {
	var actions []model.SiteAction
	err := database.DB.Where("user_id = ?", userID).Find(&actions).Error
	return actions, err
}

// Create 创建 SiteAction
func (r *SiteActionRepository) Create(action *model.SiteAction) error {
	return database.DB.Create(action).Error
}

// Update 更新 SiteAction
func (r *SiteActionRepository) Update(action *model.SiteAction) error {
	return database.DB.Save(action).Error
}

// Delete 删除 SiteAction
func (r *SiteActionRepository) Delete(actionID uint) error {
	return database.DB.Delete(&model.SiteAction{}, actionID).Error
}
