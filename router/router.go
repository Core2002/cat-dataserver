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
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化 repository 和 controller
	catRepo := repository.NewCatRepository()
	catController := controller.NewCatController(catRepo)

	catEventRepo := repository.NewCatEventRepository()
	catEventController := controller.NewCatEventController(catEventRepo)

	catActionRepo := repository.NewCatActionRepository()
	catFSMRepo := repository.NewCatFSMRepository()

	// 初始化动作处理器
	actionProcessor := middleware.NewActionProcessor(catActionRepo, catFSMRepo)
	catActionController := controller.NewCatActionController(catActionRepo, actionProcessor)

	catFSMController := controller.NewCatFSMController(catFSMRepo)

	siteRepo := repository.NewSiteRepository()
	siteController := controller.NewSiteController(siteRepo)

	// Cat CRUD 路由
	r.GET("/cats", catController.GetCats)
	r.GET("/cats/page", catController.GetCatsPage)
	r.GET("/cats/:id", catController.GetCat)
	r.POST("/cats", catController.CreateCat)
	r.PUT("/cats/:id", catController.UpdateCat)
	r.DELETE("/cats/:id", catController.DeleteCat)

	// CatEvent CRUD 路由
	r.GET("/cat-events", catEventController.GetCatEvents)
	r.GET("/cat-events/page", catEventController.GetCatEventsPage)
	r.GET("/cat-events/:id", catEventController.GetCatEvent)
	r.GET("/cat-events/cat/:cat_id", catEventController.GetCatEventsByCatID)
	r.GET("/cat-events/site/:site_id", catEventController.GetCatEventsBySiteID)
	r.POST("/cat-events", catEventController.CreateCatEvent)
	r.PUT("/cat-events/:id", catEventController.UpdateCatEvent)
	r.DELETE("/cat-events/:id", catEventController.DeleteCatEvent)

	// CatAction CRUD 路由
	r.GET("/cat-actions", catActionController.GetCatActions)
	r.GET("/cat-actions/page", catActionController.GetCatActionsPage)
	r.GET("/cat-actions/:id", catActionController.GetCatAction)
	r.GET("/cat-actions/cat/:cat_id", catActionController.GetCatActionsByCatID)
	r.GET("/cat-actions/site/:site_id", catActionController.GetCatActionsBySiteID)
	r.GET("/cat-actions/user/:user_id", catActionController.GetCatActionsByUserID)
	r.POST("/cat-actions", catActionController.CreateCatAction)
	r.PUT("/cat-actions/:id", catActionController.UpdateCatAction)
	r.DELETE("/cat-actions/:id", catActionController.DeleteCatAction)

	// CatFSM CRUD 路由
	r.GET("/cat-fsms", catFSMController.GetCatFSMs)
	r.GET("/cat-fsms/page", catFSMController.GetCatFSMsPage)
	r.GET("/cat-fsms/:id", catFSMController.GetCatFSM)
	r.GET("/cat-fsms/site/:site_id", catFSMController.GetCatFSMsBySiteID)
	r.POST("/cat-fsms", catFSMController.CreateCatFSM)
	r.PUT("/cat-fsms/:id", catFSMController.UpdateCatFSM)
	r.DELETE("/cat-fsms/:id", catFSMController.DeleteCatFSM)
	r.PATCH("/cat-fsms/:cat_id/temperature", catFSMController.UpdateTemperature)
	r.PATCH("/cat-fsms/:cat_id/weight", catFSMController.UpdateWeight)
	r.PATCH("/cat-fsms/:cat_id/trim-nails-time", catFSMController.UpdateTrimNailsTime)

	// Site CRUD 路由
	r.GET("/sites", siteController.GetSites)
	r.GET("/sites/page", siteController.GetSitesPage)
	r.GET("/sites/:id", siteController.GetSite)
	r.POST("/sites", siteController.CreateSite)
	r.PUT("/sites/:id", siteController.UpdateSite)
	r.DELETE("/sites/:id", siteController.DeleteSite)
	r.PATCH("/sites/:site_id/disinfect-time", siteController.UpdateDisinfectTime)
	r.PATCH("/sites/:site_id/feed-time", siteController.UpdateFeedTime)
	r.PATCH("/sites/:site_id/give-water-time", siteController.UpdateGiveWaterTime)
	r.PATCH("/sites/:site_id/play-time", siteController.UpdatePlayTime)

	// 健康检查
	r.GET("/health", controller.HealthCheck)

	return r
}
