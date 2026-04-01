package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"operation_admin/backend/internal/followeruser"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// stubFollowerUserService 描述用于 handler 测试的桩服务实现。
type stubFollowerUserService struct {
	// CreateFunc 描述创建接口测试时使用的桩函数。
	CreateFunc func(ctx context.Context, command followeruser.CreateCommand) (*followeruser.FollowerUserView, error)
	// ListFunc 描述列表接口测试时使用的桩函数。
	ListFunc func(ctx context.Context, query followeruser.ListQuery) ([]followeruser.FollowerUserView, int64, int, int, error)
	// GetFunc 描述详情接口测试时使用的桩函数。
	GetFunc func(ctx context.Context, userID string) (*followeruser.FollowerUserView, error)
	// UpdateFunc 描述更新接口测试时使用的桩函数。
	UpdateFunc func(ctx context.Context, command followeruser.UpdateCommand) (*followeruser.FollowerUserView, error)
	// DeleteFunc 描述删除接口测试时使用的桩函数。
	DeleteFunc func(ctx context.Context, command followeruser.DeleteCommand) (*followeruser.DeleteResult, error)
}

// Create 负责调用创建桩函数。
func (s *stubFollowerUserService) Create(ctx context.Context, command followeruser.CreateCommand) (*followeruser.FollowerUserView, error) {
	return s.CreateFunc(ctx, command)
}

// List 负责调用列表桩函数。
func (s *stubFollowerUserService) List(ctx context.Context, query followeruser.ListQuery) ([]followeruser.FollowerUserView, int64, int, int, error) {
	return s.ListFunc(ctx, query)
}

// Get 负责调用详情桩函数。
func (s *stubFollowerUserService) Get(ctx context.Context, userID string) (*followeruser.FollowerUserView, error) {
	return s.GetFunc(ctx, userID)
}

// Update 负责调用更新桩函数。
func (s *stubFollowerUserService) Update(ctx context.Context, command followeruser.UpdateCommand) (*followeruser.FollowerUserView, error) {
	return s.UpdateFunc(ctx, command)
}

// Delete 负责调用删除桩函数。
func (s *stubFollowerUserService) Delete(ctx context.Context, command followeruser.DeleteCommand) (*followeruser.DeleteResult, error) {
	return s.DeleteFunc(ctx, command)
}

// TestCreateFollowerUser 负责验证创建接口会返回统一成功响应。
func TestCreateFollowerUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// service 保存当前测试场景使用的桩服务。
	service := &stubFollowerUserService{
		CreateFunc: func(ctx context.Context, command followeruser.CreateCommand) (*followeruser.FollowerUserView, error) {
			return &followeruser.FollowerUserView{
				UserID:               command.UserID,
				AccountStatus:        command.AccountStatus,
				StrategyStatus:       command.StrategyStatus,
				BindingStatus:        command.BindingStatus,
				ResponsibilityDomain: command.ResponsibilityDomain,
				UpdatedAt:            "2026-03-31 15:30:25",
			}, nil
		},
		ListFunc:   nil,
		GetFunc:    nil,
		UpdateFunc: nil,
		DeleteFunc: nil,
	}

	// followerUserHandler 保存当前测试场景使用的 HTTP 处理器。
	followerUserHandler := NewFollowerUserHandler(service, zap.NewNop())

	// engine 用于注册本次测试需要的 Gin 路由。
	engine := gin.New()
	adminGroup := engine.Group("/admin/v1")
	followerUserHandler.RegisterRoutes(adminGroup)

	// requestBody 保存当前测试场景使用的创建请求体。
	requestBody := `{"userId":"U-1201","accountStatus":"active","strategyStatus":"enabled","bindingStatus":"pending","responsibilityDomain":"risk"}`
	// request 构造创建接口的模拟 HTTP 请求。
	request := httptest.NewRequest(http.MethodPost, "/admin/v1/follower-users", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")

	// recorder 用于捕获创建接口的响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("期望状态码为 200，实际为 %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"userId":"U-1201"`) {
		t.Fatalf("期望响应体包含 userId，实际响应为 %s", recorder.Body.String())
	}

	if !strings.Contains(recorder.Body.String(), `"code":200`) {
		t.Fatalf("期望响应体包含统一成功码，实际响应为 %s", recorder.Body.String())
	}
}

// TestListFollowerUsers 负责验证列表接口会返回统一分页结构。
func TestListFollowerUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// service 保存当前测试场景使用的桩服务。
	service := &stubFollowerUserService{
		CreateFunc: nil,
		ListFunc: func(ctx context.Context, query followeruser.ListQuery) ([]followeruser.FollowerUserView, int64, int, int, error) {
			// followerUserViews 保存列表接口返回的桩数据。
			followerUserViews := []followeruser.FollowerUserView{
				{
					UserID:               "U-1201",
					AccountStatus:        "active",
					StrategyStatus:       "enabled",
					BindingStatus:        "pending",
					ResponsibilityDomain: "risk",
					UpdatedAt:            "2026-03-31 15:30:25",
				},
			}

			return followerUserViews, 1, 1, 20, nil
		},
		GetFunc:    nil,
		UpdateFunc: nil,
		DeleteFunc: nil,
	}

	// followerUserHandler 保存当前测试场景使用的 HTTP 处理器。
	followerUserHandler := NewFollowerUserHandler(service, zap.NewNop())

	// engine 用于注册本次测试需要的 Gin 路由。
	engine := gin.New()
	adminGroup := engine.Group("/admin/v1")
	followerUserHandler.RegisterRoutes(adminGroup)

	// request 构造列表接口的模拟 HTTP 请求。
	request := httptest.NewRequest(http.MethodGet, "/admin/v1/follower-users?pageNum=1&pageSize=20", nil)

	// recorder 用于捕获列表接口的响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("期望状态码为 200，实际为 %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), `"list":[`) {
		t.Fatalf("期望响应体包含 list 数组，实际响应为 %s", recorder.Body.String())
	}

	if !strings.Contains(recorder.Body.String(), `"pageNum":1`) {
		t.Fatalf("期望响应体包含 pageNum 字段，实际响应为 %s", recorder.Body.String())
	}
}
