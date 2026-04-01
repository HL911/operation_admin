package handler

import (
	"time"

	"operation_admin/backend/internal/http/response"

	"github.com/gin-gonic/gin"
)

// HealthHandler 负责输出服务健康状态。
type HealthHandler struct {
	// AppName 保存当前服务名称，用于回填到健康检查响应中。
	AppName string
}

// HealthPayload 描述健康检查接口返回的数据结构。
type HealthPayload struct {
	// AppName 表示当前返回健康状态的服务名称。
	AppName string `json:"appName"`
	// Status 表示当前服务健康状态，当前固定返回 `ok`。
	Status string `json:"status"`
	// CheckTime 表示健康检查生成时间，统一使用 yyyy-MM-dd HH:mm:ss 格式。
	CheckTime string `json:"checkTime"`
}

// NewHealthHandler 负责构建健康检查处理器。
func NewHealthHandler(appName string) *HealthHandler {
	return &HealthHandler{
		AppName: appName,
	}
}

// Check 负责返回服务当前的基础运行状态。
func (h *HealthHandler) Check(ctx *gin.Context) {
	response.Success(ctx, "请求成功", HealthPayload{
		AppName:   h.AppName,
		Status:    "ok",
		CheckTime: time.Now().Format(response.DateTimeLayout),
	})
}
