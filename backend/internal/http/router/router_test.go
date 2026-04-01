package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"operation_admin/backend/internal/config"
	"operation_admin/backend/internal/http/handler"

	"go.uber.org/zap"
)

// TestHealthRoutes 负责验证基础健康检查路由可以正常返回统一响应结构。
func TestHealthRoutes(t *testing.T) {
	// serverConfig 为路由初始化提供测试模式下的服务配置。
	serverConfig := config.ServerConfig{
		Mode: "test",
	}

	// engine 是本次测试实际承载请求的 Gin 引擎。
	engine := New(serverConfig, zap.NewNop(), handler.NewHealthHandler("test-backend"), nil)

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
