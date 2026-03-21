package model

import "gorm.io/gorm"

type Cat struct {
	gorm.Model
	CatID             uint   `json:"cat_id" gorm:"primaryKey"`                     // 猫ID
	CatName           string `json:"cat_name" gorm:"size:100;not null"`            // 猫名
	CatPhotoUri       string `json:"cat_photo_uri" gorm:"size:1000;not null"`      // 猫照片
	CatType           string `json:"cat_type" gorm:"size:100;not null"`            // 猫种类
	CatGender         string `json:"cat_gender" gorm:"size:100;not null"`          // 猫性别
	MasterName        string `json:"master_name" gorm:"size:100;not null"`         // 主人姓名
	MasterPhoneNumber string `json:"master_phone_number" gorm:"size:100;not null"` // 主人电话
}

type CatProfile struct {
	CatAtom Cat
	FSM     CatFSM
}
