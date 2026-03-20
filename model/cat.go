package model

type Cat struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:100;not null"`
}

type CatInfo struct {
	ID          uint
	Name        string
	Age         int
	Temperature float64
}
