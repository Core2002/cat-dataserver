package model

import "time"

type CatEventType string

const (
	CatSick               = "生病"
	CatInjure             = "受伤"
	CatPreg               = "怀孕"
	CatBirth              = "分娩"
	CatDeath              = "死亡"
	CatContractTerminatio = "合同解除"
)

type CatEvent struct {
	ID        uint         `json:"id" gorm:"primaryKey"`
	CatID     uint         `json:"cat_id" gorm:"not null"`
	Type      CatEventType `json:"event_type" gorm:"size:100;not null"` // 事件项目
	Detail    string       `json:"detail" gorm:"size:1000;not null"`    // 事件详情
	EventTime time.Time    `json:"event_time" gorm:"not null"`          // 发生时间
}

type CatActionType string

const (
	CatActionFeed            = "喂食"
	CatActionGiveWater       = "喂水"
	CatActionTakeTemperature = "测体温"
	CatActionPlay            = "逗猫"
	CatActionSterilize       = "绝育"
	CatActionHealthCheck     = "体检"
	CatActionDeworm          = "驱虫"
	CatActionCleanLitter     = "清理猫砂"
	CatActionDisinfect       = "环境消毒"
	CatActionTrimNails       = "修剪指甲"
	CatActionWashFeet        = "洗脚"
	CatActionVaccinate       = "疫苗"
)

type CatAction struct {
	ID     uint          `json:"id" gorm:"primaryKey"`
	CatID  uint          `json:"cat_id" gorm:"not null"`
	UserID uint          `json:"user_id" gorm:"not null"`          // 执行人
	Time   time.Time     `json:"time" gorm:"not null"`             // 执行时间
	Type   CatActionType `json:"type" gorm:"size:100;not null"`    // 执行项目
	Detail string        `json:"detail" gorm:"size:1000;not null"` // 事件详情
}
