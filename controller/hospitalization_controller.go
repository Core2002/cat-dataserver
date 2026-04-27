package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"fifu.fun/cat-dataserver/middleware"
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
	catRepo          *repository.CatRepository
	siteRepo         *repository.SiteRepository
	catFSMRepo       *repository.CatFSMRepository
	catEventRepo     *repository.CatEventRepository
	actionProcessor  *middleware.ActionProcessor
}

// NewHospitalizationController 创建出入院管理控制器
func NewHospitalizationController(
	catRepo *repository.CatRepository,
	siteRepo *repository.SiteRepository,
	catFSMRepo *repository.CatFSMRepository,
	catEventRepo *repository.CatEventRepository,
	actionProcessor *middleware.ActionProcessor,
) *HospitalizationController {
	return &HospitalizationController{
		catRepo:         catRepo,
		siteRepo:        siteRepo,
		catFSMRepo:      catFSMRepo,
		catEventRepo:    catEventRepo,
		actionProcessor: actionProcessor,
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

	fsm, err := ctrl.catFSMRepo.FindByID(req.CatID)
	if err == nil && fsm.SiteID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cat is already admitted in a site"})
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detail, _ := json.Marshal(model.AdmissionActionDetail{
		Reason:       req.AdmissionReason,
		Notes:        req.AdmissionNote,
		TemperatureC: &req.InitialTemperatureC,
		WeightKG:     &req.InitialWeightKG,
	})

	action := model.CatAction{
		CatID:        req.CatID,
		SiteID:       req.SiteID,
		UserID:       req.UserID,
		ActionType:   model.CatActionAdmit,
		ActionDetail: string(detail),
	}
	_, err = ctrl.actionProcessor.ProcessAction(&action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event := model.CatEvent{
		EventType: model.CatAdmitted,
		SiteID:    req.SiteID,
		UserID:    req.UserID,
		CatID:     req.CatID,
		Detail:    req.AdmissionReason + "；" + req.AdmissionNote,
	}
	if err := ctrl.catEventRepo.Create(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "入院动作成功，但写入事件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"action": action, "event": event})
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

	fsm, err := ctrl.catFSMRepo.FindByID(req.CatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cat fsm not found"})
		return
	}
	if fsm.SiteID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cat is not admitted"})
		return
	}
	currentSiteID := fsm.SiteID

	detail, _ := json.Marshal(model.DischargeActionDetail{
		Reason:       req.DischargeReason,
		Notes:        req.DischargeNote,
		TemperatureC: req.FinalTemperatureC,
		WeightKG:     req.FinalWeightKG,
	})

	action := model.CatAction{
		CatID:        req.CatID,
		SiteID:       currentSiteID,
		UserID:       req.UserID,
		ActionType:   model.CatActionDischarge,
		ActionDetail: string(detail),
	}
	_, err = ctrl.actionProcessor.ProcessAction(&action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event := model.CatEvent{
		EventType: model.CatDischarged,
		SiteID:    currentSiteID,
		UserID:    req.UserID,
		CatID:     req.CatID,
		Detail:    req.DischargeReason + "；" + req.DischargeNote,
	}
	if err := ctrl.catEventRepo.Create(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "出院动作成功，但写入事件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"action": action, "event": event})
}
