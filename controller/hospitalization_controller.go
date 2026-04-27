package controller

import (
	"errors"
	"net/http"
	"time"

	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/model"
	"fifu.fun/cat-dataserver/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdmitCatRequest 办理入院请求
type AdmitCatRequest struct {
	CatID               uint    `json:"cat_id" binding:"required,min=1"`
	SiteID              uint    `json:"site_id" binding:"required,min=1"`
	UserID              uint    `json:"user_id" binding:"required,min=1"`
	AdmissionReason     string  `json:"admission_reason" binding:"required,min=1,max=200"`
	AdmissionNote       string  `json:"admission_note" binding:"omitempty,max=1000"`
	InitialTemperatureC float32 `json:"initial_temperature_c" binding:"required,gte=0,lte=50"`
	InitialWeightKG     float32 `json:"initial_weight_kg" binding:"required,gte=0.1,lte=25"`
}

// DischargeCatRequest 办理出院请求
type DischargeCatRequest struct {
	CatID             uint     `json:"cat_id" binding:"required,min=1"`
	UserID            uint     `json:"user_id" binding:"required,min=1"`
	DischargeReason   string   `json:"discharge_reason" binding:"required,min=1,max=200"`
	DischargeNote     string   `json:"discharge_note" binding:"omitempty,max=1000"`
	FinalTemperatureC *float32 `json:"final_temperature_c" binding:"omitempty,gte=0,lte=50"`
	FinalWeightKG     *float32 `json:"final_weight_kg" binding:"omitempty,gte=0.1,lte=25"`
}

// HospitalizationController 出入院管理控制器
type HospitalizationController struct {
	catRepo    *repository.CatRepository
	siteRepo   *repository.SiteRepository
}

// NewHospitalizationController 创建出入院管理控制器
func NewHospitalizationController(
	catRepo *repository.CatRepository,
	siteRepo *repository.SiteRepository,
) *HospitalizationController {
	return &HospitalizationController{
		catRepo:    catRepo,
		siteRepo:   siteRepo,
	}
}

// AdmitCat 办理入院
func (ctrl *HospitalizationController) AdmitCat(c *gin.Context) {
	var req AdmitCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 基础存在性校验
	if _, err := ctrl.catRepo.FindByID(req.CatID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CatID does not exist"})
		return
	}
	if _, err := ctrl.siteRepo.FindByID(req.SiteID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SiteID does not exist"})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		// 完整初始化：确保 CatFSM 存在，并写入入院初始状态。
		// 注意：这里不保存“入院表单记录”，仅执行业务状态变更。
		var fsm model.CatFSM
		err := tx.Where("cat_id = ?", req.CatID).First(&fsm).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fsm = model.CatFSM{
				CatID:         req.CatID,
				SiteID:        req.SiteID,
				TemperatureC:  req.InitialTemperatureC,
				WeightKG:      req.InitialWeightKG,
				TrimNailsTime: now,
			}
			return tx.Create(&fsm).Error
		}
		if err != nil {
			return err
		}
		if fsm.SiteID > 0 {
			return errors.New("cat is already admitted in a site")
		}

		fsm.SiteID = req.SiteID
		fsm.TemperatureC = req.InitialTemperatureC
		fsm.WeightKG = req.InitialWeightKG
		return tx.Save(&fsm).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":            "admitted",
		"cat_id":             req.CatID,
		"site_id":            req.SiteID,
		"initial_reason":     req.AdmissionReason,
		"initial_note":       req.AdmissionNote,
		"initial_operatorId": req.UserID,
	})
}

// DischargeCat 办理出院
func (ctrl *HospitalizationController) DischargeCat(c *gin.Context) {
	var req DischargeCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := ctrl.catRepo.FindByID(req.CatID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CatID does not exist"})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 出院结算：更新 FSM 并解绑设施（site_id=0）。
		// 注意：这里不保存“出院表单记录”，仅执行业务状态变更。
		var fsm model.CatFSM
		if err := tx.Where("cat_id = ?", req.CatID).First(&fsm).Error; err != nil {
			return err
		}
		if fsm.SiteID == 0 {
			return errors.New("cat is not admitted")
		}
		fsm.SiteID = 0
		if req.FinalTemperatureC != nil {
			fsm.TemperatureC = *req.FinalTemperatureC
		}
		if req.FinalWeightKG != nil {
			fsm.WeightKG = *req.FinalWeightKG
		}
		return tx.Save(&fsm).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "discharged",
		"cat_id":               req.CatID,
		"discharge_reason":     req.DischargeReason,
		"discharge_note":       req.DischargeNote,
		"discharge_operatorId": req.UserID,
	})
}
