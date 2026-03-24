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
