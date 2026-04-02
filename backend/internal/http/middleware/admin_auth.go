package middleware

import (
	"context"
	"errors"
	"strings"

	"operation_admin/backend/internal/adminauth"
	"operation_admin/backend/internal/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// currentAdminContextKey 表示 Gin 上下文中保存当前后台账号信息时使用的键名。
	currentAdminContextKey = "currentAdmin"
)

// AccessAuthenticator 描述鉴权中间件校验 Access Token 所需的最小业务能力。
type AccessAuthenticator interface {
	// AuthenticateAccessToken 负责校验 Access Token 并返回当前后台账号上下文。
	AuthenticateAccessToken(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error)
}

// AdminAuthMiddleware 描述后台接口 Access Token 鉴权中间件。
type AdminAuthMiddleware struct {
	// Authenticator 提供 Access Token 解析和后台账号校验能力。
	Authenticator AccessAuthenticator
	// Logger 提供处理中间件内部异常时的结构化日志能力。
	Logger *zap.Logger
}

// NewAdminAuthMiddleware 负责创建后台接口 Access Token 鉴权中间件。
func NewAdminAuthMiddleware(authenticator AccessAuthenticator, logger *zap.Logger) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		Authenticator: authenticator,
		Logger:        logger,
	}
}

// RequireAccessToken 负责校验受保护后台接口携带的 Bearer Access Token。
func (m *AdminAuthMiddleware) RequireAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// authorizationHeader 保存当前请求头中的 Authorization 原始内容。
		authorizationHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
		if authorizationHeader == "" {
			response.Unauthorized(ctx, "")
			return
		}

		// accessToken 保存从 Authorization 头中提取出的 Bearer Token 字符串。
		accessToken, ok := extractBearerToken(authorizationHeader)
		if !ok {
			response.Unauthorized(ctx, "")
			return
		}

		// currentAdmin 保存 Access Token 校验成功后得到的后台账号上下文。
		currentAdmin, err := m.Authenticator.AuthenticateAccessToken(ctx.Request.Context(), accessToken)
		if err != nil {
			switch {
			case errors.Is(err, adminauth.ErrAccessTokenInvalid):
				response.Unauthorized(ctx, "")
			case errors.Is(err, adminauth.ErrAdminUserDisabled):
				response.Forbidden(ctx, "账号已被禁用")
			default:
				m.Logger.Error("校验后台 access token 失败", zap.Error(err))
				response.InternalError(ctx, "")
			}

			return
		}

		ctx.Set(currentAdminContextKey, currentAdmin)
		ctx.Next()
	}
}

// CurrentAdmin 负责从 Gin 上下文中获取当前已鉴权的后台账号信息。
func CurrentAdmin(ctx *gin.Context) (*adminauth.CurrentAdmin, bool) {
	// rawValue 保存从 Gin 上下文中取出的原始对象。
	rawValue, exists := ctx.Get(currentAdminContextKey)
	if !exists {
		return nil, false
	}

	// currentAdmin 保存类型断言成功后的后台账号上下文对象。
	currentAdmin, ok := rawValue.(*adminauth.CurrentAdmin)
	if !ok || currentAdmin == nil {
		return nil, false
	}

	return currentAdmin, true
}

// extractBearerToken 负责从 Authorization 请求头中提取 Bearer Token。
func extractBearerToken(authorizationHeader string) (string, bool) {
	// parts 保存按空白拆分后的 Authorization 请求头片段。
	parts := strings.Fields(authorizationHeader)
	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	// accessToken 保存提取出的 Bearer Token 字符串。
	accessToken := strings.TrimSpace(parts[1])
	if accessToken == "" {
		return "", false
	}

	return accessToken, true
}
