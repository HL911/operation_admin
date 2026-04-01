package response

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

const (
	// CodeSuccess 表示请求处理成功。
	CodeSuccess = 200
	// CodeUnauthorized 表示当前请求未通过身份校验。
	CodeUnauthorized = 401
	// CodeForbidden 表示当前请求缺少权限或缺少额外校验。
	CodeForbidden = 403
	// CodeNotFound 表示请求的资源不存在。
	CodeNotFound = 404
	// CodeMethodNotAllowed 表示当前资源不支持请求方法。
	CodeMethodNotAllowed = 405
	// CodeInternalError 表示服务端出现未预期错误。
	CodeInternalError = 500
	// CodeValidationFailed 表示业务层参数校验失败。
	CodeValidationFailed = 40001
	// DateTimeLayout 定义接口中时间字符串字段统一使用的输出格式。
	DateTimeLayout = "2006-01-02 15:04:05"
)

const (
	// defaultSuccessMessage 定义普通成功响应的默认提示文案。
	defaultSuccessMessage = "请求成功"
	// defaultPageSuccessMessage 定义分页查询成功时的默认提示文案。
	defaultPageSuccessMessage = "查询成功"
	// defaultValidationMessage 定义参数校验失败时的默认提示文案。
	defaultValidationMessage = "参数校验失败"
	// defaultUnauthorizedMessage 定义登录失效时的默认提示文案。
	defaultUnauthorizedMessage = "登录已过期，请重新登录"
	// defaultForbiddenMessage 定义权限不足或缺少二次校验时的默认提示文案。
	defaultForbiddenMessage = "当前操作需要二次校验"
	// defaultNotFoundMessage 定义资源不存在时的默认提示文案。
	defaultNotFoundMessage = "请求的资源不存在"
	// defaultMethodNotAllowedMessage 定义请求方法不允许时的默认提示文案。
	defaultMethodNotAllowedMessage = "请求方法不被允许"
	// defaultInternalErrorMessage 定义服务内部错误时的默认提示文案。
	defaultInternalErrorMessage = "服务器内部错误，请稍后重试"
)

// Envelope 描述统一的 JSON 响应结构。
type Envelope struct {
	// Code 表示业务响应码，成功时固定为 200，失败时返回 HTTP 状态码或自定义业务错误码。
	Code int `json:"code"`
	// Message 表示给调用方阅读的响应说明文字。
	Message string `json:"message"`
	// Data 表示响应携带的业务载荷，无数据时必须显式返回 null。
	Data any `json:"data"`
	// Timestamp 表示响应生成时间，统一使用 Unix 毫秒时间戳。
	Timestamp int64 `json:"timestamp"`
}

// PageData 描述分页列表接口统一返回的数据结构。
type PageData[T any] struct {
	// List 表示当前页的数据列表，即使为空也应返回空数组而不是 null。
	List []T `json:"list"`
	// Total 表示满足当前筛选条件的总记录数。
	Total int64 `json:"total"`
	// PageNum 表示当前返回的数据页码。
	PageNum int `json:"pageNum"`
	// PageSize 表示当前页的分页大小。
	PageSize int `json:"pageSize"`
	// Pages 表示总页数，按总数和分页大小向上取整。
	Pages int `json:"pages"`
}

// Success 负责输出对象或普通成功响应。
func Success(ctx *gin.Context, message string, data any) {
	// finalMessage 保存最终写回客户端的成功提示文案。
	finalMessage := message
	if finalMessage == "" {
		finalMessage = defaultSuccessMessage
	}

	write(ctx, http.StatusOK, newEnvelope(CodeSuccess, finalMessage, data))
}

// SuccessPage 负责输出统一的分页列表成功响应。
func SuccessPage[T any](ctx *gin.Context, message string, list []T, total int64, pageNum, pageSize int) {
	// finalMessage 保存最终写回客户端的分页成功提示文案。
	finalMessage := message
	if finalMessage == "" {
		finalMessage = defaultPageSuccessMessage
	}

	// pageData 描述列表接口约定的分页结构。
	pageData := PageData[T]{
		List:     normalizeSlice(list),
		Total:    total,
		PageNum:  pageNum,
		PageSize: pageSize,
		Pages:    calculatePages(total, pageSize),
	}

	write(ctx, http.StatusOK, newEnvelope(CodeSuccess, finalMessage, pageData))
}

// Fail 负责输出失败响应，并中断后续处理链。
func Fail(ctx *gin.Context, statusCode int, code int, message string) {
	// finalCode 保存最终返回给前端的业务错误码。
	finalCode := code
	if finalCode == 0 {
		finalCode = statusCode
	}

	// finalMessage 保存最终写回客户端的错误提示文案。
	finalMessage := message
	if finalMessage == "" {
		finalMessage = http.StatusText(statusCode)
	}

	write(ctx, statusCode, newEnvelope(finalCode, finalMessage, nil))
	ctx.Abort()
}

// ValidationFailed 负责输出参数校验失败响应。
func ValidationFailed(ctx *gin.Context, detail string) {
	// finalMessage 保存拼接后的参数校验错误提示。
	finalMessage := defaultValidationMessage
	if detail != "" {
		finalMessage = fmt.Sprintf("%s：%s", defaultValidationMessage, detail)
	}

	Fail(ctx, http.StatusBadRequest, CodeValidationFailed, finalMessage)
}

// Unauthorized 负责输出登录失效响应。
func Unauthorized(ctx *gin.Context, message string) {
	// finalMessage 保存登录失效时最终返回的错误文案。
	finalMessage := message
	if finalMessage == "" {
		finalMessage = defaultUnauthorizedMessage
	}

	Fail(ctx, http.StatusUnauthorized, CodeUnauthorized, finalMessage)
}

// Forbidden 负责输出权限不足或缺少二次校验响应。
func Forbidden(ctx *gin.Context, message string) {
	// finalMessage 保存权限相关错误的最终提示文案。
	finalMessage := message
	if finalMessage == "" {
		finalMessage = defaultForbiddenMessage
	}

	Fail(ctx, http.StatusForbidden, CodeForbidden, finalMessage)
}

// NotFound 负责输出资源不存在响应。
func NotFound(ctx *gin.Context, message string) {
	// finalMessage 保存资源不存在时最终返回的提示文案。
	finalMessage := message
	if finalMessage == "" {
		finalMessage = defaultNotFoundMessage
	}

	Fail(ctx, http.StatusNotFound, CodeNotFound, finalMessage)
}

// MethodNotAllowed 负责输出请求方法不被允许响应。
func MethodNotAllowed(ctx *gin.Context, message string) {
	// finalMessage 保存请求方法不匹配时最终返回的提示文案。
	finalMessage := message
	if finalMessage == "" {
		finalMessage = defaultMethodNotAllowedMessage
	}

	Fail(ctx, http.StatusMethodNotAllowed, CodeMethodNotAllowed, finalMessage)
}

// InternalError 负责输出服务器内部错误响应。
func InternalError(ctx *gin.Context, message string) {
	// finalMessage 保存服务内部错误时最终返回的提示文案。
	finalMessage := message
	if finalMessage == "" {
		finalMessage = defaultInternalErrorMessage
	}

	Fail(ctx, http.StatusInternalServerError, CodeInternalError, finalMessage)
}

// newEnvelope 负责创建带有统一时间戳的响应对象。
func newEnvelope(code int, message string, data any) Envelope {
	return Envelope{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}

// write 负责统一执行 Sonic 序列化与响应写回。
func write(ctx *gin.Context, statusCode int, payload Envelope) {
	// payloadBytes 保存 Sonic 序列化后的最终响应内容。
	payloadBytes, err := sonic.Marshal(payload)
	if err != nil {
		// fallbackBody 保存 Sonic 序列化失败时写回客户端的兜底 JSON 响应。
		fallbackBody := fmt.Sprintf(
			`{"code":%d,"message":"%s","data":null,"timestamp":%d}`,
			CodeInternalError,
			defaultInternalErrorMessage,
			time.Now().UnixMilli(),
		)

		ctx.Header("Content-Type", "application/json; charset=utf-8")
		ctx.Status(http.StatusInternalServerError)
		_, _ = ctx.Writer.Write([]byte(fallbackBody))
		return
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Status(statusCode)
	_, _ = ctx.Writer.Write(payloadBytes)
}

// calculatePages 负责根据总记录数和分页大小计算总页数。
func calculatePages(total int64, pageSize int) int {
	if total <= 0 || pageSize <= 0 {
		return 0
	}

	return int((total + int64(pageSize) - 1) / int64(pageSize))
}

// normalizeSlice 负责把空切片标准化为 JSON 数组而不是 null。
func normalizeSlice[T any](items []T) []T {
	if items == nil {
		return make([]T, 0)
	}

	return items
}
