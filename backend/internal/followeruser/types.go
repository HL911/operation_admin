package followeruser

import "time"

const (
	// AccountStatusActive 表示账户处于启用状态。
	AccountStatusActive = "active"
	// AccountStatusDisabled 表示账户处于停用状态。
	AccountStatusDisabled = "disabled"
	// StrategyStatusEnabled 表示策略处于启用状态。
	StrategyStatusEnabled = "enabled"
	// StrategyStatusDisabled 表示策略处于停用状态。
	StrategyStatusDisabled = "disabled"
	// BindingStatusPending 表示绑定流程处于待处理状态。
	BindingStatusPending = "pending"
	// BindingStatusBound 表示当前用户已经完成绑定。
	BindingStatusBound = "bound"
	// BindingStatusUnbound 表示当前用户尚未绑定。
	BindingStatusUnbound = "unbound"
	// DefaultOperator 表示当前没有接入真实登录用户时使用的默认操作者。
	DefaultOperator = "system"
	// DefaultPageNum 表示列表接口的默认页码。
	DefaultPageNum = 1
	// DefaultPageSize 表示列表接口的默认分页大小。
	DefaultPageSize = 20
	// MaxPageSize 表示列表接口允许的最大分页大小，用于防止单次查询过大。
	MaxPageSize = 200
)

// CreateCommand 描述创建小龙虾用户时服务层需要的输入。
type CreateCommand struct {
	// UserID 表示要创建的用户 ID。
	UserID string
	// AccountStatus 表示要写入的账户状态。
	AccountStatus string
	// StrategyStatus 表示要写入的策略状态。
	StrategyStatus string
	// BindingStatus 表示要写入的绑定状态。
	BindingStatus string
	// ResponsibilityDomain 表示要写入的责任域。
	ResponsibilityDomain string
}

// ListQuery 描述列表查询时服务层需要的筛选条件。
type ListQuery struct {
	// PageNum 表示请求的页码。
	PageNum int
	// PageSize 表示请求的分页大小。
	PageSize int
	// UserID 表示按 user_id 精确筛选。
	UserID string
	// AccountStatus 表示按账户状态筛选。
	AccountStatus string
	// StrategyStatus 表示按策略状态筛选。
	StrategyStatus string
	// BindingStatus 表示按绑定状态筛选。
	BindingStatus string
	// ResponsibilityDomain 表示按责任域筛选。
	ResponsibilityDomain string
	// UpdatedFrom 表示更新时间范围的开始值，格式为 yyyy-MM-dd HH:mm:ss。
	UpdatedFrom string
	// UpdatedTo 表示更新时间范围的结束值，格式为 yyyy-MM-dd HH:mm:ss。
	UpdatedTo string
}

// UpdateCommand 描述更新小龙虾用户时服务层需要的输入。
type UpdateCommand struct {
	// UserID 表示要更新的目标用户 ID。
	UserID string
	// AccountStatus 表示新的账户状态。
	AccountStatus string
	// StrategyStatus 表示新的策略状态。
	StrategyStatus string
	// BindingStatus 表示新的绑定状态。
	BindingStatus string
	// ResponsibilityDomain 表示新的责任域。
	ResponsibilityDomain string
	// RowVersion 表示调用方提交的当前版本号。
	RowVersion int
}

// DeleteCommand 描述逻辑删除小龙虾用户时服务层需要的输入。
type DeleteCommand struct {
	// UserID 表示要删除的目标用户 ID。
	UserID string
	// RowVersion 表示调用方提交的当前版本号。
	RowVersion int
}

// ListFilter 描述仓储层实际使用的分页与筛选条件。
type ListFilter struct {
	// Offset 表示分页查询的起始偏移量。
	Offset int
	// Limit 表示分页查询的返回条数上限。
	Limit int
	// UserID 表示按 user_id 精确筛选。
	UserID string
	// AccountStatus 表示按账户状态筛选。
	AccountStatus string
	// StrategyStatus 表示按策略状态筛选。
	StrategyStatus string
	// BindingStatus 表示按绑定状态筛选。
	BindingStatus string
	// ResponsibilityDomain 表示按责任域筛选。
	ResponsibilityDomain string
	// UpdatedFrom 表示更新时间范围的开始值。
	UpdatedFrom *time.Time
	// UpdatedTo 表示更新时间范围的结束值。
	UpdatedTo *time.Time
}

// FollowerUserView 描述接口统一返回的小龙虾用户对象。
type FollowerUserView struct {
	// UserID 表示用户 ID。
	UserID string `json:"userId"`
	// AccountStatus 表示账户状态。
	AccountStatus string `json:"accountStatus"`
	// StrategyStatus 表示策略状态。
	StrategyStatus string `json:"strategyStatus"`
	// BindingStatus 表示绑定状态。
	BindingStatus string `json:"bindingStatus"`
	// ResponsibilityDomain 表示责任域。
	ResponsibilityDomain string `json:"responsibilityDomain"`
	// UpdatedAt 表示最近更新时间，格式为 yyyy-MM-dd HH:mm:ss。
	UpdatedAt string `json:"updatedAt"`
}

// DeleteResult 描述逻辑删除成功后的返回结构。
type DeleteResult struct {
	// UserID 表示被删除的目标用户 ID。
	UserID string `json:"userId"`
	// Success 表示本次删除动作是否成功执行。
	Success bool `json:"success"`
}
