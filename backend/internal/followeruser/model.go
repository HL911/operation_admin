package followeruser

import "time"

// FollowerUser 描述 operator_portal.follower_users 表在 Gorm 中的映射结构。
type FollowerUser struct {
	// ID 表示表主键 ID，由 PostgreSQL identity 自动生成。
	ID int64 `gorm:"column:id;primaryKey;autoIncrement"`
	// UserID 表示业务侧用户 ID，在未删除记录中必须唯一。
	UserID string `gorm:"column:user_id;type:varchar(64);not null;uniqueIndex:uk_follower_users_user_id_active,where:is_deleted = false"`
	// AccountStatus 表示账户状态，只允许 active 或 disabled。
	AccountStatus string `gorm:"column:account_status;type:varchar(24);not null;default:active;index:ix_follower_users_account_status;check:ck_follower_users_account_status,account_status IN ('active', 'disabled')"`
	// StrategyStatus 表示策略状态，只允许 enabled 或 disabled。
	StrategyStatus string `gorm:"column:strategy_status;type:varchar(24);not null;default:disabled;index:ix_follower_users_strategy_status;check:ck_follower_users_strategy_status,strategy_status IN ('enabled', 'disabled')"`
	// BindingStatus 表示绑定状态，只允许 pending、bound 或 unbound。
	BindingStatus string `gorm:"column:binding_status;type:varchar(24);not null;default:unbound;index:ix_follower_users_binding_status;check:ck_follower_users_binding_status,binding_status IN ('pending', 'bound', 'unbound')"`
	// ResponsibilityDomain 表示责任域，例如 risk 或运营自定义值。
	ResponsibilityDomain string `gorm:"column:responsibility_domain;type:varchar(32);not null;default:unknown;index:ix_follower_users_responsibility_domain;check:ck_follower_users_responsibility_domain,char_length(btrim(responsibility_domain)) > 0"`
	// CreatedAt 表示记录创建时间。
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime"`
	// CreatedBy 表示记录创建人，当前默认使用 system。
	CreatedBy string `gorm:"column:created_by;type:varchar(64);not null;default:system"`
	// UpdatedAt 表示记录最近更新时间。
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoCreateTime;autoUpdateTime;index:ix_follower_users_updated_at,sort:desc"`
	// UpdatedBy 表示记录最近更新人，当前默认使用 system。
	UpdatedBy string `gorm:"column:updated_by;type:varchar(64);not null;default:system"`
	// IsDeleted 表示逻辑删除标记。
	IsDeleted bool `gorm:"column:is_deleted;type:boolean;not null;default:false"`
	// DeletedAt 表示逻辑删除发生时间，未删除时为空。
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamptz"`
	// DeletedBy 表示逻辑删除执行人，未删除时为空。
	DeletedBy *string `gorm:"column:deleted_by;type:varchar(64)"`
	// RowVersion 表示乐观锁版本号，更新和删除时需要校验。
	RowVersion int `gorm:"column:row_version;type:integer;not null;default:1;check:ck_follower_users_row_version,row_version > 0"`
}

// TableName 负责声明 Gorm 需要操作的真实表名。
func (FollowerUser) TableName() string {
	return "operator_portal.follower_users"
}
