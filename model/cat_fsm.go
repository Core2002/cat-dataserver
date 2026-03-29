package model

import (
	"time"

	"gorm.io/gorm"
)

type CatFSM struct {
	CatID         uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	SiteID        uint           `json:"site_id" gorm:"not null" binding:"required,min=1"`              // 猫所在设施ID
	TemperatureC  float32        `json:"temperature_c" gorm:"not null" binding:"required,gte=0,lte=50"` // 体温
	WeightKG      float32        `json:"weight_kg" gorm:"not null" binding:"required,gte=0.1,lte=25"`   // 体重
	TrimNailsTime time.Time      `json:"trim_nails_time" gorm:"not null" binding:"required"`            // 爪子修剪时间
}
