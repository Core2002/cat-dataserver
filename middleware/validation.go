package middleware

import (
	"fifu.fun/cat-dataserver/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidatorMiddleware 数据校验中间件
func ValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 中间件会在 ShouldBindJSON/Query 时自动触发验证
		c.Next()
	}
}

// ValidationError 返回自定义的验证错误响应
func ValidationError(c *gin.Context, err error) {
	// 检查是否为验证错误
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			errorMsg := getErrorMessage(field, tag, param)
			errors[field] = errorMsg
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "数据验证失败",
			"errors": errors,
		})
		return
	}

	// 其他类型的错误
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

// getErrorMessage 根据验证标签返回友好的错误信息
func getErrorMessage(field, tag, param string) string {
	fieldName := getFieldName(field)

	switch tag {
	case "required":
		return fieldName + "不能为空"
	case "min":
		if param == "1" {
			return fieldName + "最小值为" + param
		}
		return fieldName + "长度不能小于" + param
	case "max":
		return fieldName + "长度不能超过" + param
	case "email":
		return fieldName + "格式不正确"
	case "oneof":
		return fieldName + "必须是以下值之一: " + param
	case "numeric":
		return fieldName + "必须是数字"
	case "gte":
		return fieldName + "必须大于等于" + param
	case "lte":
		return fieldName + "必须小于等于" + param
	case "gt":
		return fieldName + "必须大于" + param
	case "lt":
		return fieldName + "必须小于" + param
	default:
		return fieldName + "格式不正确"
	}
}

// getFieldName 将字段名转换为中文友好名称
func getFieldName(field string) string {
	// 从 model 包中定义的字段映射
	fieldMap := map[string]string{
		"CatID":                "猫ID",
		"CatName":              "猫名",
		"CatPhotoUri":          "猫照片",
		"CatType":              "猫种类",
		"CatGender":            "猫性别",
		"MasterName":           "主人姓名",
		"MasterPhoneNumber":    "主人电话",
		"SiteID":               "站点ID",
		"SiteName":             "站点名称",
		"SiteAddress":          "站点地址",
		"SiteAdminPhoneNumber": "站点管理员电话",
		"LastDisinfectTime":    "上次消毒时间",
		"LastFeedTime":         "上次喂食时间",
		"LastGiveWaterTime":    "上次喂水时间",
		"LastPlayTime":         "上次逗猫时间",
		"EventID":              "事件ID",
		"EventType":            "事件类型",
		"Detail":               "事件详情",
		"ActionID":             "操作ID",
		"UserID":               "用户ID",
		"ActionType":           "操作类型",
		"ActionDetail":         "操作详情",
		"TemperatureC":         "体温",
		"WeightKG":             "体重",
		"TrimNailsTime":        "修剪指甲时间",
		"Page":                 "页码",
		"PageSize":             "每页数量",
	}

	if name, ok := fieldMap[field]; ok {
		return name
	}

	// 转换驼峰命名
	return strings.ReplaceAll(field, "_", "")
}

// RegisterCustomValidators 注册自定义验证器
func RegisterCustomValidators(v *validator.Validate) {
	// 验证 CatEventType 是否合法
	v.RegisterValidation("catEventType", validateCatEventType)
	// 验证 CatActionType 是否合法
	v.RegisterValidation("catActionType", validateCatActionType)
	// 验证电话号码格式
	v.RegisterValidation("phone", validatePhoneNumber)
}

// validateCatEventType 验证猫事件类型
func validateCatEventType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	validTypes := []string{
		string(model.CatSick),
		string(model.CatInjure),
		string(model.CatPreg),
		string(model.CatBirth),
		string(model.CatDeath),
		string(model.CatContractTerminatio),
	}

	for _, t := range validTypes {
		if value == t {
			return true
		}
	}
	return false
}

// validateCatActionType 验证猫操作类型
func validateCatActionType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	validTypes := []string{
		string(model.CatActionFeed),
		string(model.CatActionGiveWater),
		string(model.CatActionTakeTemperature),
		string(model.CatActionPlay),
		string(model.CatActionSterilize),
		string(model.CatActionHealthCheck),
		string(model.CatActionDeworm),
		string(model.CatActionCleanLitter),
		string(model.CatActionDisinfect),
		string(model.CatActionTrimNails),
		string(model.CatActionWashFeet),
		string(model.CatActionVaccinate),
	}

	for _, t := range validTypes {
		if value == t {
			return true
		}
	}
	return false
}

// validatePhoneNumber 验证电话号码格式（中国大陆手机号）
func validatePhoneNumber(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) != 11 {
		return false
	}
	// 简单验证：以1开头，后面10位为数字
	if value[0] != '1' {
		return false
	}
	for _, c := range value {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
