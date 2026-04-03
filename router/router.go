package router

import (
	"fifu.fun/cat-dataserver/controller"
	"fifu.fun/cat-dataserver/middleware"
	"fifu.fun/cat-dataserver/repository"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		middleware.RegisterCustomValidators(v)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5000",
			"http://localhost:5173",
			"http://localhost:9000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化 repository 和 controller
	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	catController := controller.NewCatController(catRepo)

	catEventRepo := repository.NewCatEventRepository()
	catEventController := controller.NewCatEventController(catEventRepo, catRepo, siteRepo)

	catActionRepo := repository.NewCatActionRepository()
	catFSMRepo := repository.NewCatFSMRepository()

	// 初始化动作处理器
	actionProcessor := middleware.NewActionProcessor(catActionRepo, catFSMRepo)
	catActionController := controller.NewCatActionController(catActionRepo, catRepo, siteRepo, actionProcessor)

	catFSMController := controller.NewCatFSMController(catFSMRepo)

	siteFSMRepo := repository.NewSiteFSMRepository()
	siteFSMController := controller.NewSiteFSMController(siteFSMRepo, siteRepo)

	siteController := controller.NewSiteController(siteRepo, siteFSMRepo)

	// Cat CRUD 路由
	r.GET("/cats/page", catController.GetCatsPage)
	r.GET("/cats/:cat_id", catController.GetCat)
	r.POST("/cats", catController.CreateCat)
	r.PUT("/cats/:cat_id", catController.UpdateCat)
	r.DELETE("/cats/:cat_id", catController.DeleteCat)

	// CatEvent CRUD 路由
	r.GET("/cat-events/page", catEventController.GetCatEventsPage)
	r.GET("/cat-events/:cat_id", catEventController.GetCatEvent)
	r.GET("/cat-events/cat/:cat_id", catEventController.GetCatEventsByCatID)
	r.GET("/cat-events/site/:site_id", catEventController.GetCatEventsBySiteID)
	r.POST("/cat-events", catEventController.CreateCatEvent)
	r.PUT("/cat-events/:cat_id", catEventController.UpdateCatEvent)
	r.DELETE("/cat-events/:cat_id", catEventController.DeleteCatEvent)

	// CatAction CRUD 路由
	r.GET("/cat-actions/page", catActionController.GetCatActionsPage)
	r.GET("/cat-actions/:cat_id", catActionController.GetCatAction)
	r.GET("/cat-actions/cat/:cat_id", catActionController.GetCatActionsByCatID)
	r.GET("/cat-actions/site/:site_id", catActionController.GetCatActionsBySiteID)
	r.GET("/cat-actions/user/:user_id", catActionController.GetCatActionsByUserID)
	r.POST("/cat-actions", catActionController.CreateCatAction)
	r.PUT("/cat-actions/:cat_id", catActionController.UpdateCatAction)
	r.DELETE("/cat-actions/:cat_id", catActionController.DeleteCatAction)

	// CatFSM CRUD 路由
	r.GET("/cat-fsms/page", catFSMController.GetCatFSMsPage)
	r.GET("/cat-fsms/:cat_id", catFSMController.GetCatFSM)
	r.GET("/cat-fsms/site/:site_id", catFSMController.GetCatFSMsBySiteID)
	r.POST("/cat-fsms", catFSMController.CreateCatFSM)
	r.PUT("/cat-fsms/:cat_id", catFSMController.UpdateCatFSM)
	r.DELETE("/cat-fsms/:cat_id", catFSMController.DeleteCatFSM)
	r.PATCH("/cat-fsms/:cat_id/temperature", catFSMController.UpdateTemperature)
	r.PATCH("/cat-fsms/:cat_id/weight", catFSMController.UpdateWeight)
	r.PATCH("/cat-fsms/:cat_id/trim-nails-time", catFSMController.UpdateTrimNailsTime)

	// Site CRUD 路由
	r.GET("/sites/page", siteController.GetSitesPage)
	r.GET("/sites/:site_id", siteController.GetSite)
	r.POST("/sites", siteController.CreateSite)
	r.PUT("/sites/:site_id", siteController.UpdateSite)
	r.DELETE("/sites/:site_id", siteController.DeleteSite)

	// SiteFSM CRUD 路由
	r.GET("/site-fsms/page", siteFSMController.GetSiteFSMsPage)
	r.GET("/site-fsms/:site_id", siteFSMController.GetSiteFSM)
	r.GET("/site-fsms/site/:site_id", siteFSMController.GetSiteFSMBySiteID)
	r.POST("/site-fsms", siteFSMController.CreateSiteFSM)
	r.PUT("/site-fsms/:site_id", siteFSMController.UpdateSiteFSM)
	r.DELETE("/site-fsms/:site_id", siteFSMController.DeleteSiteFSM)
	r.PATCH("/site-fsms/:site_id/disinfect-time", siteFSMController.UpdateDisinfectTime)
	r.PATCH("/site-fsms/:site_id/feed-time", siteFSMController.UpdateFeedTime)
	r.PATCH("/site-fsms/:site_id/give-water-time", siteFSMController.UpdateGiveWaterTime)
	r.PATCH("/site-fsms/:site_id/play-time", siteFSMController.UpdatePlayTime)
	r.PATCH("/site-fsms/:site_id/clean-litter-time", siteFSMController.UpdateCleanLitterTime)

	// 健康检查
	r.GET("/health", controller.HealthCheck)

	return r
}
