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
			"http://localhost:5200",
			"https://tls.internal",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化 repository 和 controller
	catRepo := repository.NewCatRepository()
	siteRepo := repository.NewSiteRepository()
	catFSMRepo := repository.NewCatFSMRepository()
	catController := controller.NewCatController(catRepo, catFSMRepo, siteRepo)

	catEventRepo := repository.NewCatEventRepository()
	catEventController := controller.NewCatEventController(catEventRepo, catRepo, siteRepo)

	catActionRepo := repository.NewCatActionRepository()

	// 初始化动作处理器
	actionProcessor := middleware.NewActionProcessor(catActionRepo, catFSMRepo)
	catActionController := controller.NewCatActionController(catActionRepo, catRepo, siteRepo, actionProcessor)

	catFSMController := controller.NewCatFSMController(catFSMRepo)

	siteFSMRepo := repository.NewSiteFSMRepository()
	siteFSMController := controller.NewSiteFSMController(siteFSMRepo, siteRepo)

	siteActionRepo := repository.NewSiteActionRepository()
	// 初始化站点动作处理器
	siteActionProcessor := middleware.NewSiteActionProcessor(siteActionRepo, siteFSMRepo)
	siteActionController := controller.NewSiteActionController(siteActionRepo, siteRepo, siteActionProcessor)

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

	// CatAction 路由（只读 + 创建，操作驱动状态机）
	r.GET("/cat-actions/page", catActionController.GetCatActionsPage)
	r.GET("/cat-actions/:action_id", catActionController.GetCatAction)
	r.GET("/cat-actions/cat/:cat_id", catActionController.GetCatActionsByCatID)
	r.GET("/cat-actions/site/:site_id", catActionController.GetCatActionsBySiteID)
	r.GET("/cat-actions/user/:user_id", catActionController.GetCatActionsByUserID)
	r.POST("/cat-actions", catActionController.CreateCatAction)

	// CatFSM 路由（只读，由 Action 驱动更新）
	r.GET("/cat-fsms/page", catFSMController.GetCatFSMsPage)
	r.GET("/cat-fsms/:cat_id", catFSMController.GetCatFSM)
	r.GET("/cat-fsms/site/:site_id", catFSMController.GetCatFSMsBySiteID)

	// Site CRUD 路由
	r.GET("/sites/page", siteController.GetSitesPage)
	r.GET("/sites/:site_id", siteController.GetSite)
	r.POST("/sites", siteController.CreateSite)
	r.PUT("/sites/:site_id", siteController.UpdateSite)
	r.DELETE("/sites/:site_id", siteController.DeleteSite)

	// SiteFSM 路由（只读，由 Action 驱动更新）
	r.GET("/site-fsms/page", siteFSMController.GetSiteFSMsPage)
	r.GET("/site-fsms/:site_id", siteFSMController.GetSiteFSM)
	r.GET("/site-fsms/site/:site_id", siteFSMController.GetSiteFSMBySiteID)

	// SiteAction 路由（只读 + 创建，操作驱动状态机）
	r.GET("/site-actions/page", siteActionController.GetSiteActionsPage)
	r.GET("/site-actions/:action_id", siteActionController.GetSiteAction)
	r.GET("/site-actions/site/:site_id", siteActionController.GetSiteActionsBySiteID)
	r.GET("/site-actions/user/:user_id", siteActionController.GetSiteActionsByUserID)
	r.POST("/site-actions", siteActionController.CreateSiteAction)

	// 健康检查
	r.GET("/health", controller.HealthCheck)

	return r
}
