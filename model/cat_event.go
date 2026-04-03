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
	EventID   uint           `json:"event_id" gorm:"primaryKey"`
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
	CatActionTakeTemperature = "测体温"
	CatActionSterilize       = "绝育"
	CatActionHealthCheck     = "体检"
	CatActionDeworm          = "驱虫"
	CatActionTrimNails       = "修剪指甲"
	CatActionBathing         = "洗澡"
	CatActionVaccinate       = "疫苗"
)

type CatAction struct {
	ActionID     uint           `json:"action_id" gorm:"primaryKey"`
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

// SterilizeActionDetail 绝育动作的详细信息
type SterilizeActionDetail struct {
	Notes string `json:"notes" binding:"required"` // 备注
}

// HealthCheckActionDetail 体检动作的详细信息
type HealthCheckActionDetail struct {
	Temperature float32 `json:"temperature" binding:"required,gte=0,lte=50"` // 体温，单位：摄氏度
	Weight      float32 `json:"weight" binding:"required,gte=0.1,lte=25"`   // 体重，单位：千克
	Notes       string  `json:"notes" binding:"required"`                    // 备注
}

// DewormActionDetail 驱虫动作的详细信息
type DewormActionDetail struct {
	DrugName string `json:"drug_name" binding:"required"` // 药物名称
	Dosage   string `json:"dosage" binding:"required"`    // 剂量
}

// TrimNailsActionDetail 修剪指甲动作的详细信息
type TrimNailsActionDetail struct {
	Notes string `json:"notes" binding:"required"` // 备注
}

// BathingActionDetail 洗澡动作的详细信息
type BathingActionDetail struct {
	Notes string `json:"notes" binding:"required"` // 备注
}

// VaccinateActionDetail 疫苗动作的详细信息
type VaccinateActionDetail struct {
	VaccineName string `json:"vaccine_name" binding:"required"` // 疫苗名称
	BatchNo     string `json:"batch_no" binding:"required"`     // 批号
}

// ParseTemperatureActionDetail 解析测体温动作的详细信息
func ParseTemperatureActionDetail(detail string) (*TemperatureActionDetail, error) {
	var actionDetail TemperatureActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析测体温信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseSterilizeActionDetail 解析绝育动作的详细信息
func ParseSterilizeActionDetail(detail string) (*SterilizeActionDetail, error) {
	var actionDetail SterilizeActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析绝育信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseHealthCheckActionDetail 解析体检动作的详细信息
func ParseHealthCheckActionDetail(detail string) (*HealthCheckActionDetail, error) {
	var actionDetail HealthCheckActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析体检信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseDewormActionDetail 解析驱虫动作的详细信息
func ParseDewormActionDetail(detail string) (*DewormActionDetail, error) {
	var actionDetail DewormActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析驱虫信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseTrimNailsActionDetail 解析修剪指甲动作的详细信息
func ParseTrimNailsActionDetail(detail string) (*TrimNailsActionDetail, error) {
	var actionDetail TrimNailsActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析修剪指甲信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseBathingActionDetail 解析洗澡动作的详细信息
func ParseBathingActionDetail(detail string) (*BathingActionDetail, error) {
	var actionDetail BathingActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析洗澡信息失败: %v", err)
	}
	return &actionDetail, nil
}

// ParseVaccinateActionDetail 解析疫苗动作的详细信息
func ParseVaccinateActionDetail(detail string) (*VaccinateActionDetail, error) {
	var actionDetail VaccinateActionDetail
	if err := json.Unmarshal([]byte(detail), &actionDetail); err != nil {
		return nil, fmt.Errorf("解析疫苗信息失败: %v", err)
	}
	return &actionDetail, nil
}
