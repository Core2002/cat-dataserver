package model

import (
	"time"

	"gorm.io/gorm"
)

type Cat struct {
	CatID             uint `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	CatName           string         `json:"cat_name" gorm:"size:100;not null" binding:"required,min=1,max=100"`    // 猫名
	CatPhotoUri       string         `json:"cat_photo_uri" gorm:"size:1000;not null" binding:"required,url"`        // 猫照片
	CatType           string         `json:"cat_type" gorm:"size:100;not null" binding:"required,min=1,max=50"`     // 猫种类
	CatGender         string         `json:"cat_gender" gorm:"size:100;not null" binding:"required"`                // 猫性别
	MasterName        string         `json:"master_name" gorm:"size:100;not null" binding:"required,min=1,max=100"` // 主人姓名
	MasterPhoneNumber string         `json:"master_phone_number" gorm:"size:100;not null" binding:"required"`       // 主人电话
}

type CatProfile struct {
	CatAtom Cat
	FSM     CatFSM
}
