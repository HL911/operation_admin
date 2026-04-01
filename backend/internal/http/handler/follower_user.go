package handler

import (
	"errors"
	"net/http"

	"operation_admin/backend/internal/followeruser"
	"operation_admin/backend/internal/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FollowerUserHandler 描述小龙虾用户 HTTP 处理器。
type FollowerUserHandler struct {
	// Service 提供小龙虾用户 CRUD 所需的业务能力。
	Service followeruser.Service
	// Logger 提供处理异常场景时的结构化日志能力。
	Logger *zap.Logger
}

// followerUserURI 描述小龙虾用户详情、更新、删除接口使用的路径参数。
type followerUserURI struct {
	// UserID 表示路径中的目标用户 ID。
	UserID string `uri:"userId"`
}

// createFollowerUserRequest 描述创建接口使用的请求体结构。
type createFollowerUserRequest struct {
	// UserID 表示要创建的用户 ID。
	UserID string `json:"userId"`
	// AccountStatus 表示要写入的账户状态。
	AccountStatus string `json:"accountStatus"`
	// StrategyStatus 表示要写入的策略状态。
	StrategyStatus string `json:"strategyStatus"`
	// BindingStatus 表示要写入的绑定状态。
	BindingStatus string `json:"bindingStatus"`
	// ResponsibilityDomain 表示要写入的责任域。
	ResponsibilityDomain string `json:"responsibilityDomain"`
}

// listFollowerUsersRequest 描述列表接口使用的查询参数结构。
type listFollowerUsersRequest struct {
	// PageNum 表示请求页码。
	PageNum int `form:"pageNum"`
	// PageSize 表示请求分页大小。
	PageSize int `form:"pageSize"`
	// UserID 表示按 userId 精确筛选。
	UserID string `form:"userId"`
	// AccountStatus 表示按账户状态筛选。
	AccountStatus string `form:"accountStatus"`
	// StrategyStatus 表示按策略状态筛选。
	StrategyStatus string `form:"strategyStatus"`
	// BindingStatus 表示按绑定状态筛选。
	BindingStatus string `form:"bindingStatus"`
	// ResponsibilityDomain 表示按责任域筛选。
	ResponsibilityDomain string `form:"responsibilityDomain"`
	// UpdatedFrom 表示更新时间范围开始值。
	UpdatedFrom string `form:"updatedFrom"`
	// UpdatedTo 表示更新时间范围结束值。
	UpdatedTo string `form:"updatedTo"`
}

// updateFollowerUserRequest 描述更新接口使用的请求体结构。
type updateFollowerUserRequest struct {
	// AccountStatus 表示新的账户状态。
	AccountStatus string `json:"accountStatus"`
	// StrategyStatus 表示新的策略状态。
	StrategyStatus string `json:"strategyStatus"`
	// BindingStatus 表示新的绑定状态。
	BindingStatus string `json:"bindingStatus"`
	// ResponsibilityDomain 表示新的责任域。
	ResponsibilityDomain string `json:"responsibilityDomain"`
	// RowVersion 表示调用方提交的当前版本号。
	RowVersion int `json:"rowVersion"`
}

// deleteFollowerUserRequest 描述删除接口使用的请求体结构。
type deleteFollowerUserRequest struct {
	// RowVersion 表示调用方提交的当前版本号。
	RowVersion int `json:"rowVersion"`
}

// NewFollowerUserHandler 负责创建小龙虾用户 HTTP 处理器。
func NewFollowerUserHandler(service followeruser.Service, logger *zap.Logger) *FollowerUserHandler {
	return &FollowerUserHandler{
		Service: service,
		Logger:  logger,
	}
}

// RegisterRoutes 负责把小龙虾用户 CRUD 路由挂载到 admin 分组下。
func (h *FollowerUserHandler) RegisterRoutes(adminGroup *gin.RouterGroup) {
	adminGroup.POST("/follower-users", h.Create)
	adminGroup.GET("/follower-users", h.List)
	adminGroup.GET("/follower-users/:userId", h.Get)
	adminGroup.PATCH("/follower-users/:userId", h.Update)
	adminGroup.DELETE("/follower-users/:userId", h.Delete)
}

// Create 负责处理创建小龙虾用户请求。
func (h *FollowerUserHandler) Create(ctx *gin.Context) {
	// request 保存当前请求体绑定得到的数据。
	var request createFollowerUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ValidationFailed(ctx, "请求体格式错误")
		return
	}

	// command 保存转换后的服务层创建命令。
	command := followeruser.CreateCommand{
		UserID:               request.UserID,
		AccountStatus:        request.AccountStatus,
		StrategyStatus:       request.StrategyStatus,
		BindingStatus:        request.BindingStatus,
		ResponsibilityDomain: request.ResponsibilityDomain,
	}

	// followerUserView 保存创建成功后的返回对象。
	followerUserView, err := h.Service.Create(ctx.Request.Context(), command)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", followerUserView)
}

// List 负责处理分页查询小龙虾用户请求。
func (h *FollowerUserHandler) List(ctx *gin.Context) {
	// request 保存当前查询参数绑定得到的数据。
	var request listFollowerUsersRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		response.ValidationFailed(ctx, "查询参数格式错误")
		return
	}

	// query 保存转换后的服务层列表查询条件。
	query := followeruser.ListQuery{
		PageNum:              request.PageNum,
		PageSize:             request.PageSize,
		UserID:               request.UserID,
		AccountStatus:        request.AccountStatus,
		StrategyStatus:       request.StrategyStatus,
		BindingStatus:        request.BindingStatus,
		ResponsibilityDomain: request.ResponsibilityDomain,
		UpdatedFrom:          request.UpdatedFrom,
		UpdatedTo:            request.UpdatedTo,
	}

	// followerUserViews 保存当前页返回对象列表。
	followerUserViews, total, pageNum, pageSize, err := h.Service.List(ctx.Request.Context(), query)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.SuccessPage(ctx, "查询成功", followerUserViews, total, pageNum, pageSize)
}

// Get 负责处理小龙虾用户详情查询请求。
func (h *FollowerUserHandler) Get(ctx *gin.Context) {
	// uri 保存路径参数绑定得到的数据。
	var uri followerUserURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.ValidationFailed(ctx, "路径参数格式错误")
		return
	}

	// followerUserView 保存详情查询成功后的返回对象。
	followerUserView, err := h.Service.Get(ctx.Request.Context(), uri.UserID)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", followerUserView)
}

// Update 负责处理小龙虾用户更新请求。
func (h *FollowerUserHandler) Update(ctx *gin.Context) {
	// uri 保存路径参数绑定得到的数据。
	var uri followerUserURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.ValidationFailed(ctx, "路径参数格式错误")
		return
	}

	// request 保存请求体绑定得到的更新参数。
	var request updateFollowerUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ValidationFailed(ctx, "请求体格式错误")
		return
	}

	// command 保存转换后的服务层更新命令。
	command := followeruser.UpdateCommand{
		UserID:               uri.UserID,
		AccountStatus:        request.AccountStatus,
		StrategyStatus:       request.StrategyStatus,
		BindingStatus:        request.BindingStatus,
		ResponsibilityDomain: request.ResponsibilityDomain,
		RowVersion:           request.RowVersion,
	}

	// followerUserView 保存更新成功后的返回对象。
	followerUserView, err := h.Service.Update(ctx.Request.Context(), command)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", followerUserView)
}

// Delete 负责处理小龙虾用户逻辑删除请求。
func (h *FollowerUserHandler) Delete(ctx *gin.Context) {
	// uri 保存路径参数绑定得到的数据。
	var uri followerUserURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.ValidationFailed(ctx, "路径参数格式错误")
		return
	}

	// request 保存请求体绑定得到的删除参数。
	var request deleteFollowerUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.ValidationFailed(ctx, "请求体格式错误")
		return
	}

	// command 保存转换后的服务层删除命令。
	command := followeruser.DeleteCommand{
		UserID:     uri.UserID,
		RowVersion: request.RowVersion,
	}

	// deleteResult 保存删除成功后的返回结果。
	deleteResult, err := h.Service.Delete(ctx.Request.Context(), command)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "请求成功", deleteResult)
}

// handleServiceError 负责把服务层错误统一转换为 HTTP 响应。
func (h *FollowerUserHandler) handleServiceError(ctx *gin.Context, err error) {
	// validationError 保存从错误链中提取出的业务校验错误。
	var validationError *followeruser.ValidationError

	switch {
	case errors.As(err, &validationError):
		response.ValidationFailed(ctx, validationError.Message)
	case errors.Is(err, followeruser.ErrFollowerUserNotFound):
		response.NotFound(ctx, "小龙虾用户不存在")
	case errors.Is(err, followeruser.ErrFollowerUserDuplicated):
		response.Fail(ctx, http.StatusConflict, followeruser.ErrorCodeFollowerUserDuplicated, "小龙虾用户已存在")
	case errors.Is(err, followeruser.ErrFollowerUserVersionConflict):
		response.Fail(ctx, http.StatusConflict, followeruser.ErrorCodeFollowerUserVersionConflict, "数据版本已过期，请刷新后重试")
	default:
		h.Logger.Error("处理小龙虾用户请求失败", zap.Error(err))
		response.InternalError(ctx, "")
	}
}
