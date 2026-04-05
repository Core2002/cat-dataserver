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
	idStr := c.Param("cat_id")
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


