package controller

import (
	"net/http"
	"strconv"
	"time"

	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
)

// CreateCatRequest 创建猫请求
type CreateCatRequest struct {
	CatName           string `json:"cat_name" binding:"required,min=1,max=100"`
	CatPhotoUri       string `json:"cat_photo_uri" binding:"required,url"`
	CatType           string `json:"cat_type" binding:"required,min=1,max=50"`
	CatGender         string `json:"cat_gender" binding:"required"`
	MasterName        string `json:"master_name" binding:"required,min=1,max=100"`
	MasterPhoneNumber string `json:"master_phone_number" binding:"required"`
	SiteID            uint   `json:"site_id" binding:"omitempty,min=1"`
}

// UpdateCatRequest 更新猫请求
type UpdateCatRequest struct {
	CatName           string `json:"cat_name" binding:"omitempty,min=1,max=100"`
	CatPhotoUri       string `json:"cat_photo_uri" binding:"omitempty,url"`
	CatType           string `json:"cat_type" binding:"omitempty,min=1,max=50"`
	CatGender         string `json:"cat_gender" binding:"omitempty"`
	MasterName        string `json:"master_name" binding:"omitempty,min=1,max=100"`
	MasterPhoneNumber string `json:"master_phone_number" binding:"omitempty"`
}

// CatController Cat 处理器
type CatController struct {
	repo      *repository.CatRepository
	fsmRepo   *repository.CatFSMRepository
	siteRepo  *repository.SiteRepository
}

// NewCatController 创建 CatController 实例
func NewCatController(repo *repository.CatRepository, fsmRepo *repository.CatFSMRepository, siteRepo *repository.SiteRepository) *CatController {
	return &CatController{repo: repo, fsmRepo: fsmRepo, siteRepo: siteRepo}
}

// GetCatsPage 分页获取 Cat
func (ctrl *CatController) GetCatsPage(c *gin.Context) {
	var req model.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := req.GetPage()
	pageSize := req.GetPageSize()

	cats, total, err := ctrl.repo.FindPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := model.NewPaginationResponse(cats, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// GetCat 获取单个 Cat
func (ctrl *CatController) GetCat(c *gin.Context) {
	idStr := c.Param("cat_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	cat, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// CreateCat 创建 Cat
func (ctrl *CatController) CreateCat(c *gin.Context) {
	var req CreateCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证 SiteID 是否存在（仅当提供 SiteID 时）
	if req.SiteID > 0 {
		_, err := ctrl.siteRepo.FindByID(req.SiteID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SiteID does not exist"})
			return
		}
	}

	// 创建猫
	cat := model.Cat{
		CatName:           req.CatName,
		CatPhotoUri:       req.CatPhotoUri,
		CatType:           req.CatType,
		CatGender:         req.CatGender,
		MasterName:        req.MasterName,
		MasterPhoneNumber: req.MasterPhoneNumber,
	}
	if err := ctrl.repo.Create(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 自动创建 CatFSM 记录
	fsm := model.CatFSM{
		CatID:         cat.CatID,
		SiteID:        req.SiteID,
		TemperatureC:  37.5, // 默认体温
		WeightKG:      4.0,  // 默认体重
		TrimNailsTime: time.Now(),
	}
	if err := ctrl.fsmRepo.Create(&fsm); err != nil {
		// FSM 创建失败，回滚猫的创建
		ctrl.repo.Delete(cat.CatID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CatFSM: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"cat": cat,
		"fsm": fsm,
	})
}

// UpdateCat 更新 Cat
func (ctrl *CatController) UpdateCat(c *gin.Context) {
	idStr := c.Param("cat_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	cat, err := ctrl.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}
	var req UpdateCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := model.Cat{
		CatName:           req.CatName,
		CatPhotoUri:       req.CatPhotoUri,
		CatType:           req.CatType,
		CatGender:         req.CatGender,
		MasterName:        req.MasterName,
		MasterPhoneNumber: req.MasterPhoneNumber,
	}
	if err := ctrl.repo.Update(cat, &updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// DeleteCat 删除 Cat
func (ctrl *CatController) DeleteCat(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "Cat deleted successfully"})
}

// HealthCheck 健康检查
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
