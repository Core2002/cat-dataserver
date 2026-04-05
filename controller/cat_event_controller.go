package controller

import (
	"net/http"
	"strconv"

	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
)

// CreateCatEventRequest 创建 CatEvent 请求
type CreateCatEventRequest struct {
	EventType model.CatEventType `json:"event_type" binding:"required,catEventType"`
	SiteID    uint               `json:"site_id" binding:"required,min=1"`
	CatID     uint               `json:"cat_id" binding:"required,min=1"`
	Detail    string             `json:"detail" binding:"required,min=1,max=1000"`
}

// UpdateCatEventRequest 更新 CatEvent 请求
type UpdateCatEventRequest struct {
	EventType model.CatEventType `json:"event_type" binding:"omitempty,catEventType"`
	SiteID    uint               `json:"site_id" binding:"omitempty,min=1"`
	UserID    uint               `json:"user_id" binding:"omitempty,min=1"`
	CatID     uint               `json:"cat_id" binding:"omitempty,min=1"`
	Detail    string             `json:"detail" binding:"omitempty,min=1,max=1000"`
}

// CatEventController CatEvent 处理器
type CatEventController struct {
	repo     *repository.CatEventRepository
	catRepo  *repository.CatRepository
	siteRepo *repository.SiteRepository
}

// NewCatEventController 创建 CatEventController 实例
func NewCatEventController(repo *repository.CatEventRepository, catRepo *repository.CatRepository, siteRepo *repository.SiteRepository) *CatEventController {
	return &CatEventController{
		repo:     repo,
		catRepo:  catRepo,
		siteRepo: siteRepo,
	}
}

// GetCatEventsPage 分页获取 CatEvent
func (ctrl *CatEventController) GetCatEventsPage(c *gin.Context) {
	var req model.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := req.GetPage()
	pageSize := req.GetPageSize()

	events, total, err := ctrl.repo.FindPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := model.NewPaginationResponse(events, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetCatEvent 获取单个 CatEvent
func (ctrl *CatEventController) GetCatEvent(c *gin.Context) {
	idStr := c.Param("event_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	event, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CatEvent not found"})
		return
	}
	c.JSON(http.StatusOK, event)
}

// GetCatEventsByCatID 根据 CatID 获取事件
func (ctrl *CatEventController) GetCatEventsByCatID(c *gin.Context) {
	catIDStr := c.Param("cat_id")
	catID, err := strconv.ParseUint(catIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CatID"})
		return
	}
	events, err := ctrl.repo.FindByCatID(uint(catID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// GetCatEventsBySiteID 根据 SiteID 获取事件
func (ctrl *CatEventController) GetCatEventsBySiteID(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	events, err := ctrl.repo.FindBySiteID(uint(siteID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// CreateCatEvent 创建 CatEvent
func (ctrl *CatEventController) CreateCatEvent(c *gin.Context) {
	var req CreateCatEventRequest
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

	event := model.CatEvent{
		EventType: req.EventType,
		SiteID:    req.SiteID,
		UserID:    uint(userID),
		CatID:     req.CatID,
		Detail:    req.Detail,
	}

	if err := ctrl.repo.Create(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"event": event,
	})
}

// UpdateCatEvent 更新 CatEvent
func (ctrl *CatEventController) UpdateCatEvent(c *gin.Context) {
	idStr := c.Param("cat_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	event, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CatEvent not found"})
		return
	}
	var req UpdateCatEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 更新数据
	if req.EventType != "" {
		event.EventType = req.EventType
	}
	if req.SiteID != 0 {
		event.SiteID = req.SiteID
	}
	if req.UserID != 0 {
		event.UserID = req.UserID
	}
	if req.CatID != 0 {
		event.CatID = req.CatID
	}
	if req.Detail != "" {
		event.Detail = req.Detail
	}
	if err := ctrl.repo.Update(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

// DeleteCatEvent 删除 CatEvent
func (ctrl *CatEventController) DeleteCatEvent(c *gin.Context) {
	idStr := c.Param("cat_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := ctrl.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "CatEvent deleted successfully"})
}
