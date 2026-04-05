package controller

import (
	"net/http"
	"strconv"
	"strings"

	"fifu.fun/cat-dataserver/middleware"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
)

// CreateSiteActionRequest 创建 SiteAction 请求
type CreateSiteActionRequest struct {
	SiteID       uint               `json:"site_id" binding:"required,min=1"`
	ActionType   model.SiteActionType `json:"action_type" binding:"required,siteActionType"`
	ActionDetail string             `json:"action_detail" binding:"required,min=1,max=1000"`
}

// SiteActionController SiteAction 处理器
type SiteActionController struct {
	repo                *repository.SiteActionRepository
	siteRepo            *repository.SiteRepository
	siteActionProcessor *middleware.SiteActionProcessor
}

// NewSiteActionController 创建 SiteActionController 实例
func NewSiteActionController(repo *repository.SiteActionRepository, siteRepo *repository.SiteRepository, siteActionProcessor *middleware.SiteActionProcessor) *SiteActionController {
	return &SiteActionController{
		repo:               repo,
		siteRepo:           siteRepo,
		siteActionProcessor: siteActionProcessor,
	}
}

// GetSiteActionsPage 分页获取 SiteAction
func (ctrl *SiteActionController) GetSiteActionsPage(c *gin.Context) {
	var req model.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := req.GetPage()
	pageSize := req.GetPageSize()

	actions, total, err := ctrl.repo.FindPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := model.NewPaginationResponse(actions, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetSiteAction 获取单个 SiteAction
func (ctrl *SiteActionController) GetSiteAction(c *gin.Context) {
	idStr := c.Param("action_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	action, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SiteAction not found"})
		return
	}
	c.JSON(http.StatusOK, action)
}

// GetSiteActionsBySiteID 根据 SiteID 获取操作记录
func (ctrl *SiteActionController) GetSiteActionsBySiteID(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	actions, err := ctrl.repo.FindBySiteID(uint(siteID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, actions)
}

// GetSiteActionsByUserID 根据 UserID 获取操作记录
func (ctrl *SiteActionController) GetSiteActionsByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
		return
	}
	actions, err := ctrl.repo.FindByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, actions)
}

// CreateSiteAction 创建 SiteAction
func (ctrl *SiteActionController) CreateSiteAction(c *gin.Context) {
	var req CreateSiteActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从请求头获取 UserID
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-User-ID header"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid X-User-ID header"})
		return
	}

	// 校验 SiteID 是否存在
	_, err = ctrl.siteRepo.FindByID(req.SiteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SiteID does not exist"})
		return
	}

	action := model.SiteAction{
		SiteID:       req.SiteID,
		UserID:       uint(userID),
		ActionType:   req.ActionType,
		ActionDetail: req.ActionDetail,
	}

	// 使用动作处理器处理动作，自动更新状态机
	updatedFSM, err := ctrl.siteActionProcessor.ProcessAction(&action)
	if err != nil {
		// 判断错误类型，如果是记录不存在相关的错误，返回 400 Bad Request
		if containsSiteRecordNotFoundError(err.Error()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"action": action,
	}

	// 如果状态机被更新，返回更新后的状态
	if updatedFSM != nil {
		response["fsm"] = updatedFSM
	}

	c.JSON(http.StatusCreated, response)
}

// containsSiteRecordNotFoundError 检查错误信息是否包含记录不存在的错误
func containsSiteRecordNotFoundError(errMsg string) bool {
	return strings.Contains(errMsg, "record not found") ||
		strings.Contains(errMsg, "SiteFSM not found") ||
		strings.Contains(errMsg, "SiteID does not exist") ||
		strings.Contains(errMsg, "更新状态机失败")
}
