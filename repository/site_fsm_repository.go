package repository

import (
	"time"

	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
)

// SiteFSMRepository SiteFSM 数据访问层
type SiteFSMRepository struct{}

// NewSiteFSMRepository 创建 SiteFSMRepository 实例
func NewSiteFSMRepository() *SiteFSMRepository {
	return &SiteFSMRepository{}
}

// FindPage 分页查找 SiteFSM
func (r *SiteFSMRepository) FindPage(page, pageSize int) ([]model.SiteFSM, int64, error) {
	var fsms []model.SiteFSM
	var total int64

	offset := (page - 1) * pageSize

	if err := database.DB.Model(&model.SiteFSM{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Offset(offset).Limit(pageSize).Find(&fsms).Error; err != nil {
		return nil, 0, err
	}

	return fsms, total, nil
}

// FindByID 根据 ID 查找 SiteFSM
func (r *SiteFSMRepository) FindByID(id uint) (*model.SiteFSM, error) {
	var fsm model.SiteFSM
	err := database.DB.First(&fsm, id).Error
	if err != nil {
		return nil, err
	}
	return &fsm, nil
}

// FindBySiteID 根据 SiteID 查找 SiteFSM
func (r *SiteFSMRepository) FindBySiteID(siteID uint) (*model.SiteFSM, error) {
	var fsm model.SiteFSM
	err := database.DB.Where("site_id = ?", siteID).First(&fsm).Error
	if err != nil {
		return nil, err
	}
	return &fsm, nil
}

// GetOrCreateBySiteID 获取或创建 SiteFSM（原子操作）
func (r *SiteFSMRepository) GetOrCreateBySiteID(siteID uint) (*model.SiteFSM, error) {
	var fsm model.SiteFSM
	result := database.DB.Where(model.SiteFSM{SiteID: siteID}).FirstOrCreate(&fsm)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fsm, nil
}

// Create 创建 SiteFSM
func (r *SiteFSMRepository) Create(fsm *model.SiteFSM) error {
	return database.DB.Create(fsm).Error
}

// Update 更新 SiteFSM
func (r *SiteFSMRepository) Update(fsm *model.SiteFSM) error {
	return database.DB.Save(fsm).Error
}

// Delete 删除 SiteFSM
func (r *SiteFSMRepository) Delete(id uint) error {
	return database.DB.Delete(&model.SiteFSM{}, id).Error
}

// parseTime 解析时间字符串，空字符串返回 nil 表示不更新
func parseTime(timeStr string) (*time.Time, error) {
	// 空字符串表示不更新
	if timeStr == "" {
		return nil, nil
	}

	// 解析 RFC3339 格式时间
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateDisinfectTime 更新消毒时间
func (r *SiteFSMRepository) UpdateDisinfectTime(siteID uint, timeStr string) error {
	t, err := parseTime(timeStr)
	if err != nil {
		return err
	}
	if t == nil {
		return nil // 空字符串，跳过更新
	}
	return database.DB.Model(&model.SiteFSM{}).Where("site_id = ?", siteID).Update("last_disinfect_time", t).Error
}

// UpdateFeedTime 更新喂食时间
func (r *SiteFSMRepository) UpdateFeedTime(siteID uint, timeStr string) error {
	t, err := parseTime(timeStr)
	if err != nil {
		return err
	}
	if t == nil {
		return nil // 空字符串，跳过更新
	}
	return database.DB.Model(&model.SiteFSM{}).Where("site_id = ?", siteID).Update("last_feed_time", t).Error
}

// UpdateGiveWaterTime 更新喂水时间
func (r *SiteFSMRepository) UpdateGiveWaterTime(siteID uint, timeStr string) error {
	t, err := parseTime(timeStr)
	if err != nil {
		return err
	}
	if t == nil {
		return nil // 空字符串，跳过更新
	}
	return database.DB.Model(&model.SiteFSM{}).Where("site_id = ?", siteID).Update("last_give_water_time", t).Error
}

// UpdatePlayTime 更新逗猫时间
func (r *SiteFSMRepository) UpdatePlayTime(siteID uint, timeStr string) error {
	t, err := parseTime(timeStr)
	if err != nil {
		return err
	}
	if t == nil {
		return nil // 空字符串，跳过更新
	}
	return database.DB.Model(&model.SiteFSM{}).Where("site_id = ?", siteID).Update("last_play_time", t).Error
}

// UpdateCleanLitterTime 更新清理猫砂时间
func (r *SiteFSMRepository) UpdateCleanLitterTime(siteID uint, timeStr string) error {
	t, err := parseTime(timeStr)
	if err != nil {
		return err
	}
	if t == nil {
		return nil // 空字符串，跳过更新
	}
	return database.DB.Model(&model.SiteFSM{}).Where("site_id = ?", siteID).Update("last_clean_litter_time", t).Error
}
