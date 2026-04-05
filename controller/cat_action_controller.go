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

// CreateCatActionRequest 创建 CatAction 请求
type CreateCatActionRequest struct {
	CatID        uint              `json:"cat_id" binding:"required,min=1"`
	SiteID       uint              `json:"site_id" binding:"required,min=1"`
	ActionType   model.CatActionType `json:"action_type" binding:"required,catActionType"`
	ActionDetail string            `json:"action_detail" binding:"required,min=1,max=1000"`
}

// CatActionController CatAction 处理器
type CatActionController struct {
	repo            *repository.CatActionRepository
	catRepo         *repository.CatRepository
	siteRepo        *repository.SiteRepository
	actionProcessor *middleware.ActionProcessor
}

// NewCatActionController 创建 CatActionController 实例
func NewCatActionController(repo *repository.CatActionRepository, catRepo *repository.CatRepository, siteRepo *repository.SiteRepository, actionProcessor *middleware.ActionProcessor) *CatActionController {
	return &CatActionController{
		repo:            repo,
		catRepo:         catRepo,
		siteRepo:        siteRepo,
		actionProcessor: actionProcessor,
	}
}

// GetCatActionsPage 分页获取 CatAction
func (ctrl *CatActionController) GetCatActionsPage(c *gin.Context) {
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

// GetCatAction 获取单个 CatAction
func (ctrl *CatActionController) GetCatAction(c *gin.Context) {
	idStr := c.Param("action_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	action, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CatAction not found"})
		return
	}
	c.JSON(http.StatusOK, action)
}

// GetCatActionsByCatID 根据 CatID 获取操作记录
func (ctrl *CatActionController) GetCatActionsByCatID(c *gin.Context) {
	catIDStr := c.Param("cat_id")
	catID, err := strconv.ParseUint(catIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CatID"})
		return
	}
	actions, err := ctrl.repo.FindByCatID(uint(catID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, actions)
}

// GetCatActionsBySiteID 根据 SiteID 获取操作记录
func (ctrl *CatActionController) GetCatActionsBySiteID(c *gin.Context) {
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

// GetCatActionsByUserID 根据 UserID 获取操作记录
func (ctrl *CatActionController) GetCatActionsByUserID(c *gin.Context) {
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

// CreateCatAction 创建 CatAction
func (ctrl *CatActionController) CreateCatAction(c *gin.Context) {
	var req CreateCatActionRequest
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

	// 校验 CatID 和 SiteID 是否存在
	var errors []string

	_, err = ctrl.catRepo.FindByID(req.CatID)
	if err != nil {
		errors = append(errors, "CatID does not exist")
	}

	_, err = ctrl.siteRepo.FindByID(req.SiteID)
	if err != nil {
		errors = append(errors, "SiteID does not exist")
	}

	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	action := model.CatAction{
		CatID:        req.CatID,
		SiteID:       req.SiteID,
		UserID:       uint(userID),
		ActionType:   req.ActionType,
		ActionDetail: req.ActionDetail,
	}

	// 使用动作处理器处理动作，自动更新状态机
	updatedFSM, err := ctrl.actionProcessor.ProcessAction(&action)
	if err != nil {
		// 判断错误类型，如果是记录不存在相关的错误，返回 400 Bad Request
		if containsRecordNotFoundError(err.Error()) {
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

// containsRecordNotFoundError 检查错误信息是否包含记录不存在的错误
func containsRecordNotFoundError(errMsg string) bool {
	return strings.Contains(errMsg, "record not found") ||
		strings.Contains(errMsg, "CatFSM not found") ||
		strings.Contains(errMsg, "CatID does not exist") ||
		strings.Contains(errMsg, "更新状态机失败")
}
