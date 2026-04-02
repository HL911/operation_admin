package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"operation_admin/backend/internal/adminauth"
	"operation_admin/backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// stubAuthService 描述用于后台鉴权 Handler 测试的业务桩实现。
type stubAuthService struct {
	// LoginFunc 描述当前测试场景下的登录桩函数。
	LoginFunc func(ctx context.Context, command adminauth.LoginCommand) (*adminauth.LoginResult, error)
	// RefreshFunc 描述当前测试场景下的刷新桩函数。
	RefreshFunc func(ctx context.Context, command adminauth.RefreshCommand) (*adminauth.RefreshResult, error)
	// LogoutFunc 描述当前测试场景下的登出桩函数。
	LogoutFunc func(ctx context.Context, command adminauth.LogoutCommand) (*adminauth.LogoutResult, error)
}

// Login 负责调用登录桩函数。
func (s *stubAuthService) Login(ctx context.Context, command adminauth.LoginCommand) (*adminauth.LoginResult, error) {
	return s.LoginFunc(ctx, command)
}

// Refresh 负责调用刷新桩函数。
func (s *stubAuthService) Refresh(ctx context.Context, command adminauth.RefreshCommand) (*adminauth.RefreshResult, error) {
	return s.RefreshFunc(ctx, command)
}

// Logout 负责调用登出桩函数。
func (s *stubAuthService) Logout(ctx context.Context, command adminauth.LogoutCommand) (*adminauth.LogoutResult, error) {
	return s.LogoutFunc(ctx, command)
}

// stubAuthAccessAuthenticator 描述用于后台鉴权 Handler 测试的 Access Token 校验桩实现。
type stubAuthAccessAuthenticator struct {
	// AuthenticateFunc 描述当前测试场景下的 Access Token 校验桩函数。
	AuthenticateFunc func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error)
}

// AuthenticateAccessToken 负责调用 Access Token 校验桩函数。
func (s *stubAuthAccessAuthenticator) AuthenticateAccessToken(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
	return s.AuthenticateFunc(ctx, accessToken)
}

// TestAuthHandlerLogin 负责验证登录接口会返回统一成功响应。
func TestAuthHandlerLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// authHandler 保存当前测试场景使用的后台鉴权 Handler。
	authHandler := NewAuthHandler(&stubAuthService{
		LoginFunc: func(ctx context.Context, command adminauth.LoginCommand) (*adminauth.LoginResult, error) {
			// lastLoginAt 保存当前测试场景构造的最近登录时间。
			lastLoginAt := "2026-04-02 10:20:30"

			return &adminauth.LoginResult{
				AccessToken:           "access-token",
				AccessTokenExpiresIn:  7200,
				RefreshToken:          "refresh-token",
				RefreshTokenExpiresIn: 604800,
				User: adminauth.AdminProfileView{
					AdminUserID: "ADM-1001",
					LoginName:   "admin.root",
					DisplayName: "系统管理员",
					RoleCode:    adminauth.RoleCodeAdmin,
					Status:      adminauth.AdminStatusActive,
					LastLoginAt: &lastLoginAt,
				},
			}, nil
		},
		RefreshFunc: nil,
		LogoutFunc:  nil,
	}, zap.NewNop())

	// engine 用于注册当前测试场景需要的公开后台路由。
	engine := gin.New()
	adminGroup := engine.Group("/admin/v1")
	authHandler.RegisterPublicRoutes(adminGroup)

	// requestBody 保存当前测试场景使用的登录请求体。
	requestBody := `{"loginName":"admin.root","password":"12345678"}`
	// request 构造登录接口的模拟 HTTP 请求。
	request := httptest.NewRequest(http.MethodPost, "/admin/v1/auth/login", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")

	// recorder 用于捕获登录接口的响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("期望状态码为 200，实际为 %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"accessToken":"access-token"`) {
		t.Fatalf("期望响应体包含 accessToken，实际响应为 %s", recorder.Body.String())
	}

	if !strings.Contains(recorder.Body.String(), `"loginName":"admin.root"`) {
		t.Fatalf("期望响应体包含 user.loginName，实际响应为 %s", recorder.Body.String())
	}
}

// TestAuthHandlerLoginInvalidCredentials 负责验证登录名或密码错误时会返回 401。
func TestAuthHandlerLoginInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// authHandler 保存当前测试场景使用的后台鉴权 Handler。
	authHandler := NewAuthHandler(&stubAuthService{
		LoginFunc: func(ctx context.Context, command adminauth.LoginCommand) (*adminauth.LoginResult, error) {
			return nil, adminauth.ErrInvalidCredentials
		},
		RefreshFunc: nil,
		LogoutFunc:  nil,
	}, zap.NewNop())

	// engine 用于注册当前测试场景需要的公开后台路由。
	engine := gin.New()
	adminGroup := engine.Group("/admin/v1")
	authHandler.RegisterPublicRoutes(adminGroup)

	// requestBody 保存当前测试场景使用的登录请求体。
	requestBody := `{"loginName":"admin.root","password":"wrong-pass"}`
	// request 构造登录接口的模拟 HTTP 请求。
	request := httptest.NewRequest(http.MethodPost, "/admin/v1/auth/login", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")

	// recorder 用于捕获登录接口的响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("期望状态码为 401，实际为 %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `登录名或密码错误`) {
		t.Fatalf("期望响应体包含登录失败文案，实际响应为 %s", recorder.Body.String())
	}
}

// TestAuthHandlerMe 负责验证 me 接口会返回中间件注入的当前后台账号信息。
func TestAuthHandlerMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// lastLoginAt 保存当前测试场景构造的最近登录时间。
	lastLoginAt := time.Date(2026, 4, 2, 11, 22, 33, 0, time.Local)

	// authHandler 保存当前测试场景使用的后台鉴权 Handler。
	authHandler := NewAuthHandler(&stubAuthService{
		LoginFunc:   nil,
		RefreshFunc: nil,
		LogoutFunc:  nil,
	}, zap.NewNop())

	// authMiddleware 保存当前测试场景使用的后台鉴权中间件。
	authMiddleware := middleware.NewAdminAuthMiddleware(&stubAuthAccessAuthenticator{
		AuthenticateFunc: func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
			return &adminauth.CurrentAdmin{
				AdminUserID: "ADM-1001",
				LoginName:   "admin.root",
				DisplayName: "系统管理员",
				RoleCode:    adminauth.RoleCodeAdmin,
				Status:      adminauth.AdminStatusActive,
				LastLoginAt: &lastLoginAt,
			}, nil
		},
	}, zap.NewNop())

	// engine 用于注册当前测试场景需要的受保护后台路由。
	engine := gin.New()
	adminGroup := engine.Group("/admin/v1")
	adminGroup.Use(authMiddleware.RequireAccessToken())
	authHandler.RegisterProtectedRoutes(adminGroup)

	// request 构造 me 接口的模拟 HTTP 请求。
	request := httptest.NewRequest(http.MethodGet, "/admin/v1/auth/me", nil)
	request.Header.Set("Authorization", "Bearer valid-token")

	// recorder 用于捕获 me 接口的响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("期望状态码为 200，实际为 %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"adminUserId":"ADM-1001"`) {
		t.Fatalf("期望响应体包含 adminUserId，实际响应为 %s", recorder.Body.String())
	}

	if !strings.Contains(recorder.Body.String(), `"lastLoginAt":"2026-04-02 11:22:33"`) {
		t.Fatalf("期望响应体包含 lastLoginAt，实际响应为 %s", recorder.Body.String())
	}
}
