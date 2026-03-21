package controller

import (
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CatEventController CatEvent 处理器
type CatEventController struct {
	repo *repository.CatEventRepository
}

// NewCatEventController 创建 CatEventController 实例
func NewCatEventController(repo *repository.CatEventRepository) *CatEventController {
	return &CatEventController{repo: repo}
}

// GetCatEvents 获取所有 CatEvent
func (ctrl *CatEventController) GetCatEvents(c *gin.Context) {
	events, err := ctrl.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
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
	idStr := c.Param("id")
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
	var event model.CatEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.Create(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)
}

// UpdateCatEvent 更新 CatEvent
func (ctrl *CatEventController) UpdateCatEvent(c *gin.Context) {
	idStr := c.Param("id")
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
	var updates model.CatEvent
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.Update(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

// DeleteCatEvent 删除 CatEvent
func (ctrl *CatEventController) DeleteCatEvent(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "CatEvent deleted successfully"})
}
