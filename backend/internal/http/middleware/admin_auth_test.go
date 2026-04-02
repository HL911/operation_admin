package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"operation_admin/backend/internal/adminauth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// stubAccessAuthenticator 描述用于鉴权中间件测试的 Access Token 校验桩实现。
type stubAccessAuthenticator struct {
	// AuthenticateFunc 描述当前测试场景下的 Access Token 校验桩函数。
	AuthenticateFunc func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error)
}

// AuthenticateAccessToken 负责调用 Access Token 校验桩函数。
func (s *stubAccessAuthenticator) AuthenticateAccessToken(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
	return s.AuthenticateFunc(ctx, accessToken)
}

// TestRequireAccessToken 负责验证后台鉴权中间件的核心校验路径。
func TestRequireAccessToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// testCases 保存本次测试需要覆盖的后台鉴权中间件场景。
	testCases := []struct {
		// Name 表示当前子测试名称。
		Name string
		// AuthorizationHeader 表示当前请求头中的 Authorization 内容。
		AuthorizationHeader string
		// AuthenticateFunc 表示当前场景使用的 Access Token 校验桩函数。
		AuthenticateFunc func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error)
		// WantStatusCode 表示当前场景期望返回的 HTTP 状态码。
		WantStatusCode int
		// WantBodyContains 表示当前场景期望响应体包含的关键片段。
		WantBodyContains string
	}{
		{
			Name:                "缺少 Authorization 头返回 401",
			AuthorizationHeader: "",
			AuthenticateFunc: func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
				return nil, nil
			},
			WantStatusCode:   http.StatusUnauthorized,
			WantBodyContains: `"code":401`,
		},
		{
			Name:                "错误 Bearer 格式返回 401",
			AuthorizationHeader: "Token test-token",
			AuthenticateFunc: func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
				return nil, nil
			},
			WantStatusCode:   http.StatusUnauthorized,
			WantBodyContains: `"code":401`,
		},
		{
			Name:                "失效 token 返回 401",
			AuthorizationHeader: "Bearer invalid-token",
			AuthenticateFunc: func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
				return nil, adminauth.ErrAccessTokenInvalid
			},
			WantStatusCode:   http.StatusUnauthorized,
			WantBodyContains: `"code":401`,
		},
		{
			Name:                "禁用账号返回 403",
			AuthorizationHeader: "Bearer disabled-token",
			AuthenticateFunc: func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
				return nil, adminauth.ErrAdminUserDisabled
			},
			WantStatusCode:   http.StatusForbidden,
			WantBodyContains: `"code":403`,
		},
		{
			Name:                "合法 token 可以放行并注入当前账号信息",
			AuthorizationHeader: "Bearer valid-token",
			AuthenticateFunc: func(ctx context.Context, accessToken string) (*adminauth.CurrentAdmin, error) {
				// lastLoginAt 保存当前测试场景构造的最近登录时间。
				lastLoginAt := time.Date(2026, 4, 2, 10, 0, 0, 0, time.Local)

				return &adminauth.CurrentAdmin{
					AdminUserID: "ADM-1001",
					LoginName:   "admin.root",
					DisplayName: "系统管理员",
					RoleCode:    adminauth.RoleCodeAdmin,
					Status:      adminauth.AdminStatusActive,
					LastLoginAt: &lastLoginAt,
				}, nil
			},
			WantStatusCode:   http.StatusOK,
			WantBodyContains: `"loginName":"admin.root"`,
		},
	}

	for _, testCase := range testCases {
		// testCase 保存当前循环执行的测试用例副本，避免闭包复用同一地址。
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			// authMiddleware 保存当前测试场景使用的后台鉴权中间件实例。
			authMiddleware := NewAdminAuthMiddleware(&stubAccessAuthenticator{
				AuthenticateFunc: testCase.AuthenticateFunc,
			}, zap.NewNop())

			// engine 用于承载当前测试场景的 Gin 路由和中间件。
			engine := gin.New()
			protectedGroup := engine.Group("/admin/v1")
			protectedGroup.Use(authMiddleware.RequireAccessToken())
			protectedGroup.GET("/protected", func(ctx *gin.Context) {
				// currentAdmin 保存当前请求上下文中的后台账号信息。
				currentAdmin, ok := CurrentAdmin(ctx)
				if !ok {
					ctx.String(http.StatusInternalServerError, "missing current admin")
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"loginName": currentAdmin.LoginName,
				})
			})

			// request 构造当前测试场景的模拟 HTTP 请求。
			request := httptest.NewRequest(http.MethodGet, "/admin/v1/protected", nil)
			if testCase.AuthorizationHeader != "" {
				request.Header.Set("Authorization", testCase.AuthorizationHeader)
			}

			// recorder 用于捕获中间件处理后的响应结果。
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, request)

			if recorder.Code != testCase.WantStatusCode {
				t.Fatalf("期望状态码为 %d，实际为 %d", testCase.WantStatusCode, recorder.Code)
			}

			if !strings.Contains(recorder.Body.String(), testCase.WantBodyContains) {
				t.Fatalf("期望响应体包含 %s，实际响应为 %s", testCase.WantBodyContains, recorder.Body.String())
			}
		})
	}
}
