package model

import (
	"gorm.io/gorm"
)

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
	gorm.Model
	EventID   uint         `json:"event_id" gorm:"primaryKey" binding:"min=1"`
	EventType CatEventType `json:"event_type" gorm:"size:100;not null" binding:"required,catEventType"` // 事件项目
	SiteID    uint         `json:"site_id" gorm:"not null" binding:"required,min=1"`                    // 事件地点
	CatID     uint         `json:"cat_id" gorm:"not null" binding:"required,min=1"`                     // 发生事件的猫
	Detail    string       `json:"detail" gorm:"size:1000;not null" binding:"required,min=1,max=1000"`  // 事件详情
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
	gorm.Model
	ActionID     uint          `json:"action_id" gorm:"primaryKey" binding:"min=1"`
	CatID        uint          `json:"cat_id" gorm:"not null" binding:"required,min=1"`                           // 被执行的猫
	SiteID       uint          `json:"site_id" gorm:"not null" binding:"required,min=1"`                          // 执行地点
	UserID       uint          `json:"user_id" gorm:"not null" binding:"required,min=1"`                          // 执行人
	ActionType   CatActionType `json:"action_type" gorm:"size:100;not null" binding:"required,catActionType"`     // 执行项目
	ActionDetail string        `json:"action_detail" gorm:"size:1000;not null" binding:"required,min=1,max=1000"` // 事件详情
}
