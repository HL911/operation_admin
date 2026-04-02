package handler

import (
	"context"
	"errors"
	"time"

	"operation_admin/backend/internal/adminauth"
	"operation_admin/backend/internal/http/middleware"
	"operation_admin/backend/internal/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthService 描述后台鉴权 Handler 所需的业务能力。
type AuthService interface {
	// Login 负责校验后台账号凭证并返回新的登录会话。
	Login(ctx context.Context, command adminauth.LoginCommand) (*adminauth.LoginResult, error)
	// Refresh 负责基于有效 Refresh Token 轮换登录会话。
	Refresh(ctx context.Context, command adminauth.RefreshCommand) (*adminauth.RefreshResult, error)
	// Logout 负责撤销当前后台账号名下的指定 Refresh Token。
	Logout(ctx context.Context, command adminauth.LogoutCommand) (*adminauth.LogoutResult, error)
}

// AuthHandler 描述后台鉴权相关 HTTP 处理器。
type AuthHandler struct {
	// Service 提供登录、刷新、登出等后台鉴权业务能力。
	Service AuthService
	// Logger 提供处理异常场景时的结构化日志能力。
	Logger *zap.Logger
}

// loginRequest 描述登录接口使用的请求体结构。
type loginRequest struct {
	// LoginName 表示当前登录请求提交的后台账号登录名。
	LoginName string `json:"loginName"`
	// Password 表示当前登录请求提交的明文密码。
	Password string `json:"password"`
}

// refreshRequest 描述刷新令牌接口使用的请求体结构。
type refreshRequest struct {
	// RefreshToken 表示当前刷新请求提交的原始 Refresh Token。
	RefreshToken string `json:"refreshToken"`
}

// logoutRequest 描述登出接口使用的请求体结构。
type logoutRequest struct {
	// RefreshToken 表示当前登出请求需要撤销的原始 Refresh Token。
	RefreshToken string `json:"refreshToken"`
}

// NewAuthHandler 负责创建后台鉴权 HTTP 处理器。
func NewAuthHandler(service AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		Service: service,
		Logger:  logger,
	}
}

// RegisterPublicRoutes 负责注册后台鉴权模块对外公开的路由。
func (h *AuthHandler) RegisterPublicRoutes(adminGroup *gin.RouterGroup) {
	adminGroup.POST("/auth/login", h.Login)
	adminGroup.POST("/auth/refresh", h.Refresh)
}

// RegisterProtectedRoutes 负责注册需要 Access Token 鉴权的后台鉴权路由。
func (h *AuthHandler) RegisterProtectedRoutes(adminGroup *gin.RouterGroup) {
	adminGroup.POST("/auth/logout", h.Logout)
	adminGroup.GET("/auth/me", h.Me)
}

// Login 负责处理后台开发者登录请求。
func (h *AuthHandler) Login(ctx *gin.Context) {
	// request 保存当前请求体绑定得到的数据。
	var request loginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ValidationFailed(ctx, "请求体格式错误")
		return
	}

	// command 保存转换后的服务层登录命令。
	command := adminauth.LoginCommand{
		LoginName: request.LoginName,
		Password:  request.Password,
		ClientIP:  ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
	}

	// loginResult 保存登录成功后需要返回给前端的结果对象。
	loginResult, err := h.Service.Login(ctx.Request.Context(), command)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "登录成功", loginResult)
}

// Refresh 负责处理后台登录会话刷新请求。
func (h *AuthHandler) Refresh(ctx *gin.Context) {
	// request 保存当前请求体绑定得到的数据。
	var request refreshRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ValidationFailed(ctx, "请求体格式错误")
		return
	}

	// command 保存转换后的服务层刷新命令。
	command := adminauth.RefreshCommand{
		RefreshToken: request.RefreshToken,
		ClientIP:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
	}

	// refreshResult 保存刷新成功后需要返回给前端的结果对象。
	refreshResult, err := h.Service.Refresh(ctx.Request.Context(), command)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "刷新成功", refreshResult)
}

// Logout 负责处理后台登录会话登出请求。
func (h *AuthHandler) Logout(ctx *gin.Context) {
	// currentAdmin 保存当前已通过中间件鉴权的后台账号信息。
	currentAdmin, ok := middleware.CurrentAdmin(ctx)
	if !ok {
		response.Unauthorized(ctx, "")
		return
	}

	// request 保存当前请求体绑定得到的数据。
	var request logoutRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ValidationFailed(ctx, "请求体格式错误")
		return
	}

	// command 保存转换后的服务层登出命令。
	command := adminauth.LogoutCommand{
		AdminUserID:  currentAdmin.AdminUserID,
		RefreshToken: request.RefreshToken,
	}

	// logoutResult 保存登出成功后需要返回给前端的结果对象。
	logoutResult, err := h.Service.Logout(ctx.Request.Context(), command)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", logoutResult)
}

// Me 负责返回当前已登录后台开发者的信息。
func (h *AuthHandler) Me(ctx *gin.Context) {
	// currentAdmin 保存当前已通过中间件鉴权的后台账号信息。
	currentAdmin, ok := middleware.CurrentAdmin(ctx)
	if !ok {
		response.Unauthorized(ctx, "")
		return
	}

	// profileView 保存当前登录后台账号对应的接口返回对象。
	profileView := adminauth.AdminProfileView{
		AdminUserID: currentAdmin.AdminUserID,
		LoginName:   currentAdmin.LoginName,
		DisplayName: currentAdmin.DisplayName,
		RoleCode:    currentAdmin.RoleCode,
		Status:      currentAdmin.Status,
		LastLoginAt: formatCurrentAdminLastLoginAt(currentAdmin.LastLoginAt),
	}

	response.Success(ctx, "请求成功", profileView)
}

// handleServiceError 负责把后台鉴权服务层错误统一转换为 HTTP 响应。
func (h *AuthHandler) handleServiceError(ctx *gin.Context, err error) {
	// validationError 保存从错误链中提取出的参数校验错误。
	var validationError *adminauth.ValidationError

	switch {
	case errors.As(err, &validationError):
		response.ValidationFailed(ctx, validationError.Message)
	case errors.Is(err, adminauth.ErrInvalidCredentials):
		response.Unauthorized(ctx, "登录名或密码错误")
	case errors.Is(err, adminauth.ErrRefreshTokenInvalid):
		response.Unauthorized(ctx, "")
	case errors.Is(err, adminauth.ErrAccessTokenInvalid):
		response.Unauthorized(ctx, "")
	case errors.Is(err, adminauth.ErrAdminUserDisabled):
		response.Forbidden(ctx, "账号已被禁用")
	case errors.Is(err, adminauth.ErrAdminPermissionDenied):
		response.Forbidden(ctx, "当前账号无权限执行该操作")
	default:
		h.Logger.Error("处理后台鉴权请求失败", zap.Error(err))
		response.InternalError(ctx, "")
	}
}

// formatCurrentAdminLastLoginAt 负责把当前后台账号最近登录时间格式化为统一的时间字符串。
func formatCurrentAdminLastLoginAt(lastLoginAtValue *time.Time) *string {
	if lastLoginAtValue == nil {
		return nil
	}

	// formattedValue 保存格式化后的最近登录时间字符串。
	formattedValue := lastLoginAtValue.Format(response.DateTimeLayout)
	return &formattedValue
}
