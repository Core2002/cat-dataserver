package controller

import (
	"fifu.fun/cat-dataserver/middleware"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CatActionController CatAction 处理器
type CatActionController struct {
	repo           *repository.CatActionRepository
	actionProcessor *middleware.ActionProcessor
}

// NewCatActionController 创建 CatActionController 实例
func NewCatActionController(repo *repository.CatActionRepository, actionProcessor *middleware.ActionProcessor) *CatActionController {
	return &CatActionController{
		repo:           repo,
		actionProcessor: actionProcessor,
	}
}

// GetCatActions 获取所有 CatAction
func (ctrl *CatActionController) GetCatActions(c *gin.Context) {
	actions, err := ctrl.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, actions)
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
	idStr := c.Param("id")
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
	var action model.CatAction
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用动作处理器处理动作，自动更新状态机
	updatedFSM, err := ctrl.actionProcessor.ProcessAction(&action)
	if err != nil {
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

// UpdateCatAction 更新 CatAction
func (ctrl *CatActionController) UpdateCatAction(c *gin.Context) {
	idStr := c.Param("id")
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
	var updates model.CatAction
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.Update(action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, action)
}

// DeleteCatAction 删除 CatAction
func (ctrl *CatActionController) DeleteCatAction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := ctrl.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "CatAction deleted successfully"})
}
