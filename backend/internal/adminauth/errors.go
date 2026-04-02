package adminauth

import "errors"

const (
	// ErrorCodeAdminUserDuplicated 表示创建后台账号时遇到重复 loginName 的业务错误码。
	ErrorCodeAdminUserDuplicated = 40010
)

var (
	// ErrAdminUserDuplicated 表示当前 loginName 已被其他后台账号占用。
	ErrAdminUserDuplicated = errors.New("后台账号已存在")
	// ErrInvalidCredentials 表示登录名不存在或密码不正确。
	ErrInvalidCredentials = errors.New("登录名或密码错误")
	// ErrAccessTokenInvalid 表示 Access Token 无效、过期或对应账号不可用。
	ErrAccessTokenInvalid = errors.New("access token 无效")
	// ErrRefreshTokenInvalid 表示 Refresh Token 无效、过期、伪造或已撤销。
	ErrRefreshTokenInvalid = errors.New("refresh token 无效")
	// ErrAdminUserDisabled 表示当前后台账号已被禁用。
	ErrAdminUserDisabled = errors.New("后台账号已被禁用")
	// ErrAdminPermissionDenied 表示当前后台账号没有执行目标操作的权限。
	ErrAdminPermissionDenied = errors.New("当前后台账号无权执行该操作")
)

// ValidationError 描述后台鉴权模块中的可直接返回给调用方的参数校验错误。
type ValidationError struct {
	// Message 保存需要返回给前端的具体校验提示。
	Message string
}

// Error 负责把参数校验错误转换为标准 error 文本。
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError 负责创建统一的参数校验错误对象。
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}
