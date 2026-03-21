package repository

import (
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
)

// CatFSMRepository CatFSM 数据访问层
type CatFSMRepository struct{}

// NewCatFSMRepository 创建 CatFSMRepository 实例
func NewCatFSMRepository() *CatFSMRepository {
	return &CatFSMRepository{}
}

// FindAll 查找所有 CatFSM
func (r *CatFSMRepository) FindAll() ([]model.CatFSM, error) {
	var fsms []model.CatFSM
	err := database.DB.Find(&fsms).Error
	return fsms, err
}

// FindPage 分页查询 CatFSM
func (r *CatFSMRepository) FindPage(page, pageSize int) ([]model.CatFSM, int64, error) {
	var fsms []model.CatFSM
	var total int64

	if err := database.DB.Model(&model.CatFSM{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Find(&fsms).Error
	return fsms, total, err
}

// FindByID 根据 ID 查找 CatFSM
func (r *CatFSMRepository) FindByID(catID uint) (*model.CatFSM, error) {
	var fsm model.CatFSM
	err := database.DB.First(&fsm, catID).Error
	if err != nil {
		return nil, err
	}
	return &fsm, nil
}

// FindBySiteID 根据 SiteID 查找猫状态
func (r *CatFSMRepository) FindBySiteID(siteID uint) ([]model.CatFSM, error) {
	var fsms []model.CatFSM
	err := database.DB.Where("site_id = ?", siteID).Find(&fsms).Error
	return fsms, err
}

// Create 创建 CatFSM
func (r *CatFSMRepository) Create(fsm *model.CatFSM) error {
	return database.DB.Create(fsm).Error
}

// Update 更新 CatFSM
func (r *CatFSMRepository) Update(fsm *model.CatFSM) error {
	return database.DB.Save(fsm).Error
}

// Delete 删除 CatFSM
func (r *CatFSMRepository) Delete(catID uint) error {
	return database.DB.Delete(&model.CatFSM{}, catID).Error
}

// UpdateTemperature 更新体温
func (r *CatFSMRepository) UpdateTemperature(catID uint, temperature float32) error {
	return database.DB.Model(&model.CatFSM{}).Where("cat_id = ?", catID).Update("temperature_c", temperature).Error
}

// UpdateWeight 更新体重
func (r *CatFSMRepository) UpdateWeight(catID uint, weight float32) error {
	return database.DB.Model(&model.CatFSM{}).Where("cat_id = ?", catID).Update("weight_kg", weight).Error
}

// UpdateTrimNailsTime 更新修剪指甲时间
func (r *CatFSMRepository) UpdateTrimNailsTime(catID uint, time interface{}) error {
	return database.DB.Model(&model.CatFSM{}).Where("cat_id = ?", catID).Update("trim_nails_time", time).Error
}
