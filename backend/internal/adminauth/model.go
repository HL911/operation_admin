package adminauth

import "time"

// AdminUser 描述 operator_portal.admin_users 表在 Gorm 中的映射结构。
type AdminUser struct {
	// AdminUserID 表示对外暴露的后台账号 ID，例如 ADM-1001。
	AdminUserID string `gorm:"column:admin_user_id;type:varchar(64);primaryKey"`
	// LoginName 表示后台账号登录名，要求全局唯一且统一以小写存储。
	LoginName string `gorm:"column:login_name;type:varchar(32);not null;uniqueIndex:uk_admin_users_login_name"`
	// PasswordHash 表示后台账号密码对应的 bcrypt 哈希值。
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null"`
	// DisplayName 表示后台账号显示名。
	DisplayName string `gorm:"column:display_name;type:varchar(64);not null"`
	// RoleCode 表示后台账号角色编码，当前一期固定为 admin。
	RoleCode string `gorm:"column:role_code;type:varchar(24);not null;default:admin;index:ix_admin_users_role_code;check:ck_admin_users_role_code,role_code = 'admin'"`
	// Status 表示后台账号状态，只允许 active 或 disabled。
	Status string `gorm:"column:status;type:varchar(24);not null;default:active;index:ix_admin_users_status;check:ck_admin_users_status,status IN ('active', 'disabled')"`
	// LastLoginAt 表示最近一次成功登录时间，未登录过时为空。
	LastLoginAt *time.Time `gorm:"column:last_login_at;type:timestamptz"`
	// LastLoginIP 表示最近一次成功登录请求的客户端 IP。
	LastLoginIP string `gorm:"column:last_login_ip;type:varchar(64)"`
	// CreatedAt 表示后台账号创建时间。
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime"`
	// UpdatedAt 表示后台账号最近更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoCreateTime;autoUpdateTime"`
	// IsDeleted 表示后台账号是否已被逻辑删除。
	IsDeleted bool `gorm:"column:is_deleted;type:boolean;not null;default:false;index:ix_admin_users_is_deleted"`
	// RowVersion 表示后台账号记录的乐观锁版本号。
	RowVersion int `gorm:"column:row_version;type:integer;not null;default:1;check:ck_admin_users_row_version,row_version > 0"`
}

// TableName 负责声明 Gorm 需要操作的后台账号真实表名。
func (AdminUser) TableName() string {
	return "operator_portal.admin_users"
}

// AdminRefreshToken 描述 operator_portal.admin_refresh_tokens 表在 Gorm 中的映射结构。
type AdminRefreshToken struct {
	// TokenJTI 表示当前 Refresh Token 记录的唯一标识。
	TokenJTI string `gorm:"column:token_jti;type:varchar(96);primaryKey"`
	// AdminUserID 表示该 Refresh Token 所属的后台账号 ID。
	AdminUserID string `gorm:"column:admin_user_id;type:varchar(64);not null;index:ix_admin_refresh_tokens_admin_user_id"`
	// RefreshTokenHash 表示原始 Refresh Token 的 SHA-256 哈希值。
	RefreshTokenHash string `gorm:"column:refresh_token_hash;type:char(64);not null;uniqueIndex:uk_admin_refresh_tokens_hash"`
	// ExpiresAt 表示 Refresh Token 的过期时间。
	ExpiresAt time.Time `gorm:"column:expires_at;type:timestamptz;not null;index:ix_admin_refresh_tokens_expires_at"`
	// RevokedAt 表示 Refresh Token 被撤销的时间，未撤销时为空。
	RevokedAt *time.Time `gorm:"column:revoked_at;type:timestamptz"`
	// ClientIP 表示签发当前 Refresh Token 时的客户端 IP。
	ClientIP string `gorm:"column:client_ip;type:varchar(64)"`
	// UserAgent 表示签发当前 Refresh Token 时的客户端 User-Agent。
	UserAgent string `gorm:"column:user_agent;type:varchar(255)"`
	// CreatedAt 表示当前 Refresh Token 记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime"`
	// UpdatedAt 表示当前 Refresh Token 记录最近更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoCreateTime;autoUpdateTime"`
}

// TableName 负责声明 Gorm 需要操作的后台 Refresh Token 真实表名。
func (AdminRefreshToken) TableName() string {
	return "operator_portal.admin_refresh_tokens"
}

// AdminLoginAudit 描述 operator_portal.admin_login_audits 表在 Gorm 中的映射结构。
type AdminLoginAudit struct {
	// AdminUserID 表示本次登录尝试对应的后台账号 ID，不存在账号时为空。
	AdminUserID *string `gorm:"column:admin_user_id;type:varchar(64);index:ix_admin_login_audits_admin_user_id"`
	// LoginName 表示本次登录尝试使用的登录名。
	LoginName string `gorm:"column:login_name;type:varchar(32);not null;index:ix_admin_login_audits_login_name"`
	// Result 表示本次登录尝试结果，只允许 success 或 failed。
	Result string `gorm:"column:result;type:varchar(16);not null;check:ck_admin_login_audits_result,result IN ('success', 'failed')"`
	// Reason 表示本次登录尝试结果原因，例如 invalid_credentials。
	Reason string `gorm:"column:reason;type:varchar(64);not null"`
	// ClientIP 表示本次登录尝试的客户端 IP。
	ClientIP string `gorm:"column:client_ip;type:varchar(64)"`
	// UserAgent 表示本次登录尝试携带的 User-Agent。
	UserAgent string `gorm:"column:user_agent;type:varchar(255)"`
	// CreatedAt 表示本次登录审计记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime;index:ix_admin_login_audits_created_at,sort:desc"`
}

// TableName 负责声明 Gorm 需要操作的后台登录审计真实表名。
func (AdminLoginAudit) TableName() string {
	return "operator_portal.admin_login_audits"
}
