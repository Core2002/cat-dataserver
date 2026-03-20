package model

type CatEventType string

const (
	CatEventTypeSick               = "生病"
	CatEventTypeInjure             = "受伤"
	CatEventTypePreg               = "怀孕"
	CatEventTypeBirth              = "分娩"
	CatEventTypeDeath              = "死亡"
	CatEventTypeContractTerminatio = "合同解除"
)

type CatEvent struct {
	ID     uint         `json:"id" gorm:"primaryKey"`
	CatID  uint         `json:"cat_id" gorm:"not null"`
	Type   CatEventType `json:"type" gorm:"size:100;not null"`
	Detail string       `json:"detail" gorm:"size:1000;not null"` // 事件详情
	Time   string       `json:"time" gorm:"size:100;not null"`    // 发生时间
}

type CatActionType string

const (
	CatActionTypeFeed            = "喂食"
	CatActionTypeGiveWater       = "喂水"
	CatActionTypeTakeTemperature = "测体温"
	CatActionTypePlay            = "逗猫"
	CatActionTypeSterilize       = "绝育"
	CatActionTypeHealthCheck     = "体检"
	CatActionTypeDeworm          = "驱虫"
	CatActionTypeCleanLitter     = "清理猫砂"
	CatActionTypeDisinfect       = "环境消毒"
	CatActionTypeTrimNails       = "修剪指甲"
	CatActionTypeWashFeet        = "洗脚"
	CatActionTypeVaccinate       = "疫苗"
)

type CatAction struct {
	ID     uint          `json:"id" gorm:"primaryKey"`
	CatID  uint          `json:"cat_id" gorm:"not null"`
	Type   CatActionType `json:"type" gorm:"size:100;not null"`
	Time   string        `json:"time" gorm:"size:100;not null"` // 执行时间
	UserID uint          `json:"user_id" gorm:"not null"`       // 执行人
}
