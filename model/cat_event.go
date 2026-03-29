package model

import (
	"encoding/json"
	"fmt"
	"time"

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
	ID        uint           `json:"event_id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	EventType CatEventType   `json:"event_type" gorm:"size:100;not null" binding:"required,catEventType"` // 事件项目
	SiteID    uint           `json:"site_id" gorm:"not null" binding:"required,min=1"`                    // 事件地点
	UserID    uint           `json:"user_id" gorm:"not null" binding:"omitempty,min=1"`                   // 记录人
	CatID     uint           `json:"cat_id" gorm:"not null" binding:"required,min=1"`                     // 发生事件的猫
	Detail    string         `json:"detail" gorm:"size:1000;not null" binding:"required,min=1,max=1000"`  // 事件详情
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
	ID           uint           `json:"action_id" gorm:"primaryKey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CatID        uint           `json:"cat_id" gorm:"not null" binding:"required,min=1"`                           // 被执行的猫
	SiteID       uint           `json:"site_id" gorm:"not null" binding:"required,min=1"`                          // 执行地点
	UserID       uint           `json:"user_id" gorm:"not null" binding:"omitempty,min=1"`                         // 执行人
	ActionType   CatActionType  `json:"action_type" gorm:"size:100;not null" binding:"required,catActionType"`     // 执行项目
	ActionDetail string         `json:"action_detail" gorm:"size:1000;not null" binding:"required,min=1,max=1000"` // 事件详情
}

// TemperatureActionDetail 测体温动作的详细信息
type TemperatureActionDetail struct {
	Temperature float32 `json:"temperature" binding:"required,gte=0,lte=50"` // 体温，单位：摄氏度
}

// WeightActionDetail 测体重动作的详细信息
type WeightActionDetail struct {
	Weight float32 `json:"weight" binding:"required,gte=0.1,lte=25"` // 体重，单位：千克
}

// ParseTemperatureActionDetail 解析测体温动作的详细信息
func ParseTemperatureActionDetail(detail string) (*TemperatureActionDetail, error) {
	var actionDetail TemperatureActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析测体温信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseWeightActionDetail 解析测体重动作的详细信息
func ParseWeightActionDetail(detail string) (*WeightActionDetail, error) {
	var actionDetail WeightActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析测体重信息失败: %v", err)
	}
	return &actionDetail, nil
}
