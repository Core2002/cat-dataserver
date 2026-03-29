package model

import (
	"time"

	"gorm.io/gorm"
)

type Site struct {
	SiteID               uint `gorm:"primarykey"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	SiteName             string         `json:"site_name" gorm:"size:100;not null" binding:"required,min=1,max=100"`    // 站点名称
	SiteAddress          string         `json:"site_address" gorm:"size:100;not null" binding:"required,min=1,max=100"` // 站点地址
	SiteAdminPhoneNumber string         `json:"site_admin_phone_number" gorm:"size:100;not null" binding:"required"`    // 站点管理员电话
}
