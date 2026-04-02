package adminauth

import "time"

const (
	// RoleCodeAdmin 表示当前一期唯一支持的后台角色编码。
	RoleCodeAdmin = "admin"
	// AdminStatusActive 表示后台账号处于可登录、可访问受保护接口的状态。
	AdminStatusActive = "active"
	// AdminStatusDisabled 表示后台账号已被禁用，禁止登录和访问受保护接口。
	AdminStatusDisabled = "disabled"
	// LoginAuditResultSuccess 表示登录审计中的成功结果值。
	LoginAuditResultSuccess = "success"
	// LoginAuditResultFailed 表示登录审计中的失败结果值。
	LoginAuditResultFailed = "failed"
	// LoginAuditReasonInvalidCredentials 表示登录失败原因为登录名不存在或密码错误。
	LoginAuditReasonInvalidCredentials = "invalid_credentials"
	// LoginAuditReasonAccountDisabled 表示登录失败原因为账号已被禁用。
	LoginAuditReasonAccountDisabled = "account_disabled"
	// AccessTokenType 表示 JWT Claims 中的 Access Token 类型标识。
	AccessTokenType = "access"
)

// LoginCommand 描述登录接口传入的业务命令。
type LoginCommand struct {
	// LoginName 表示调用方提交的后台登录名。
	LoginName string
	// Password 表示调用方提交的明文密码。
	Password string
	// ClientIP 表示当前登录请求的客户端 IP。
	ClientIP string
	// UserAgent 表示当前登录请求携带的 User-Agent。
	UserAgent string
}

// RefreshCommand 描述刷新令牌接口传入的业务命令。
type RefreshCommand struct {
	// RefreshToken 表示调用方提交的原始 Refresh Token。
	RefreshToken string
	// ClientIP 表示当前刷新请求的客户端 IP。
	ClientIP string
	// UserAgent 表示当前刷新请求携带的 User-Agent。
	UserAgent string
}

// LogoutCommand 描述登出接口传入的业务命令。
type LogoutCommand struct {
	// AdminUserID 表示当前 Access Token 对应的后台账号 ID。
	AdminUserID string
	// RefreshToken 表示本次需要撤销的原始 Refresh Token。
	RefreshToken string
}

// CreateAdminUserCommand 描述管理员创建后台账号时使用的业务命令。
type CreateAdminUserCommand struct {
	// LoginName 表示新后台账号的登录名。
	LoginName string
	// Password 表示新后台账号的明文密码。
	Password string
	// DisplayName 表示新后台账号的显示名。
	DisplayName string
}

// LoginResult 描述登录成功后返回给接口层的业务结果。
type LoginResult struct {
	// AccessToken 表示签发给调用方的 Access Token。
	AccessToken string `json:"accessToken"`
	// AccessTokenExpiresIn 表示 Access Token 的剩余有效期秒数。
	AccessTokenExpiresIn int64 `json:"accessTokenExpiresIn"`
	// RefreshToken 表示签发给调用方的 Refresh Token。
	RefreshToken string `json:"refreshToken"`
	// RefreshTokenExpiresIn 表示 Refresh Token 的剩余有效期秒数。
	RefreshTokenExpiresIn int64 `json:"refreshTokenExpiresIn"`
	// User 表示当前登录成功的后台账号信息。
	User AdminProfileView `json:"user"`
}

// RefreshResult 描述刷新令牌成功后的业务结果。
type RefreshResult struct {
	// AccessToken 表示新签发的 Access Token。
	AccessToken string `json:"accessToken"`
	// AccessTokenExpiresIn 表示新 Access Token 的剩余有效期秒数。
	AccessTokenExpiresIn int64 `json:"accessTokenExpiresIn"`
	// RefreshToken 表示新签发的 Refresh Token。
	RefreshToken string `json:"refreshToken"`
	// RefreshTokenExpiresIn 表示新 Refresh Token 的剩余有效期秒数。
	RefreshTokenExpiresIn int64 `json:"refreshTokenExpiresIn"`
}

// LogoutResult 描述登出接口统一返回的业务结果。
type LogoutResult struct {
	// Success 表示本次登出是否按约定完成撤销动作或幂等处理。
	Success bool `json:"success"`
}

// AdminProfileView 描述登录成功和当前用户接口返回的后台账号信息。
type AdminProfileView struct {
	// AdminUserID 表示对外暴露的后台账号 ID。
	AdminUserID string `json:"adminUserId"`
	// LoginName 表示后台账号登录名，统一为小写。
	LoginName string `json:"loginName"`
	// DisplayName 表示后台账号显示名。
	DisplayName string `json:"displayName"`
	// RoleCode 表示后台账号角色编码。
	RoleCode string `json:"roleCode"`
	// Status 表示后台账号当前状态。
	Status string `json:"status"`
	// LastLoginAt 表示最近一次成功登录时间，未登录过时返回 null。
	LastLoginAt *string `json:"lastLoginAt"`
}

// CreatedAdminUserView 描述管理员创建后台账号后返回的结果对象。
type CreatedAdminUserView struct {
	// AdminUserID 表示新创建后台账号的对外 ID。
	AdminUserID string `json:"adminUserId"`
	// LoginName 表示新创建后台账号的登录名。
	LoginName string `json:"loginName"`
	// DisplayName 表示新创建后台账号的显示名。
	DisplayName string `json:"displayName"`
	// RoleCode 表示新创建后台账号的角色编码，当前固定为 admin。
	RoleCode string `json:"roleCode"`
	// Status 表示新创建后台账号的状态，当前固定为 active。
	Status string `json:"status"`
	// CreatedAt 表示后台账号创建时间。
	CreatedAt string `json:"createdAt"`
}

// CurrentAdmin 描述通过 Access Token 鉴权后注入请求上下文的后台账号信息。
type CurrentAdmin struct {
	// AdminUserID 表示当前请求所属的后台账号 ID。
	AdminUserID string
	// LoginName 表示当前请求所属的后台账号登录名。
	LoginName string
	// DisplayName 表示当前请求所属的后台账号显示名。
	DisplayName string
	// RoleCode 表示当前请求所属的后台账号角色编码。
	RoleCode string
	// Status 表示当前请求所属的后台账号状态。
	Status string
	// LastLoginAt 表示当前请求所属后台账号最近一次成功登录时间。
	LastLoginAt *time.Time
}
