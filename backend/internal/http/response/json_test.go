package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

// TestSuccess 负责验证对象成功响应符合统一结构约定。
func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// engine 用于注册测试场景下的最小 Gin 路由。
	engine := gin.New()
	engine.GET("/success", func(ctx *gin.Context) {
		// userData 描述对象查询成功时返回的业务数据。
		userData := gin.H{
			"id":         10001,
			"username":   "admin",
			"phone":      "13800138000",
			"createTime": "2026-03-31 15:30:25",
		}

		Success(ctx, "请求成功", userData)
	})

	// request 构造对象成功响应的模拟请求。
	request := httptest.NewRequest(http.MethodGet, "/success", nil)
	// recorder 用于捕获 HTTP 响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("期望 HTTP 状态码为 200，实际为 %d", recorder.Code)
	}

	// responseBody 保存解析后的统一响应结构。
	responseBody := decodeBody(t, recorder.Body.Bytes())

	if responseBody["code"] != float64(CodeSuccess) {
		t.Fatalf("期望业务码为 200，实际为 %#v", responseBody["code"])
	}

	if responseBody["message"] != "请求成功" {
		t.Fatalf("期望 message 为 请求成功，实际为 %#v", responseBody["message"])
	}

	// dataValue 保存对象成功响应中的业务数据部分。
	dataValue, ok := responseBody["data"].(map[string]any)
	if !ok {
		t.Fatalf("期望 data 为对象，实际为 %#v", responseBody["data"])
	}

	if dataValue["createTime"] != "2026-03-31 15:30:25" {
		t.Fatalf("期望 createTime 保持 yyyy-MM-dd HH:mm:ss 格式，实际为 %#v", dataValue["createTime"])
	}

	if responseBody["timestamp"] == nil {
		t.Fatal("期望统一响应包含毫秒级 timestamp 字段")
	}
}

// TestSuccessPage 负责验证分页成功响应符合统一结构约定。
func TestSuccessPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// engine 用于注册测试场景下的最小 Gin 路由。
	engine := gin.New()
	engine.GET("/page", func(ctx *gin.Context) {
		// orderList 描述分页查询成功时返回的列表数据。
		orderList := []gin.H{
			{
				"id":     1,
				"name":   "订单A",
				"status": 1,
			},
			{
				"id":     2,
				"name":   "订单B",
				"status": 2,
			},
		}

		SuccessPage(ctx, "查询成功", orderList, 128, 1, 10)
	})

	// request 构造分页成功响应的模拟请求。
	request := httptest.NewRequest(http.MethodGet, "/page", nil)
	// recorder 用于捕获 HTTP 响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	// responseBody 保存解析后的统一响应结构。
	responseBody := decodeBody(t, recorder.Body.Bytes())
	// dataValue 保存分页成功响应中的 data 部分。
	dataValue, ok := responseBody["data"].(map[string]any)
	if !ok {
		t.Fatalf("期望分页 data 为对象，实际为 %#v", responseBody["data"])
	}

	if dataValue["pageNum"] != float64(1) {
		t.Fatalf("期望 pageNum 为 1，实际为 %#v", dataValue["pageNum"])
	}

	if dataValue["pageSize"] != float64(10) {
		t.Fatalf("期望 pageSize 为 10，实际为 %#v", dataValue["pageSize"])
	}

	if dataValue["pages"] != float64(13) {
		t.Fatalf("期望总页数为 13，实际为 %#v", dataValue["pages"])
	}

	// listValue 保存分页返回的列表字段。
	listValue, ok := dataValue["list"].([]any)
	if !ok {
		t.Fatalf("期望 list 为数组，实际为 %#v", dataValue["list"])
	}

	if len(listValue) != 2 {
		t.Fatalf("期望 list 长度为 2，实际为 %d", len(listValue))
	}
}

// TestValidationFailed 负责验证失败响应会固定返回 data:null 和业务错误码。
func TestValidationFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// engine 用于注册测试场景下的最小 Gin 路由。
	engine := gin.New()
	engine.GET("/validation", func(ctx *gin.Context) {
		ValidationFailed(ctx, "手机号格式错误")
	})

	// request 构造参数校验失败的模拟请求。
	request := httptest.NewRequest(http.MethodGet, "/validation", nil)
	// recorder 用于捕获 HTTP 响应结果。
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("期望 HTTP 状态码为 400，实际为 %d", recorder.Code)
	}

	// responseBody 保存解析后的统一响应结构。
	responseBody := decodeBody(t, recorder.Body.Bytes())

	if responseBody["code"] != float64(CodeValidationFailed) {
		t.Fatalf("期望业务错误码为 40001，实际为 %#v", responseBody["code"])
	}

	if responseBody["message"] != "参数校验失败：手机号格式错误" {
		t.Fatalf("期望错误文案符合规范，实际为 %#v", responseBody["message"])
	}

	if dataValue, exists := responseBody["data"]; !exists || dataValue != nil {
		t.Fatalf("期望失败响应显式返回 data:null，实际为 %#v", responseBody["data"])
	}
}

// decodeBody 负责把测试响应体解析为通用 map 结构，便于断言 JSON 字段。
func decodeBody(t *testing.T, rawBody []byte) map[string]any {
	t.Helper()

	// decodedBody 保存反序列化后的 JSON 响应体。
	decodedBody := make(map[string]any)
	if err := sonic.Unmarshal(rawBody, &decodedBody); err != nil {
		t.Fatalf("解析响应体失败: %v", err)
	}

	return decodedBody
}
