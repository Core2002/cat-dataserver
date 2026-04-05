package controller

import (
	"net/http"
	"strconv"

	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
)

// CreateSiteRequest 创建 Site 请求
type CreateSiteRequest struct {
	SiteName             string `json:"site_name" binding:"required,min=1,max=100"`
	SiteAddress          string `json:"site_address" binding:"required,min=1,max=100"`
	SiteAdminPhoneNumber string `json:"site_admin_phone_number" binding:"required"`
}

// UpdateSiteRequest 更新 Site 请求
type UpdateSiteRequest struct {
	SiteName             string `json:"site_name" binding:"omitempty,min=1,max=100"`
	SiteAddress          string `json:"site_address" binding:"omitempty,min=1,max=100"`
	SiteAdminPhoneNumber string `json:"site_admin_phone_number" binding:"omitempty"`
}

// SiteController Site 处理器
type SiteController struct {
	repo        *repository.SiteRepository
	siteFSMRepo *repository.SiteFSMRepository
}

// NewSiteController 创建 SiteController 实例
func NewSiteController(repo *repository.SiteRepository, siteFSMRepo *repository.SiteFSMRepository) *SiteController {
	return &SiteController{repo: repo, siteFSMRepo: siteFSMRepo}
}

// GetSitesPage 分页获取 Site
func (ctrl *SiteController) GetSitesPage(c *gin.Context) {
	var req model.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := req.GetPage()
	pageSize := req.GetPageSize()

	sites, total, err := ctrl.repo.FindPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := model.NewPaginationResponse(sites, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetSite 获取单个 Site
func (ctrl *SiteController) GetSite(c *gin.Context) {
	idStr := c.Param("site_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	site, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}
	c.JSON(http.StatusOK, site)
}

// CreateSite 创建 Site
func (ctrl *SiteController) CreateSite(c *gin.Context) {
	var req CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	site := model.Site{
		SiteName:             req.SiteName,
		SiteAddress:          req.SiteAddress,
		SiteAdminPhoneNumber: req.SiteAdminPhoneNumber,
	}
	if err := ctrl.repo.Create(&site); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 自动创建对应的 SiteFSM 记录
	siteFSM := &model.SiteFSM{SiteID: site.SiteID}
	if err := ctrl.siteFSMRepo.Create(siteFSM); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, site)
}

// UpdateSite 更新 Site
func (ctrl *SiteController) UpdateSite(c *gin.Context) {
	idStr := c.Param("site_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	site, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}
	var req UpdateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 更新字段
	if req.SiteName != "" {
		site.SiteName = req.SiteName
	}
	if req.SiteAddress != "" {
		site.SiteAddress = req.SiteAddress
	}
	if req.SiteAdminPhoneNumber != "" {
		site.SiteAdminPhoneNumber = req.SiteAdminPhoneNumber
	}
	if err := ctrl.repo.Update(site); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, site)
}

// DeleteSite 删除 Site
func (ctrl *SiteController) DeleteSite(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "Site deleted successfully"})
}
