package handler

import (
	"context"
	"errors"
	"net/http"

	"operation_admin/backend/internal/adminauth"
	"operation_admin/backend/internal/http/middleware"
	"operation_admin/backend/internal/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AdminUserService 描述管理员创建后台账号 Handler 所需的业务能力。
type AdminUserService interface {
	// CreateAdminUser 负责由管理员创建新的后台开发者账号。
	CreateAdminUser(ctx context.Context, operator adminauth.CurrentAdmin, command adminauth.CreateAdminUserCommand) (*adminauth.CreatedAdminUserView, error)
}

// AdminUserHandler 描述后台开发者账号管理 HTTP 处理器。
type AdminUserHandler struct {
	// Service 提供管理员创建后台开发者账号所需的业务能力。
	Service AdminUserService
	// Logger 提供处理异常场景时的结构化日志能力。
	Logger *zap.Logger
}

// createAdminUserRequest 描述管理员创建后台账号接口使用的请求体结构。
type createAdminUserRequest struct {
	// LoginName 表示新后台账号的登录名。
	LoginName string `json:"loginName"`
	// Password 表示新后台账号的明文密码。
	Password string `json:"password"`
	// DisplayName 表示新后台账号的显示名。
	DisplayName string `json:"displayName"`
}

// NewAdminUserHandler 负责创建后台开发者账号管理 HTTP 处理器。
func NewAdminUserHandler(service AdminUserService, logger *zap.Logger) *AdminUserHandler {
	return &AdminUserHandler{
		Service: service,
		Logger:  logger,
	}
}

// RegisterRoutes 负责把后台开发者账号管理路由挂载到受保护后台分组下。
func (h *AdminUserHandler) RegisterRoutes(adminGroup *gin.RouterGroup) {
	adminGroup.POST("/admin-users", h.Create)
}

// Create 负责处理管理员创建后台开发者账号请求。
func (h *AdminUserHandler) Create(ctx *gin.Context) {
	// currentAdmin 保存当前已通过中间件鉴权的后台账号信息。
	currentAdmin, ok := middleware.CurrentAdmin(ctx)
	if !ok {
		response.Unauthorized(ctx, "")
		return
	}

	// request 保存当前请求体绑定得到的数据。
	var request createAdminUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ValidationFailed(ctx, "请求体格式错误")
		return
	}

	// command 保存转换后的服务层创建命令。
	command := adminauth.CreateAdminUserCommand{
		LoginName:   request.LoginName,
		Password:    request.Password,
		DisplayName: request.DisplayName,
	}

	// createdAdminUserView 保存创建成功后需要返回给前端的结果对象。
	createdAdminUserView, err := h.Service.CreateAdminUser(ctx.Request.Context(), *currentAdmin, command)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", createdAdminUserView)
}

// handleServiceError 负责把后台开发者账号管理服务层错误统一转换为 HTTP 响应。
func (h *AdminUserHandler) handleServiceError(ctx *gin.Context, err error) {
	// validationError 保存从错误链中提取出的参数校验错误。
	var validationError *adminauth.ValidationError

	switch {
	case errors.As(err, &validationError):
		response.ValidationFailed(ctx, validationError.Message)
	case errors.Is(err, adminauth.ErrAdminUserDuplicated):
		response.Fail(ctx, http.StatusConflict, adminauth.ErrorCodeAdminUserDuplicated, "后台账号已存在")
	case errors.Is(err, adminauth.ErrAdminPermissionDenied):
		response.Forbidden(ctx, "当前账号无权限执行该操作")
	default:
		h.Logger.Error("处理后台开发者账号请求失败", zap.Error(err))
		response.InternalError(ctx, "")
	}
}
