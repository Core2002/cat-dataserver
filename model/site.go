package model

import (
	"time"

	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	SiteID               uint      `json:"site_id" gorm:"primaryKey"`                        // 站点ID
	SiteName             string    `json:"site_name" gorm:"size:100;not null"`               // 站点名称
	SiteAddress          string    `json:"site_address" gorm:"size:100;not null"`            // 站点地址
	SiteAdminPhoneNumber string    `json:"site_admin_phone_number" gorm:"size:100;not null"` // 站点管理员电话
	LastDisinfectTime    time.Time `json:"disinfect_time" gorm:"size:100;not null"`          // 上次消毒时间
	LastFeedTime         time.Time `json:"feed_time" gorm:"not null"`                        // 上次喂食时间
	LastGiveWaterTime    time.Time `json:"give_water_time" gorm:"not null"`                  // 上次喂水时间
	LastPlayTime         time.Time `json:"play_time" gorm:"not null"`                        // 上次逗猫时间
}
