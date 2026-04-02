package router

import (
	"strings"

	"operation_admin/backend/internal/config"
	"operation_admin/backend/internal/http/handler"
	"operation_admin/backend/internal/http/middleware"
	"operation_admin/backend/internal/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// New 负责注册 Gin 路由、中间件与基础健康检查接口。
func New(
	cfg config.ServerConfig,
	logger *zap.Logger,
	healthHandler *handler.HealthHandler,
	authMiddleware *middleware.AdminAuthMiddleware,
	authHandler *handler.AuthHandler,
	adminUserHandler *handler.AdminUserHandler,
	followerUserHandler *handler.FollowerUserHandler,
) *gin.Engine {
	gin.SetMode(normalizeGinMode(cfg.Mode))

	// engine 是当前应用的根路由引擎。
	engine := gin.New()
	engine.GET("/healthz", healthHandler.Check)

	// publicAdminGroup 负责承载无需登录即可访问的后台接口。
	publicAdminGroup := engine.Group("/admin/v1")
	publicAdminGroup.GET("/healthz", healthHandler.Check)

	if authHandler != nil {
		authHandler.RegisterPublicRoutes(publicAdminGroup)
	}

	// protectedAdminGroup 负责承载必须通过 Access Token 鉴权的后台接口。
	protectedAdminGroup := engine.Group("/admin/v1")
	if authMiddleware != nil {
		protectedAdminGroup.Use(authMiddleware.RequireAccessToken())
	}

	if authHandler != nil {
		authHandler.RegisterProtectedRoutes(protectedAdminGroup)
	}

	if adminUserHandler != nil {
		adminUserHandler.RegisterRoutes(protectedAdminGroup)
	}

	if followerUserHandler != nil {
		followerUserHandler.RegisterRoutes(protectedAdminGroup)
	}

	engine.NoRoute(func(ctx *gin.Context) {
		response.NotFound(ctx, "")
	})

	engine.NoMethod(func(ctx *gin.Context) {
		response.MethodNotAllowed(ctx, "")
	})

	return engine
}

// normalizeGinMode 负责把外部配置映射为 Gin 支持的运行模式。
func normalizeGinMode(mode string) string {
	// normalizedMode 用于统一处理大小写和空白字符差异。
	normalizedMode := strings.TrimSpace(strings.ToLower(mode))

	switch normalizedMode {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		return normalizedMode
	default:
		return gin.ReleaseMode
	}
}
