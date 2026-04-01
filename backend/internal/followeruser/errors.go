package followeruser

import "errors"

const (
	// ErrorCodeFollowerUserDuplicated 表示创建时遇到重复 userId 的业务错误码。
	ErrorCodeFollowerUserDuplicated = 40002
	// ErrorCodeFollowerUserVersionConflict 表示更新或删除时版本号冲突的业务错误码。
	ErrorCodeFollowerUserVersionConflict = 40003
)

var (
	// ErrFollowerUserNotFound 表示当前 userId 对应的小龙虾用户不存在或已被逻辑删除。
	ErrFollowerUserNotFound = errors.New("小龙虾用户不存在")
	// ErrFollowerUserDuplicated 表示创建时 userId 已存在。
	ErrFollowerUserDuplicated = errors.New("小龙虾用户已存在")
	// ErrFollowerUserVersionConflict 表示 rowVersion 与数据库中的当前版本不一致。
	ErrFollowerUserVersionConflict = errors.New("小龙虾用户版本冲突")
)

// ValidationError 描述业务层发现的可直接返回给调用方的校验错误。
type ValidationError struct {
	// Message 保存需要返回给前端的具体校验提示。
	Message string
}

// Error 负责把校验错误转换为标准 error 文本。
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError 负责创建统一的业务校验错误对象。
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}
