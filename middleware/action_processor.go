package middleware

import (
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"fmt"
	"log"
	"time"
)

// ActionProcessor 动作处理器中间件，负责自动处理动作并更新状态机
type ActionProcessor struct {
	actionRepo *repository.CatActionRepository
	fsmRepo    *repository.CatFSMRepository
}

// NewActionProcessor 创建动作处理器实例
func NewActionProcessor(
	actionRepo *repository.CatActionRepository,
	fsmRepo *repository.CatFSMRepository,
) *ActionProcessor {
	return &ActionProcessor{
		actionRepo: actionRepo,
		fsmRepo:    fsmRepo,
	}
}

// ProcessAction 处理动作并更新状态机
// 在创建动作记录后调用此方法，会自动更新相应的状态机
func (p *ActionProcessor) ProcessAction(action *model.CatAction) (*model.CatFSM, error) {
	// 1. 记录动作到数据库
	if err := p.actionRepo.Create(action); err != nil {
		return nil, err
	}

	// 2. 根据动作类型更新状态机
	fsm, err := p.updateFSM(action)
	if err != nil {
		log.Printf("更新状态机失败: CatID=%d, ActionType=%s, Error=%v",
			action.CatID, action.ActionType, err)
		// 即使更新失败，动作记录已经保存，所以返回成功
		return nil, nil
	}

	log.Printf("动作处理完成: CatID=%d, ActionType=%s", action.CatID, action.ActionType)
	return fsm, nil
}

// updateFSM 根据动作类型更新状态机
func (p *ActionProcessor) updateFSM(action *model.CatAction) (*model.CatFSM, error) {
	// 查找状态机记录
	fsm, err := p.fsmRepo.FindByID(action.CatID)
	if err != nil {
		return nil, err
	}

	switch action.ActionType {
	case model.CatActionTakeTemperature:
		return p.updateTemperature(action, fsm)
	case model.CatActionTrimNails:
		return p.updateTrimNailsTime(action, fsm)
	case model.CatActionGiveWater, model.CatActionFeed, model.CatActionPlay:
		// 喂食、喂水、逗猫不需要更新FSM
		return fsm, nil
	case model.CatActionHealthCheck:
		return p.updateWeight(action, fsm)
	case model.CatActionSterilize, model.CatActionDeworm, model.CatActionVaccinate:
		return fsm, nil // 医疗类动作暂不更新FSM
	case model.CatActionCleanLitter, model.CatActionDisinfect, model.CatActionWashFeet:
		return fsm, nil // 卫生类动作暂不更新FSM
	default:
		return fsm, nil
	}
}

// updateTemperature 更新体温
func (p *ActionProcessor) updateTemperature(action *model.CatAction, fsm *model.CatFSM) (*model.CatFSM, error) {
	temperature := parseTemperature(action.ActionDetail)
	if temperature == 0 {
		return fsm, nil // 解析失败则不更新
	}

	if err := p.fsmRepo.UpdateTemperature(action.CatID, temperature); err != nil {
		return nil, err
	}

	fsm.TemperatureC = temperature
	return fsm, nil
}

// updateWeight 更新体重
func (p *ActionProcessor) updateWeight(action *model.CatAction, fsm *model.CatFSM) (*model.CatFSM, error) {
	weight := parseWeight(action.ActionDetail)
	if weight == 0 {
		return fsm, nil // 解析失败则不更新
	}

	if err := p.fsmRepo.UpdateWeight(action.CatID, weight); err != nil {
		return nil, err
	}

	fsm.WeightKG = weight
	return fsm, nil
}

// updateTrimNailsTime 更新修剪指甲时间
func (p *ActionProcessor) updateTrimNailsTime(action *model.CatAction, fsm *model.CatFSM) (*model.CatFSM, error) {
	now := time.Now()
	if err := p.fsmRepo.UpdateTrimNailsTime(action.CatID, now); err != nil {
		return nil, err
	}

	fsm.TrimNailsTime = now
	return fsm, nil
}

// parseTemperature 从详情中解析体温值
func parseTemperature(detail string) float32 {
	// 简化版本，提取数字部分
	var temp float32
	_, _ = fmt.Sscanf(detail, "%f", &temp)
	return temp
}

// parseWeight 从详情中解析体重值
func parseWeight(detail string) float32 {
	// 简化版本，提取数字部分
	var weight float32
	_, _ = fmt.Sscanf(detail, "%f", &weight)
	return weight
}
