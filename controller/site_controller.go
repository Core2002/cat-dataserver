package controller

import (
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SiteController Site 处理器
type SiteController struct {
	repo *repository.SiteRepository
}

// NewSiteController 创建 SiteController 实例
func NewSiteController(repo *repository.SiteRepository) *SiteController {
	return &SiteController{repo: repo}
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
	idStr := c.Param("id")
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
	var site model.Site
	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.repo.Create(&site); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, site)
}

// UpdateSite 更新 Site
func (ctrl *SiteController) UpdateSite(c *gin.Context) {
	idStr := c.Param("id")
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
	var updates model.Site
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 更新字段
	if updates.SiteID != 0 {
		site.SiteID = updates.SiteID
	}
	if updates.SiteName != "" {
		site.SiteName = updates.SiteName
	}
	if updates.SiteAddress != "" {
		site.SiteAddress = updates.SiteAddress
	}
	if updates.SiteAdminPhoneNumber != "" {
		site.SiteAdminPhoneNumber = updates.SiteAdminPhoneNumber
	}
	if err := ctrl.repo.Update(site); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, site)
}

// DeleteSite 删除 Site
func (ctrl *SiteController) DeleteSite(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "Site deleted successfully"})
}
