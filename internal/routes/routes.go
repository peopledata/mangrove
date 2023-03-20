package routes

import (
	"patronus/internal/controller"
	"patronus/internal/controller/admin"
	"patronus/internal/controller/api"
	"patronus/internal/logger"
	"patronus/internal/routes/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 跨域配置
	config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://google.com"}
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// Serve static files from the "static" directory
	r.Use(static.Serve("/", static.LocalFile("./ui/dist", true)))

	r.GET("/ping", controller.Ping)

	// 后台管理相关接口
	adminV1 := r.Group("/admin/api/v1")

	adminV1.POST("/signup", admin.SignUpHandler)
	adminV1.POST("/login", admin.LoginHandler)
	adminV1.POST("/refresh_token", admin.RefreshTokenHandler)
	adminV1.GET("/dashboard", middleware.JWTAuthMiddleware(), admin.DashboardHandler)
	adminV1.GET("/routes", middleware.JWTAuthMiddleware(), admin.RoutesHandler)
	adminV1.GET("/user", middleware.JWTAuthMiddleware(), admin.GetUserHandler)

	// 需求管理
	adminV1.GET("/demand", middleware.JWTAuthMiddleware(), admin.DemandListHandler)          // 需求列表
	adminV1.GET("/demand/:id/info", middleware.JWTAuthMiddleware(), admin.DemandInfoHandler) // 基本信息                                   // 需求详情
	adminV1.GET("/demand/:id/detail", middleware.JWTAuthMiddleware(), admin.DemandDetailHandler)
	adminV1.GET("/demand/:id/contract_record", middleware.JWTAuthMiddleware(), admin.ContractRecordsHandler)                    // 签约记录                  //
	adminV1.POST("/demand", middleware.JWTAuthMiddleware(), admin.DemandCreateHandler)                                          // 新建需求
	adminV1.POST("/demand/:id", middleware.JWTAuthMiddleware(), admin.DemandUpdateHandler)                                      // 更新需求
	adminV1.POST("/demand/:id/publish", middleware.JWTAuthMiddleware(), middleware.EthMiddleware(), admin.DemandPublishHandler) // 发布需求
	adminV1.POST("/demand/:id/delete", middleware.JWTAuthMiddleware(), admin.DemandDeleteHandler)                               // 删除需求

	// 任务管理
	adminV1.POST("/demand/:id/task", middleware.JWTAuthMiddleware(), admin.TaskCreateHandler) // 触发一次新的任务
	adminV1.GET("/demand/:id/task", middleware.JWTAuthMiddleware(), admin.TaskListHandler)    // 任务列表

	// 外部API接口
	apiV1 := r.Group("/api/v1")
	apiV1.GET("/demand", api.DemandListHandler)                            // 获取发布的需求列表
	apiV1.GET("/demand/contract/:category", api.DemandContractListHandler) // 获取某个分类下发布的需求合约列表
	apiV1.GET("/demand/:id", api.DemandDetailHandler)                      // 获取需求详细信息

	return r
}
