package model

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SiteActionType string

const (
	SiteActionDisinfect    = "消毒"
	SiteActionFeed         = "喂食"
	SiteActionGiveWater    = "喂水"
	SiteActionPlay         = "逗猫"
	SiteActionCleanLitter  = "清理猫砂"
)

type SiteAction struct {
	ActionID     uint            `json:"action_id" gorm:"primaryKey"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `json:"deleted_at" gorm:"index"`
	SiteID       uint            `json:"site_id" gorm:"not null" binding:"required,min=1"`                          // 执行地点
	UserID       uint            `json:"user_id" gorm:"not null" binding:"omitempty,min=1"`                         // 执行人
	ActionType   SiteActionType  `json:"action_type" gorm:"size:100;not null" binding:"required,siteActionType"`    // 执行项目
	ActionDetail string          `json:"action_detail" gorm:"size:1000;not null" binding:"required,min=1,max=1000"` // 事件详情
}

// TableName 指定表名
func (SiteAction) TableName() string {
	return "site_actions"
}

// DisinfectActionDetail 消毒动作的详细信息
type DisinfectActionDetail struct {
	Disinfectant string `json:"disinfectant" binding:"required"` // 消毒剂名称
	Notes        string `json:"notes"`                           // 备注
}

// FeedActionDetail 喂食动作的详细信息
type FeedActionDetail struct {
	FoodType string `json:"food_type" binding:"required"` // 食物类型
	Amount   string `json:"amount" binding:"required"`    // 食物量
	Notes    string `json:"notes"`                        // 备注
}

// GiveWaterActionDetail 喂水动作的详细信息
type GiveWaterActionDetail struct {
	WaterType string `json:"water_type" binding:"required"` // 水类型
	Notes     string `json:"notes"`                         // 备注
}

// PlayActionDetail 逗猫动作的详细信息
type PlayActionDetail struct {
	Duration int    `json:"duration" binding:"required,gte=1"` // 持续时间（分钟）
	Notes    string `json:"notes"`                             // 备注
}

// CleanLitterActionDetail 清理猫砂动作的详细信息
type CleanLitterActionDetail struct {
	LitterType string `json:"litter_type" binding:"required"` // 猫砂类型
	Notes      string `json:"notes"`                          // 备注
}

// ParseDisinfectActionDetail 解析消毒动作的详细信息
func ParseDisinfectActionDetail(detail string) (*DisinfectActionDetail, error) {
	var actionDetail DisinfectActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析消毒信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseFeedActionDetail 解析喂食动作的详细信息
func ParseFeedActionDetail(detail string) (*FeedActionDetail, error) {
	var actionDetail FeedActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析喂食信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseGiveWaterActionDetail 解析喂水动作的详细信息
func ParseGiveWaterActionDetail(detail string) (*GiveWaterActionDetail, error) {
	var actionDetail GiveWaterActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析喂水信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParsePlayActionDetail 解析逗猫动作的详细信息
func ParsePlayActionDetail(detail string) (*PlayActionDetail, error) {
	var actionDetail PlayActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析逗猫信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseCleanLitterActionDetail 解析清理猫砂动作的详细信息
func ParseCleanLitterActionDetail(detail string) (*CleanLitterActionDetail, error) {
	var actionDetail CleanLitterActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析清理猫砂信息失败: %v", err)
	}
	return &actionDetail, nil
}
