package followeruser

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Repository 描述小龙虾用户仓储层需要提供的数据访问能力。
type Repository interface {
	// EnsureSchema 负责确保业务 schema 和 follower_users 表已经就绪。
	EnsureSchema(ctx context.Context) error
	// Create 负责把新用户记录写入数据库。
	Create(ctx context.Context, followerUser *FollowerUser) error
	// GetByUserID 负责按 user_id 查询单条未删除记录。
	GetByUserID(ctx context.Context, userID string) (*FollowerUser, error)
	// List 负责按筛选条件分页查询未删除记录。
	List(ctx context.Context, filter ListFilter) ([]FollowerUser, int64, error)
	// Update 负责基于 user_id 和 row_version 更新记录。
	Update(ctx context.Context, userID string, rowVersion int, updates map[string]any) (bool, error)
	// SoftDelete 负责基于 user_id 和 row_version 执行逻辑删除。
	SoftDelete(ctx context.Context, userID string, rowVersion int, deletedAt time.Time, deletedBy string) (bool, error)
}

// GormRepository 描述基于 Gorm 的小龙虾用户仓储实现。
type GormRepository struct {
	// DB 持有当前服务使用的 Gorm 数据库连接实例。
	DB *gorm.DB
}

// NewRepository 负责创建基于 Gorm 的小龙虾用户仓储。
func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		DB: db,
	}
}

// EnsureSchema 负责确保业务 schema 和 follower_users 表已经就绪。
func (r *GormRepository) EnsureSchema(ctx context.Context) error {
	// schemaSQL 用于确保 operator_portal schema 已存在。
	schemaSQL := "CREATE SCHEMA IF NOT EXISTS operator_portal"
	if err := r.DB.WithContext(ctx).Exec(schemaSQL).Error; err != nil {
		return fmt.Errorf("创建 operator_portal schema 失败: %w", err)
	}

	if err := r.DB.WithContext(ctx).AutoMigrate(&FollowerUser{}); err != nil {
		return fmt.Errorf("迁移 follower_users 表失败: %w", err)
	}

	return nil
}

// Create 负责把新用户记录写入数据库。
func (r *GormRepository) Create(ctx context.Context, followerUser *FollowerUser) error {
	if err := r.DB.WithContext(ctx).Create(followerUser).Error; err != nil {
		return fmt.Errorf("创建小龙虾用户失败: %w", err)
	}

	return nil
}

// GetByUserID 负责按 user_id 查询单条未删除记录。
func (r *GormRepository) GetByUserID(ctx context.Context, userID string) (*FollowerUser, error) {
	// followerUser 保存查询到的小龙虾用户记录。
	var followerUser FollowerUser

	if err := r.DB.WithContext(ctx).
		Where("user_id = ? AND is_deleted = ?", userID, false).
		First(&followerUser).Error; err != nil {
		return nil, err
	}

	return &followerUser, nil
}

// List 负责按筛选条件分页查询未删除记录。
func (r *GormRepository) List(ctx context.Context, filter ListFilter) ([]FollowerUser, int64, error) {
	// query 用于逐步叠加查询条件并最终执行总数与列表查询。
	query := r.DB.WithContext(ctx).Model(&FollowerUser{}).Where("is_deleted = ?", false)
	query = applyListFilters(query, filter)

	// total 保存满足当前筛选条件的总记录数。
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计小龙虾用户总数失败: %w", err)
	}

	// followerUsers 保存当前页实际返回的小龙虾用户列表。
	var followerUsers []FollowerUser
	if err := query.
		Order("updated_at DESC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&followerUsers).Error; err != nil {
		return nil, 0, fmt.Errorf("查询小龙虾用户列表失败: %w", err)
	}

	return followerUsers, total, nil
}

// Update 负责基于 user_id 和 row_version 更新记录。
func (r *GormRepository) Update(ctx context.Context, userID string, rowVersion int, updates map[string]any) (bool, error) {
	// result 保存本次更新语句的执行结果，用于判断是否命中记录。
	result := r.DB.WithContext(ctx).
		Model(&FollowerUser{}).
		Where("user_id = ? AND is_deleted = ? AND row_version = ?", userID, false, rowVersion).
		Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("更新小龙虾用户失败: %w", result.Error)
	}

	return result.RowsAffected > 0, nil
}

// SoftDelete 负责基于 user_id 和 row_version 执行逻辑删除。
func (r *GormRepository) SoftDelete(ctx context.Context, userID string, rowVersion int, deletedAt time.Time, deletedBy string) (bool, error) {
	// updates 保存逻辑删除时需要写入的数据列。
	updates := map[string]any{
		"is_deleted":  true,
		"deleted_at":  deletedAt,
		"deleted_by":  deletedBy,
		"updated_at":  deletedAt,
		"updated_by":  deletedBy,
		"row_version": gorm.Expr("row_version + 1"),
	}

	// result 保存本次逻辑删除语句的执行结果，用于判断是否命中记录。
	result := r.DB.WithContext(ctx).
		Model(&FollowerUser{}).
		Where("user_id = ? AND is_deleted = ? AND row_version = ?", userID, false, rowVersion).
		Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("删除小龙虾用户失败: %w", result.Error)
	}

	return result.RowsAffected > 0, nil
}

// applyListFilters 负责把列表查询筛选条件追加到 Gorm 查询对象上。
func applyListFilters(query *gorm.DB, filter ListFilter) *gorm.DB {
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.AccountStatus != "" {
		query = query.Where("account_status = ?", filter.AccountStatus)
	}

	if filter.StrategyStatus != "" {
		query = query.Where("strategy_status = ?", filter.StrategyStatus)
	}

	if filter.BindingStatus != "" {
		query = query.Where("binding_status = ?", filter.BindingStatus)
	}

	if filter.ResponsibilityDomain != "" {
		query = query.Where("responsibility_domain = ?", filter.ResponsibilityDomain)
	}

	if filter.UpdatedFrom != nil {
		query = query.Where("updated_at >= ?", *filter.UpdatedFrom)
	}

	if filter.UpdatedTo != nil {
		query = query.Where("updated_at <= ?", *filter.UpdatedTo)
	}

	return query
}
