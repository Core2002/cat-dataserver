package controller

import (
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SiteFSMController SiteFSM 处理器
type SiteFSMController struct {
	repo *repository.SiteFSMRepository
}

// NewSiteFSMController 创建 SiteFSMController 实例
func NewSiteFSMController(repo *repository.SiteFSMRepository) *SiteFSMController {
	return &SiteFSMController{repo: repo}
}

// GetSiteFSM 获取单个 SiteFSM
func (ctrl *SiteFSMController) GetSiteFSM(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	fsm, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SiteFSM not found"})
		return
	}
	c.JSON(http.StatusOK, fsm)
}

// GetSiteFSMBySiteID 根据 SiteID 获取 SiteFSM
func (ctrl *SiteFSMController) GetSiteFSMBySiteID(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	fsm, err := ctrl.repo.FindBySiteID(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SiteFSM not found"})
		return
	}
	c.JSON(http.StatusOK, fsm)
}

// CreateSiteFSM 创建 SiteFSM
func (ctrl *SiteFSMController) CreateSiteFSM(c *gin.Context) {
	var fsm model.SiteFSM
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

// UpdateSiteFSM 更新 SiteFSM
func (ctrl *SiteFSMController) UpdateSiteFSM(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	fsm, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SiteFSM not found"})
		return
	}
	var updates model.SiteFSM
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 更新字段
	if updates.SiteID != 0 {
		fsm.SiteID = updates.SiteID
	}
	if updates.LastDisinfectTime != nil {
		fsm.LastDisinfectTime = updates.LastDisinfectTime
	}
	if updates.LastFeedTime != nil {
		fsm.LastFeedTime = updates.LastFeedTime
	}
	if updates.LastGiveWaterTime != nil {
		fsm.LastGiveWaterTime = updates.LastGiveWaterTime
	}
	if updates.LastPlayTime != nil {
		fsm.LastPlayTime = updates.LastPlayTime
	}
	if err := ctrl.repo.Update(fsm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fsm)
}

// DeleteSiteFSM 删除 SiteFSM
func (ctrl *SiteFSMController) DeleteSiteFSM(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "SiteFSM deleted successfully"})
}

// UpdateDisinfectTime 更新消毒时间
func (ctrl *SiteFSMController) UpdateDisinfectTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
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
	if err := ctrl.repo.UpdateDisinfectTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Disinfect time updated successfully"})
}

// UpdateFeedTime 更新喂食时间
func (ctrl *SiteFSMController) UpdateFeedTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
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
	if err := ctrl.repo.UpdateFeedTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Feed time updated successfully"})
}

// UpdateGiveWaterTime 更新喂水时间
func (ctrl *SiteFSMController) UpdateGiveWaterTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
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
	if err := ctrl.repo.UpdateGiveWaterTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Give water time updated successfully"})
}

// UpdatePlayTime 更新逗猫时间
func (ctrl *SiteFSMController) UpdatePlayTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
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
	if err := ctrl.repo.UpdatePlayTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Play time updated successfully"})
}
