package model

import "time"

type Site struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"size:100;not null"`
	Address           string    `json:"address" gorm:"size:100;not null"`
	AdminPhoneNumber  string    `json:"phone_number" gorm:"size:100;not null"`
	LastDisinfectTime time.Time `json:"disinfect_time" gorm:"size:100;not null"` // 上次消毒时间
	LastFeedTime      time.Time `json:"feed_time" gorm:"not null"`               // 上次喂食时间
	LastGiveWaterTime time.Time `json:"give_water_time" gorm:"not null"`         // 上次喂水时间
	LastPlayTime      time.Time `json:"play_time" gorm:"not null"`               // 上次逗猫时间
}
