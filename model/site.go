package model

import "time"

type Site struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	SiteName             string    `json:"site_name" gorm:"size:100;not null"`
	SiteAddress          string    `json:"site_address" gorm:"size:100;not null"`
	SiteAdminPhoneNumber string    `json:"site_admin_phone_number" gorm:"size:100;not null"`
	LastDisinfectTime    time.Time `json:"disinfect_time" gorm:"size:100;not null"` // 上次消毒时间
	LastFeedTime         time.Time `json:"feed_time" gorm:"not null"`               // 上次喂食时间
	LastGiveWaterTime    time.Time `json:"give_water_time" gorm:"not null"`         // 上次喂水时间
	LastPlayTime         time.Time `json:"play_time" gorm:"not null"`               // 上次逗猫时间
}
