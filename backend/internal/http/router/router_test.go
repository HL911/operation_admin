package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"operation_admin/backend/internal/adminauth"
	"operation_admin/backend/internal/config"
	"operation_admin/backend/internal/followeruser"
	"operation_admin/backend/internal/http/handler"
	"operation_admin/backend/internal/http/middleware"

	"go.uber.org/zap"
)

// stubRouterFollowerUserService 描述用于路由测试的小龙虾用户业务桩实现。
type stubRouterFollowerUserService struct {
	// CreateFunc 描述当前测试场景下的创建桩函数。
	CreateFunc func(ctx context.Context, command followeruser.CreateCommand) (*followeruser.FollowerUserView, error)
	// ListFunc 描述当前测试场景下的列表桩函数。
	ListFunc func(ctx context.Context, query followeruser.ListQuery) ([]followeruser.FollowerUserView, int64, int, int, error)
	// GetFunc 描述当前测试场景下的详情桩函数。
	GetFunc func(ctx context.Context, userID string) (*followeruser.FollowerUserView, error)
	// UpdateFunc 描述当前测试场景下的更新桩函数。
	UpdateFunc func(ctx context.Context, command followeruser.UpdateCommand) (*followeruser.FollowerUserView, error)
	// DeleteFunc 描述当前测试场景下的删除桩函数。
	DeleteFunc func(ctx context.Context, command followeruser.DeleteCommand) (*followeruser.DeleteResult, error)
}

// Create 负责调用创建桩函数。
func (s *stubRouterFollowerUserService) Create(ctx context.Context, command followeruser.CreateCommand) (*followeruser.FollowerUserView, error) {
	return s.CreateFunc(ctx, command)
}

// List 负责调用列表桩函数。
func (s *stubRouterFollowerUserService) List(ctx context.Context, query followeruser.ListQuery) ([]followeruser.FollowerUserView, int64, int, int, error) {
	return s.ListFunc(ctx, query)
}

// Get 负责调用详情桩函数。
func (s *stubRouterFollowerUserService) Get(ctx context.Context, userID string) (*followeruser.FollowerUserView, error) {
	return s.GetFunc(ctx, userID)
}

// Update 负责调用更新桩函数。
func (s *stubRouterFollowerUserService) Update(ctx context.Context, command followeruser.UpdateCommand) (*followeruser.FollowerUserView, error) {
	return s.UpdateFunc(ctx, command)
}

// Delete 负责调用删除桩函数。
func (s *stubRouterFollowerUserService) Delete(ctx context.Context, command followeruser.DeleteCommand) (*followeruser.DeleteResult, error) {
	return s.DeleteFunc(ctx, command)
}

// TestHealthRoutes 负责验证基础健康检查路由可以正常返回统一响应结构。
func TestHealthRoutes(t *testing.T) {
	// serverConfig 为路由初始化提供测试模式下的服务配置。
	serverConfig := config.ServerConfig{
		Mode: "test",
	}

	// engine 是本次测试实际承载请求的 Gin 引擎。
	engine := New(serverConfig, zap.NewNop(), handler.NewHealthHandler("test-backend"), nil, nil, nil, nil)

	// request 构造了一次对健康检查接口的模拟访问。
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	// recorder 用于捕获 HTTP 处理后的响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("期望状态码为 200，实际为 %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"code":200`) {
		t.Fatalf("期望响应体包含数值型成功码，实际响应为 %s", recorder.Body.String())
	}

	if !strings.Contains(recorder.Body.String(), `"appName":"test-backend"`) {
		t.Fatalf("期望响应体包含小驼峰 appName 字段，实际响应为 %s", recorder.Body.String())
	}
}

// stubRouterAccessAuthenticator 描述用于路由测试的 Access Token 校验桩实现。
type stubRouterAccessAuthenticator struct {
	// AuthenticateFunc 描述当前测试场景下的 Access Token 校验桩函数。
	AuthenticateFunc func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error)
}

// AuthenticateAccessToken 负责调用 Access Token 校验桩函数。
func (s *stubRouterAccessAuthenticator) AuthenticateAccessToken(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
	return s.AuthenticateFunc(ctx, accessToken)
}

// TestProtectedFollowerUserRoutes 负责验证现有后台业务路由已经纳入 Access Token 鉴权保护。
func TestProtectedFollowerUserRoutes(t *testing.T) {
	// serverConfig 为路由初始化提供测试模式下的服务配置。
	serverConfig := config.ServerConfig{
		Mode: "test",
	}

	// authMiddleware 保存当前测试场景使用的后台鉴权中间件。
	authMiddleware := middleware.NewAdminAuthMiddleware(&stubRouterAccessAuthenticator{
		AuthenticateFunc: func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
			return &adminauth.CurrentAdmin{
				AdminUserID: "ADM-1001",
				LoginName:   "admin.root",
				DisplayName: "系统管理员",
				RoleCode:    adminauth.RoleCodeAdmin,
				Status:      adminauth.AdminStatusActive,
			}, nil
		},
	}, zap.NewNop())

	// followerUserHandler 保存当前测试场景使用的小龙虾用户 Handler。
	followerUserHandler := handler.NewFollowerUserHandler(&stubRouterFollowerUserService{
		CreateFunc: nil,
		ListFunc: func(ctx context.Context, query followeruser.ListQuery) ([]followeruser.FollowerUserView, int64, int, int, error) {
			return []followeruser.FollowerUserView{
				{
					UserID:               "U-1201",
					AccountStatus:        "active",
					StrategyStatus:       "enabled",
					BindingStatus:        "pending",
					ResponsibilityDomain: "risk",
					UpdatedAt:            "2026-03-31 15:30:25",
				},
			}, 1, 1, 20, nil
		},
		GetFunc:    nil,
		UpdateFunc: nil,
		DeleteFunc: nil,
	}, zap.NewNop())

	// engine 是本次测试实际承载请求的 Gin 引擎。
	engine := New(
		serverConfig,
		zap.NewNop(),
		handler.NewHealthHandler("test-backend"),
		authMiddleware,
		nil,
		nil,
		followerUserHandler,
	)

	// unauthorizedRequest 构造未携带 Access Token 的列表请求。
	unauthorizedRequest := httptest.NewRequest(http.MethodGet, "/admin/v1/follower-users?pageNum=1&pageSize=20", nil)
	// unauthorizedRecorder 用于捕获未授权请求的响应结果。
	unauthorizedRecorder := httptest.NewRecorder()
	engine.ServeHTTP(unauthorizedRecorder, unauthorizedRequest)

	if unauthorizedRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("期望未登录请求返回 401，实际为 %d", unauthorizedRecorder.Code)
	}

	// authorizedRequest 构造携带合法 Access Token 的列表请求。
	authorizedRequest := httptest.NewRequest(http.MethodGet, "/admin/v1/follower-users?pageNum=1&pageSize=20", nil)
	authorizedRequest.Header.Set("Authorization", "Bearer valid-token")

	// authorizedRecorder 用于捕获已授权请求的响应结果。
	authorizedRecorder := httptest.NewRecorder()
	engine.ServeHTTP(authorizedRecorder, authorizedRequest)

	if authorizedRecorder.Code != http.StatusOK {
		t.Fatalf("期望已登录请求返回 200，实际为 %d", authorizedRecorder.Code)
	}

	if !strings.Contains(authorizedRecorder.Body.String(), `"userId":"U-1201"`) {
		t.Fatalf("期望已登录请求返回小龙虾用户列表，实际响应为 %s", authorizedRecorder.Body.String())
	}
}
