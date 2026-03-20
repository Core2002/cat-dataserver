package model

type Cat struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	Name              string `json:"name" gorm:"size:100;not null"`
	CatPhotoUri       string `json:"cat_photo_uri" gorm:"size:1000;not null"`
	MasterName        string `json:"master_name" gorm:"size:100;not null"`
	MasterPhoneNumber string `json:"master_phone_number" gorm:"size:100;not null"`
	CatType           string `json:"cat_type" gorm:"size:100;not null"`
	CatGender         string `json:"cat_gender" gorm:"size:100;not null"`
}

type CatProfile struct {
	CatAtom Cat
	FSM     CatFSM
}
