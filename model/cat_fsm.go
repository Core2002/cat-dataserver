package model

import "time"

type CatFSM struct {
	CatID         uint      `json:"cat_id" gorm:"primaryKey"`
	HospitalID    uint      `json:"hospital_id" gorm:"not null"`
	TemperatureC  float32   `json:"temperature" gorm:"not null"`
	WeightKG      float32   `json:"weight_kg" gorm:"not null"`
	TrimNailsTime time.Time `json:"trim_nails_time" gorm:"not null"`
}
