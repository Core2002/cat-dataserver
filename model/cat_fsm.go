package model

import "time"

type CatFSM struct {
	CatID         uint      `json:"cat_id" gorm:"primaryKey"`        // 猫ID
	SiteID        uint      `json:"site_id" gorm:"not null"`         // 猫所在设施ID
	TemperatureC  float32   `json:"temperature_c" gorm:"not null"`   // 体温
	WeightKG      float32   `json:"weight_kg" gorm:"not null"`       // 体重
	TrimNailsTime time.Time `json:"trim_nails_time" gorm:"not null"` // 爪子修剪时间
}
