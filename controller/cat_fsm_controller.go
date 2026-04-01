package controller

import (
	"net/http"
	"strconv"

	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
)

// CatFSMController CatFSM 处理器
type CatFSMController struct {
	repo *repository.CatFSMRepository
}

// NewCatFSMController 创建 CatFSMController 实例
func NewCatFSMController(repo *repository.CatFSMRepository) *CatFSMController {
	return &CatFSMController{repo: repo}
}

// GetCatFSMsPage 分页获取 CatFSM
func (ctrl *CatFSMController) GetCatFSMsPage(c *gin.Context) {
	var req model.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := req.GetPage()
	pageSize := req.GetPageSize()

	fsms, total, err := ctrl.repo.FindPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := model.NewPaginationResponse(fsms, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetCatFSM 获取单个 CatFSM
func (ctrl *CatFSMController) GetCatFSM(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	fsm, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CatFSM not found"})
		return
	}
	c.JSON(http.StatusOK, fsm)
}

// GetCatFSMsBySiteID 根据 SiteID 获取猫状态
func (ctrl *CatFSMController) GetCatFSMsBySiteID(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	fsms, err := ctrl.repo.FindBySiteID(uint(siteID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fsms)
}

// CreateCatFSM 创建 CatFSM
func (ctrl *CatFSMController) CreateCatFSM(c *gin.Context) {
	var fsm model.CatFSM
	if err := c.ShouldBindJSON(&fsm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.Create(&fsm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, fsm)
}

// UpdateCatFSM 更新 CatFSM
func (ctrl *CatFSMController) UpdateCatFSM(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	fsm, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CatFSM not found"})
		return
	}
	var updates model.CatFSM
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.Update(fsm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fsm)
}

// DeleteCatFSM 删除 CatFSM
func (ctrl *CatFSMController) DeleteCatFSM(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "CatFSM deleted successfully"})
}

// UpdateTemperature 更新体温
func (ctrl *CatFSMController) UpdateTemperature(c *gin.Context) {
	catIDStr := c.Param("cat_id")
	catID, err := strconv.ParseUint(catIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CatID"})
		return
	}
	type TemperatureUpdate struct {
		Temperature float32 `json:"temperature"`
	}
	var update TemperatureUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.UpdateTemperature(uint(catID), update.Temperature); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Temperature updated successfully"})
}

// UpdateWeight 更新体重
func (ctrl *CatFSMController) UpdateWeight(c *gin.Context) {
	catIDStr := c.Param("cat_id")
	catID, err := strconv.ParseUint(catIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CatID"})
		return
	}
	type WeightUpdate struct {
		Weight float32 `json:"weight"`
	}
	var update WeightUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.UpdateWeight(uint(catID), update.Weight); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Weight updated successfully"})
}

// UpdateTrimNailsTime 更新修剪指甲时间
func (ctrl *CatFSMController) UpdateTrimNailsTime(c *gin.Context) {
	catIDStr := c.Param("cat_id")
	catID, err := strconv.ParseUint(catIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CatID"})
		return
	}
	type TimeUpdate struct {
		Time string `json:"time"`
	}
	var update TimeUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.UpdateTrimNailsTime(uint(catID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Trim nails time updated successfully"})
}
