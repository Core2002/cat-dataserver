package router

import (
	"fifu.fun/cat-dataserver/controller"
	"fifu.fun/cat-dataserver/repository"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 初始化 repository 和 controller
	catRepo := repository.NewCatRepository()
	catController := controller.NewCatController(catRepo)

	// Cat CRUD 路由
	r.GET("/cats", catController.GetCats)
	r.GET("/cats/:id", catController.GetCat)
	r.POST("/cats", catController.CreateCat)
	r.PUT("/cats/:id", catController.UpdateCat)
	r.DELETE("/cats/:id", catController.DeleteCat)

	// 健康检查
	r.GET("/health", controller.HealthCheck)

	return r
}
