package repository

import (
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
)

// SiteRepository Site 数据访问层
type SiteRepository struct{}

// NewSiteRepository 创建 SiteRepository 实例
func NewSiteRepository() *SiteRepository {
	return &SiteRepository{}
}

// FindAll 查找所有 Site
func (r *SiteRepository) FindAll() ([]model.Site, error) {
	var sites []model.Site
	err := database.DB.Find(&sites).Error
	return sites, err
}

// FindPage 分页查询 Site
func (r *SiteRepository) FindPage(page, pageSize int) ([]model.Site, int64, error) {
	var sites []model.Site
	var total int64

	if err := database.DB.Model(&model.Site{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Find(&sites).Error
	return sites, total, err
}

// FindByID 根据 ID 查找 Site
func (r *SiteRepository) FindByID(siteID uint) (*model.Site, error) {
	var site model.Site
	err := database.DB.First(&site, siteID).Error
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// FindByName 根据名称查找 Site
func (r *SiteRepository) FindByName(name string) (*model.Site, error) {
	var site model.Site
	err := database.DB.Where("site_name = ?", name).First(&site).Error
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// Create 创建 Site
func (r *SiteRepository) Create(site *model.Site) error {
	return database.DB.Create(site).Error
}

// Update 更新 Site
func (r *SiteRepository) Update(site *model.Site) error {
	return database.DB.Save(site).Error
}

// Delete 删除 Site
func (r *SiteRepository) Delete(siteID uint) error {
	return database.DB.Delete(&model.Site{}, siteID).Error
}

// UpdateDisinfectTime 更新消毒时间
func (r *SiteRepository) UpdateDisinfectTime(siteID uint, time interface{}) error {
	return database.DB.Model(&model.Site{}).Where("site_id = ?", siteID).Update("last_disinfect_time", time).Error
}

// UpdateFeedTime 更新喂食时间
func (r *SiteRepository) UpdateFeedTime(siteID uint, time interface{}) error {
	return database.DB.Model(&model.Site{}).Where("site_id = ?", siteID).Update("last_feed_time", time).Error
}

// UpdateGiveWaterTime 更新喂水时间
func (r *SiteRepository) UpdateGiveWaterTime(siteID uint, time interface{}) error {
	return database.DB.Model(&model.Site{}).Where("site_id = ?", siteID).Update("last_give_water_time", time).Error
}

// UpdatePlayTime 更新逗猫时间
func (r *SiteRepository) UpdatePlayTime(siteID uint, time interface{}) error {
	return database.DB.Model(&model.Site{}).Where("site_id = ?", siteID).Update("last_play_time", time).Error
}
