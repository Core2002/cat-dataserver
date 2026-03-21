package repository

import (
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
)

// CatEventRepository CatEvent 数据访问层
type CatEventRepository struct{}

// NewCatEventRepository 创建 CatEventRepository 实例
func NewCatEventRepository() *CatEventRepository {
	return &CatEventRepository{}
}

// FindAll 查找所有 CatEvent
func (r *CatEventRepository) FindAll() ([]model.CatEvent, error) {
	var events []model.CatEvent
	err := database.DB.Find(&events).Error
	return events, err
}

// FindPage 分页查询 CatEvent
func (r *CatEventRepository) FindPage(page, pageSize int) ([]model.CatEvent, int64, error) {
	var events []model.CatEvent
	var total int64

	if err := database.DB.Model(&model.CatEvent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := database.DB.Offset(offset).Limit(pageSize).Find(&events).Error
	return events, total, err
}

// FindByID 根据 ID 查找 CatEvent
func (r *CatEventRepository) FindByID(eventID uint) (*model.CatEvent, error) {
	var event model.CatEvent
	err := database.DB.First(&event, eventID).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// FindByCatID 根据 CatID 查找事件
func (r *CatEventRepository) FindByCatID(catID uint) ([]model.CatEvent, error) {
	var events []model.CatEvent
	err := database.DB.Where("cat_id = ?", catID).Find(&events).Error
	return events, err
}

// FindBySiteID 根据 SiteID 查找事件
func (r *CatEventRepository) FindBySiteID(siteID uint) ([]model.CatEvent, error) {
	var events []model.CatEvent
	err := database.DB.Where("site_id = ?", siteID).Find(&events).Error
	return events, err
}

// Create 创建 CatEvent
func (r *CatEventRepository) Create(event *model.CatEvent) error {
	return database.DB.Create(event).Error
}

// Update 更新 CatEvent
func (r *CatEventRepository) Update(event *model.CatEvent) error {
	return database.DB.Save(event).Error
}

// Delete 删除 CatEvent
func (r *CatEventRepository) Delete(eventID uint) error {
	return database.DB.Delete(&model.CatEvent{}, eventID).Error
}
