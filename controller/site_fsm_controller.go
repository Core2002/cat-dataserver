package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
)

// SiteFSMController SiteFSM 处理器
type SiteFSMController struct {
	repo     *repository.SiteFSMRepository
	siteRepo *repository.SiteRepository
}

// NewSiteFSMController 创建 SiteFSMController 实例
func NewSiteFSMController(repo *repository.SiteFSMRepository, siteRepo *repository.SiteRepository) *SiteFSMController {
	return &SiteFSMController{repo: repo, siteRepo: siteRepo}
}

// GetSiteFSMsPage 分页获取 SiteFSM
func (ctrl *SiteFSMController) GetSiteFSMsPage(c *gin.Context) {
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

// GetSiteFSM 获取单个 SiteFSM
func (ctrl *SiteFSMController) GetSiteFSM(c *gin.Context) {
	idStr := c.Param("site_id")
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
	idStr := c.Param("site_id")
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
	if !updates.LastDisinfectTime.IsZero() {
		fsm.LastDisinfectTime = updates.LastDisinfectTime
	}
	if !updates.LastFeedTime.IsZero(){
		fsm.LastFeedTime = updates.LastFeedTime
	}
	if !updates.LastGiveWaterTime.IsZero() {
		fsm.LastGiveWaterTime = updates.LastGiveWaterTime
	}
	if !updates.LastPlayTime.IsZero() {
		fsm.LastPlayTime = updates.LastPlayTime
	}
	if !updates.LastCleanLitter.IsZero() {
		fsm.LastCleanLitter = updates.LastCleanLitter
	}
	if err := ctrl.repo.Update(fsm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fsm)
}

// DeleteSiteFSM 删除 SiteFSM
func (ctrl *SiteFSMController) DeleteSiteFSM(c *gin.Context) {
	idStr := c.Param("site_id")
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

// getOrCreateSiteFSM 获取或创建 SiteFSM，如果不存在则自动创建
func (ctrl *SiteFSMController) getOrCreateSiteFSM(siteID uint) (*model.SiteFSM, error) {
	// 验证 Site 是否存在
	_, err := ctrl.siteRepo.FindByID(siteID)
	if err != nil {
		return nil, fmt.Errorf("site with ID %d not found", siteID)
	}
	// 获取或创建 SiteFSM（原子操作）
	return ctrl.repo.GetOrCreateBySiteID(siteID)
}

// UpdateDisinfectTime 更新消毒时间
func (ctrl *SiteFSMController) UpdateDisinfectTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	// 获取或创建 SiteFSM
	fsm, err := ctrl.getOrCreateSiteFSM(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var update struct {
		Time string `json:"last_disinfect_time"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.UpdateDisinfectTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回更新后的数据
	fsm, _ = ctrl.repo.FindBySiteID(uint(siteID))
	c.JSON(http.StatusOK, fsm)
}

// UpdateFeedTime 更新喂食时间
func (ctrl *SiteFSMController) UpdateFeedTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	// 获取或创建 SiteFSM
	fsm, err := ctrl.getOrCreateSiteFSM(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var update struct {
		Time string `json:"last_feed_time"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.UpdateFeedTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回更新后的数据
	fsm, _ = ctrl.repo.FindBySiteID(uint(siteID))
	c.JSON(http.StatusOK, fsm)
}

// UpdateGiveWaterTime 更新喂水时间
func (ctrl *SiteFSMController) UpdateGiveWaterTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	// 获取或创建 SiteFSM
	fsm, err := ctrl.getOrCreateSiteFSM(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var update struct {
		Time string `json:"last_give_water_time"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.UpdateGiveWaterTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回更新后的数据
	fsm, _ = ctrl.repo.FindBySiteID(uint(siteID))
	c.JSON(http.StatusOK, fsm)
}

// UpdatePlayTime 更新逗猫时间
func (ctrl *SiteFSMController) UpdatePlayTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	// 获取或创建 SiteFSM
	fsm, err := ctrl.getOrCreateSiteFSM(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var update struct {
		Time string `json:"last_play_time"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.UpdatePlayTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回更新后的数据
	fsm, _ = ctrl.repo.FindBySiteID(uint(siteID))
	c.JSON(http.StatusOK, fsm)
}

// UpdateCleanLitterTime 更新清理猫砂时间
func (ctrl *SiteFSMController) UpdateCleanLitterTime(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiteID"})
		return
	}
	// 获取或创建 SiteFSM
	fsm, err := ctrl.getOrCreateSiteFSM(uint(siteID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var update struct {
		Time string `json:"last_clean_litter_time"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.UpdateCleanLitterTime(uint(siteID), update.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回更新后的数据
	fsm, _ = ctrl.repo.FindBySiteID(uint(siteID))
	c.JSON(http.StatusOK, fsm)
}
