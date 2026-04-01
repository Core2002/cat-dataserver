package middleware

import (
	"fmt"
	"log"
	"time"

	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
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
		// 状态机更新失败，需要回滚已创建的动作记录
		if deleteErr := p.actionRepo.Delete(action.ID); deleteErr != nil {
			return nil, fmt.Errorf("更新状态机失败: %v, 回滚动作记录失败: %v", err, deleteErr)
		}
		return nil, fmt.Errorf("更新状态机失败: %v", err)
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
	case model.CatActionHealthCheck:
		return p.updateWeight(action, fsm)
	case model.CatActionSterilize, model.CatActionDeworm, model.CatActionVaccinate:
		return fsm, nil // 医疗类动作暂不更新FSM
	default:
		return fsm, nil
	}
}

// updateTemperature 更新体温
func (p *ActionProcessor) updateTemperature(action *model.CatAction, fsm *model.CatFSM) (*model.CatFSM, error) {
	detail, err := model.ParseTemperatureActionDetail(action.ActionDetail)
	if err != nil {
		return nil, fmt.Errorf("解析测体温信息失败: %v", err)
	}

	if err := p.fsmRepo.UpdateTemperature(action.CatID, detail.Temperature); err != nil {
		return nil, fmt.Errorf("更新体温失败: %v", err)
	}

	fsm.TemperatureC = detail.Temperature
	return fsm, nil
}

// updateWeight 更新体重
func (p *ActionProcessor) updateWeight(action *model.CatAction, fsm *model.CatFSM) (*model.CatFSM, error) {
	detail, err := model.ParseWeightActionDetail(action.ActionDetail)
	if err != nil {
		return nil, fmt.Errorf("解析测体重信息失败: %v", err)
	}

	if err := p.fsmRepo.UpdateWeight(action.CatID, detail.Weight); err != nil {
		return nil, fmt.Errorf("更新体重失败: %v", err)
	}

	fsm.WeightKG = detail.Weight
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
