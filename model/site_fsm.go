package model

import (
	"time"

	"gorm.io/gorm"
)

// SiteFSM 站点状态机信息
type SiteFSM struct {
	SiteID            uint           `json:"site_id" gorm:"primarykey"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	LastDisinfectTime *time.Time     `json:"last_disinfect_time" gorm:"type:datetime"`  // 上次消毒时间
	LastFeedTime      *time.Time     `json:"last_feed_time" gorm:"type:datetime"`       // 上次喂食时间
	LastGiveWaterTime *time.Time     `json:"last_give_water_time" gorm:"type:datetime"` // 上次喂水时间
	LastPlayTime      *time.Time     `json:"last_play_time" gorm:"type:datetime"`       // 上次逗猫时间
}

// TableName 指定表名
func (SiteFSM) TableName() string {
	return "site_fsm"
}
