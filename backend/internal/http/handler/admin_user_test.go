package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"operation_admin/backend/internal/adminauth"
	"operation_admin/backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// stubAdminUserService 描述用于后台开发者账号 Handler 测试的业务桩实现。
type stubAdminUserService struct {
	// CreateFunc 描述当前测试场景下的创建后台账号桩函数。
	CreateFunc func(ctx context.Context, operator adminauth.CurrentAdmin, command adminauth.CreateAdminUserCommand) (*adminauth.CreatedAdminUserView, error)
}

// CreateAdminUser 负责调用创建后台账号桩函数。
func (s *stubAdminUserService) CreateAdminUser(ctx context.Context, operator adminauth.CurrentAdmin, command adminauth.CreateAdminUserCommand) (*adminauth.CreatedAdminUserView, error) {
	return s.CreateFunc(ctx, operator, command)
}

// stubAdminUserAccessAuthenticator 描述用于后台开发者账号 Handler 测试的 Access Token 校验桩实现。
type stubAdminUserAccessAuthenticator struct {
	// AuthenticateFunc 描述当前测试场景下的 Access Token 校验桩函数。
	AuthenticateFunc func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error)
}

// AuthenticateAccessToken 负责调用 Access Token 校验桩函数。
func (s *stubAdminUserAccessAuthenticator) AuthenticateAccessToken(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
	return s.AuthenticateFunc(ctx, accessToken)
}

// TestAdminUserHandlerCreate 负责验证管理员创建后台账号接口会返回统一成功响应。
func TestAdminUserHandlerCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// adminUserHandler 保存当前测试场景使用的后台开发者账号 Handler。
	adminUserHandler := NewAdminUserHandler(&stubAdminUserService{
		CreateFunc: func(ctx context.Context, operator adminauth.CurrentAdmin, command adminauth.CreateAdminUserCommand) (*adminauth.CreatedAdminUserView, error) {
			return &adminauth.CreatedAdminUserView{
				AdminUserID: "ADM-1002",
				LoginName:   "dev.user",
				DisplayName: "开发同学",
				RoleCode:    adminauth.RoleCodeAdmin,
				Status:      adminauth.AdminStatusActive,
				CreatedAt:   "2026-04-02 12:30:45",
			}, nil
		},
	}, zap.NewNop())

	// authMiddleware 保存当前测试场景使用的后台鉴权中间件。
	authMiddleware := middleware.NewAdminAuthMiddleware(&stubAdminUserAccessAuthenticator{
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

	// engine 用于注册当前测试场景需要的受保护后台路由。
	engine := gin.New()
	adminGroup := engine.Group("/admin/v1")
	adminGroup.Use(authMiddleware.RequireAccessToken())
	adminUserHandler.RegisterRoutes(adminGroup)

	// requestBody 保存当前测试场景使用的创建请求体。
	requestBody := `{"loginName":"dev.user","password":"12345678","displayName":"开发同学"}`
	// request 构造创建后台账号接口的模拟 HTTP 请求。
	request := httptest.NewRequest(http.MethodPost, "/admin/v1/admin-users", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer valid-token")

	// recorder 用于捕获创建后台账号接口的响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("期望状态码为 200，实际为 %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"adminUserId":"ADM-1002"`) {
		t.Fatalf("期望响应体包含 adminUserId，实际响应为 %s", recorder.Body.String())
	}
}

// TestAdminUserHandlerCreateDuplicate 负责验证重复 loginName 会返回 40010。
func TestAdminUserHandlerCreateDuplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// adminUserHandler 保存当前测试场景使用的后台开发者账号 Handler。
	adminUserHandler := NewAdminUserHandler(&stubAdminUserService{
		CreateFunc: func(ctx context.Context, operator adminauth.CurrentAdmin, command adminauth.CreateAdminUserCommand) (*adminauth.CreatedAdminUserView, error) {
			return nil, adminauth.ErrAdminUserDuplicated
		},
	}, zap.NewNop())

	// authMiddleware 保存当前测试场景使用的后台鉴权中间件。
	authMiddleware := middleware.NewAdminAuthMiddleware(&stubAdminUserAccessAuthenticator{
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

	// engine 用于注册当前测试场景需要的受保护后台路由。
	engine := gin.New()
	adminGroup := engine.Group("/admin/v1")
	adminGroup.Use(authMiddleware.RequireAccessToken())
	adminUserHandler.RegisterRoutes(adminGroup)

	// requestBody 保存当前测试场景使用的创建请求体。
	requestBody := `{"loginName":"dev.user","password":"12345678","displayName":"开发同学"}`
	// request 构造创建后台账号接口的模拟 HTTP 请求。
	request := httptest.NewRequest(http.MethodPost, "/admin/v1/admin-users", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer valid-token")

	// recorder 用于捕获创建后台账号接口的响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("期望状态码为 409，实际为 %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"code":40010`) {
		t.Fatalf("期望响应体包含 40010 业务错误码，实际响应为 %s", recorder.Body.String())
	}
}
