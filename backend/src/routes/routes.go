package routes

import (
	"time"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/controllers"
	"groundTurn/src/middlewares"
	"groundTurn/src/repositories"
	"groundTurn/src/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Start 装配依赖并启动 HTTP 服务。
// 依赖按 repositories -> services -> controllers 手工注入（分层拆分，禁止单文件聚合业务）。
func Start(addr string, cfg config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 仓储
	turnaroundRepo := repositories.NewFlightTurnaroundRepository(db)
	taskRepo := repositories.NewGroundTaskRepository(db)
	delayRepo := repositories.NewDelayEventRepository(db)
	resourceRepo := repositories.NewGroundResourceRepository(db)
	bookingRepo := repositories.NewResourceBookingRepository(db)
	userRepo := repositories.NewUserRepository(db)
	auditRepo := repositories.NewAuditLogRepository(db)

	// 服务
	viewService := services.NewViewService(turnaroundRepo, taskRepo, delayRepo)
	taskService := services.NewGroundTaskService(taskRepo, turnaroundRepo, delayRepo)
	delayService := services.NewDelayEventService(delayRepo, turnaroundRepo)
	turnaroundService := services.NewFlightTurnaroundService(turnaroundRepo, viewService)
	resourceService := services.NewGroundResourceService(resourceRepo)
	bookingService := services.NewResourceBookingService(bookingRepo, resourceRepo)
	authService := services.NewAuthService(userRepo, auditRepo, cfg)
	dashboardService := services.NewDashboardService(viewService)

	// 控制器
	taskCtl := controllers.NewGroundTaskController(taskService)
	viewCtl := controllers.NewViewController(viewService)
	delayCtl := controllers.NewDelayEventController(delayService)
	turnaroundCtl := controllers.NewFlightTurnaroundController(turnaroundService)
	resourceCtl := controllers.NewGroundResourceController(resourceService)
	bookingCtl := controllers.NewResourceBookingController(bookingService)
	authCtl := controllers.NewAuthController(authService)
	dashboardCtl := controllers.NewDashboardController(dashboardService)

	r.Use(middlewares.ErrorHandlerMiddleware())
	r.Use(middlewares.RateLimitMiddleware(300, time.Minute))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "ground-turn"})
	})

	api := r.Group("/api")
	authCtl.RegisterPublicRoutes(api) // POST /api/auth/login

	// 读接口：任意已登录角色
	read := writeGroup(api, cfg, auditRepo) // 仅认证，不限制角色
	viewCtl.RegisterRoutes(read)
	delayCtl.RegisterReadRoutes(read)
	resourceCtl.RegisterReadRoutes(read)
	bookingCtl.RegisterReadRoutes(read)
	dashboardCtl.RegisterRoutes(read)
	authCtl.RegisterRoutes(read)

	// 派工 / 航班 / 延误登记关闭：地勤调度、运行督导
	dispatch := writeGroup(api, cfg, auditRepo, constants.RoleDispatcher, constants.RoleSupervisor)
	taskCtl.RegisterDispatchRoutes(dispatch)
	turnaroundCtl.RegisterRoutes(dispatch)
	delayCtl.RegisterWriteRoutes(dispatch)

	// 签收 / 完成 / 阻塞：班组、地勤调度、运行督导
	team := writeGroup(api, cfg, auditRepo, constants.RoleTeam, constants.RoleDispatcher, constants.RoleSupervisor)
	taskCtl.RegisterActionRoutes(team)

	// 资源台账 / 预约：资源管理员、运行督导
	rm := writeGroup(api, cfg, auditRepo, constants.RoleResourceManager, constants.RoleSupervisor)
	resourceCtl.RegisterWriteRoutes(rm)
	bookingCtl.RegisterWriteRoutes(rm)

	return r
}

// writeGroup 组装 认证 + RBAC(可选) + 审计 中间件；不传角色表示仅要求登录。
func writeGroup(api *gin.RouterGroup, cfg config.Config, auditRepo *repositories.AuditLogRepository, roles ...string) *gin.RouterGroup {
	g := api.Group("")
	g.Use(middlewares.AuthMiddleware(cfg.JWTSecret))
	if len(roles) > 0 {
		g.Use(middlewares.RBACMiddleware(roles...))
	}
	g.Use(middlewares.AuditLogMiddleware(auditRepo))
	return g
}
