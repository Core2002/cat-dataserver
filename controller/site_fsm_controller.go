package controller

import (
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


