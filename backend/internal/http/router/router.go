package router

import (
	"strings"

	"operation_admin/backend/internal/config"
	"operation_admin/backend/internal/http/handler"
	"operation_admin/backend/internal/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// New 负责注册 Gin 路由、中间件与基础健康检查接口。
func New(
	cfg config.ServerConfig,
	logger *zap.Logger,
	healthHandler *handler.HealthHandler,
	followerUserHandler *handler.FollowerUserHandler,
) *gin.Engine {
	gin.SetMode(normalizeGinMode(cfg.Mode))

	// engine 是当前应用的根路由引擎。
	engine := gin.New()
	engine.GET("/healthz", healthHandler.Check)

	// adminGroup 预留后台管理接口的统一分组前缀。
	adminGroup := engine.Group("/admin/v1")
	adminGroup.GET("/healthz", healthHandler.Check)

	if followerUserHandler != nil {
		followerUserHandler.RegisterRoutes(adminGroup)
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
