package middleware

import (
	"fmt"
	"log"
	"time"

	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
)

// SiteActionProcessor 站点动作处理器中间件，负责自动处理动作并更新状态机
type SiteActionProcessor struct {
	actionRepo *repository.SiteActionRepository
	fsmRepo    *repository.SiteFSMRepository
}

// NewSiteActionProcessor 创建站点动作处理器实例
func NewSiteActionProcessor(
	actionRepo *repository.SiteActionRepository,
	fsmRepo *repository.SiteFSMRepository,
) *SiteActionProcessor {
	return &SiteActionProcessor{
		actionRepo: actionRepo,
		fsmRepo:    fsmRepo,
	}
}

// ProcessAction 处理动作并更新状态机
// 在创建动作记录后调用此方法，会自动更新相应的状态机
func (p *SiteActionProcessor) ProcessAction(action *model.SiteAction) (*model.SiteFSM, error) {
	// 1. 记录动作到数据库
	if err := p.actionRepo.Create(action); err != nil {
		return nil, err
	}

	// 2. 根据动作类型更新状态机
	fsm, err := p.updateFSM(action)
	if err != nil {
		// 状态机更新失败，需要回滚已创建的动作记录
		if deleteErr := p.actionRepo.Delete(action.ActionID); deleteErr != nil {
			return nil, fmt.Errorf("更新状态机失败: %v, 回滚动作记录失败: %v", err, deleteErr)
		}
		return nil, fmt.Errorf("更新状态机失败: %v", err)
	}

	log.Printf("站点动作处理完成: SiteID=%d, ActionType=%s", action.SiteID, action.ActionType)
	return fsm, nil
}

// updateFSM 根据动作类型更新状态机
func (p *SiteActionProcessor) updateFSM(action *model.SiteAction) (*model.SiteFSM, error) {
	// 查找或创建状态机记录
	fsm, err := p.fsmRepo.GetOrCreateBySiteID(action.SiteID)
	if err != nil {
		return nil, fmt.Errorf("获取或创建 SiteFSM 失败: %v", err)
	}

	now := time.Now()

	switch action.ActionType {
	case model.SiteActionDisinfect:
		return p.updateDisinfectTime(action, fsm, now)
	case model.SiteActionFeed:
		return p.updateFeedTime(action, fsm, now)
	case model.SiteActionGiveWater:
		return p.updateGiveWaterTime(action, fsm, now)
	case model.SiteActionPlay:
		return p.updatePlayTime(action, fsm, now)
	case model.SiteActionCleanLitter:
		return p.updateCleanLitterTime(action, fsm, now)
	default:
		return fsm, nil
	}
}

// updateDisinfectTime 更新消毒时间
func (p *SiteActionProcessor) updateDisinfectTime(action *model.SiteAction, fsm *model.SiteFSM, now time.Time) (*model.SiteFSM, error) {
	if err := p.fsmRepo.UpdateDisinfectTime(action.SiteID, now.Format(time.RFC3339)); err != nil {
		return nil, fmt.Errorf("更新消毒时间失败: %v", err)
	}

	fsm.LastDisinfectTime = &now
	return fsm, nil
}

// updateFeedTime 更新喂食时间
func (p *SiteActionProcessor) updateFeedTime(action *model.SiteAction, fsm *model.SiteFSM, now time.Time) (*model.SiteFSM, error) {
	if err := p.fsmRepo.UpdateFeedTime(action.SiteID, now.Format(time.RFC3339)); err != nil {
		return nil, fmt.Errorf("更新喂食时间失败: %v", err)
	}

	fsm.LastFeedTime = &now
	return fsm, nil
}

// updateGiveWaterTime 更新喂水时间
func (p *SiteActionProcessor) updateGiveWaterTime(action *model.SiteAction, fsm *model.SiteFSM, now time.Time) (*model.SiteFSM, error) {
	if err := p.fsmRepo.UpdateGiveWaterTime(action.SiteID, now.Format(time.RFC3339)); err != nil {
		return nil, fmt.Errorf("更新喂水时间失败: %v", err)
	}

	fsm.LastGiveWaterTime = &now
	return fsm, nil
}

// updatePlayTime 更新逗猫时间
func (p *SiteActionProcessor) updatePlayTime(action *model.SiteAction, fsm *model.SiteFSM, now time.Time) (*model.SiteFSM, error) {
	if err := p.fsmRepo.UpdatePlayTime(action.SiteID, now.Format(time.RFC3339)); err != nil {
		return nil, fmt.Errorf("更新逗猫时间失败: %v", err)
	}

	fsm.LastPlayTime = &now
	return fsm, nil
}

// updateCleanLitterTime 更新清理猫砂时间
func (p *SiteActionProcessor) updateCleanLitterTime(action *model.SiteAction, fsm *model.SiteFSM, now time.Time) (*model.SiteFSM, error) {
	if err := p.fsmRepo.UpdateCleanLitterTime(action.SiteID, now.Format(time.RFC3339)); err != nil {
		return nil, fmt.Errorf("更新清理猫砂时间失败: %v", err)
	}

	fsm.LastCleanLitter = &now
	return fsm, nil
}
