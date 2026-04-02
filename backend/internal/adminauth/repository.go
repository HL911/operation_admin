package adminauth

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Repository 描述后台鉴权模块所需的数据访问能力。
type Repository interface {
	// EnsureSchema 负责确保 operator_portal schema、序列和后台鉴权相关表结构已经就绪。
	EnsureSchema(ctx context.Context) error
	// NextAdminUserID 负责基于 PostgreSQL 序列生成新的后台账号 ID。
	NextAdminUserID(ctx context.Context) (string, error)
	// CreateAdminUser 负责写入新的后台账号记录。
	CreateAdminUser(ctx context.Context, adminUser *AdminUser) error
	// GetAdminUserByLoginName 负责按登录名查询后台账号。
	GetAdminUserByLoginName(ctx context.Context, loginName string) (*AdminUser, error)
	// GetAdminUserByAdminUserID 负责按后台账号 ID 查询后台账号。
	GetAdminUserByAdminUserID(ctx context.Context, adminUserID string) (*AdminUser, error)
	// UpdateAdminUserLoginMetadata 负责更新后台账号最近登录时间和 IP。
	UpdateAdminUserLoginMetadata(ctx context.Context, adminUserID string, lastLoginAt time.Time, lastLoginIP string) error
	// CreateRefreshToken 负责写入新的 Refresh Token 记录。
	CreateRefreshToken(ctx context.Context, refreshToken *AdminRefreshToken) error
	// GetRefreshTokenByHash 负责按 Refresh Token 哈希查询记录。
	GetRefreshTokenByHash(ctx context.Context, refreshTokenHash string) (*AdminRefreshToken, error)
	// RevokeRefreshTokenByJTI 负责按 JTI 撤销指定后台账号名下的 Refresh Token。
	RevokeRefreshTokenByJTI(ctx context.Context, adminUserID string, tokenJTI string, revokedAt time.Time) (bool, error)
	// RevokeRefreshTokenByHash 负责按哈希撤销指定后台账号名下的 Refresh Token。
	RevokeRefreshTokenByHash(ctx context.Context, adminUserID string, refreshTokenHash string, revokedAt time.Time) (bool, error)
	// CreateLoginAudit 负责写入后台登录审计记录。
	CreateLoginAudit(ctx context.Context, audit *AdminLoginAudit) error
	// Transaction 负责在数据库事务中执行一组后台鉴权数据操作。
	Transaction(ctx context.Context, fn func(repository Repository) error) error
}

// GormRepository 描述基于 Gorm 的后台鉴权仓储实现。
type GormRepository struct {
	// DB 持有当前服务使用的 Gorm 数据库连接实例。
	DB *gorm.DB
}

// NewRepository 负责创建基于 Gorm 的后台鉴权仓储。
func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		DB: db,
	}
}

// EnsureSchema 负责确保 operator_portal schema、序列和后台鉴权相关表结构已经就绪。
func (r *GormRepository) EnsureSchema(ctx context.Context) error {
	// schemaSQL 用于确保 operator_portal schema 已存在。
	schemaSQL := "CREATE SCHEMA IF NOT EXISTS operator_portal"
	if err := r.DB.WithContext(ctx).Exec(schemaSQL).Error; err != nil {
		return fmt.Errorf("创建 operator_portal schema 失败: %w", err)
	}

	// sequenceSQL 用于确保后台账号 ID 所依赖的序列已经创建。
	sequenceSQL := `
CREATE SEQUENCE IF NOT EXISTS operator_portal.admin_user_seq
    START WITH 1001
    INCREMENT BY 1
    MINVALUE 1001
    CACHE 1`
	if err := r.DB.WithContext(ctx).Exec(sequenceSQL).Error; err != nil {
		return fmt.Errorf("创建 admin_user_seq 序列失败: %w", err)
	}

	if err := r.DB.WithContext(ctx).AutoMigrate(&AdminUser{}, &AdminRefreshToken{}, &AdminLoginAudit{}); err != nil {
		return fmt.Errorf("迁移后台鉴权相关表失败: %w", err)
	}

	return nil
}

// NextAdminUserID 负责基于 PostgreSQL 序列生成新的后台账号 ID。
func (r *GormRepository) NextAdminUserID(ctx context.Context) (string, error) {
	// nextValue 保存序列返回的下一个数值。
	var nextValue int64
	if err := r.DB.WithContext(ctx).
		Raw("SELECT nextval('operator_portal.admin_user_seq')").
		Scan(&nextValue).Error; err != nil {
		return "", fmt.Errorf("生成后台账号 ID 失败: %w", err)
	}

	return fmt.Sprintf("ADM-%d", nextValue), nil
}

// CreateAdminUser 负责写入新的后台账号记录。
func (r *GormRepository) CreateAdminUser(ctx context.Context, adminUser *AdminUser) error {
	if err := r.DB.WithContext(ctx).Create(adminUser).Error; err != nil {
		return fmt.Errorf("创建后台账号失败: %w", err)
	}

	return nil
}

// GetAdminUserByLoginName 负责按登录名查询后台账号。
func (r *GormRepository) GetAdminUserByLoginName(ctx context.Context, loginName string) (*AdminUser, error) {
	// adminUser 保存查询到的后台账号记录。
	var adminUser AdminUser
	if err := r.DB.WithContext(ctx).
		Where("login_name = ?", loginName).
		First(&adminUser).Error; err != nil {
		return nil, err
	}

	return &adminUser, nil
}

// GetAdminUserByAdminUserID 负责按后台账号 ID 查询后台账号。
func (r *GormRepository) GetAdminUserByAdminUserID(ctx context.Context, adminUserID string) (*AdminUser, error) {
	// adminUser 保存查询到的后台账号记录。
	var adminUser AdminUser
	if err := r.DB.WithContext(ctx).
		Where("admin_user_id = ?", adminUserID).
		First(&adminUser).Error; err != nil {
		return nil, err
	}

	return &adminUser, nil
}

// UpdateAdminUserLoginMetadata 负责更新后台账号最近登录时间和 IP。
func (r *GormRepository) UpdateAdminUserLoginMetadata(ctx context.Context, adminUserID string, lastLoginAt time.Time, lastLoginIP string) error {
	// updates 保存最近登录元数据需要写回数据库的列值。
	updates := map[string]any{
		"last_login_at": lastLoginAt,
		"last_login_ip": lastLoginIP,
		"updated_at":    lastLoginAt,
		"row_version":   gorm.Expr("row_version + 1"),
	}

	if err := r.DB.WithContext(ctx).
		Model(&AdminUser{}).
		Where("admin_user_id = ?", adminUserID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("更新后台账号最近登录信息失败: %w", err)
	}

	return nil
}

// CreateRefreshToken 负责写入新的 Refresh Token 记录。
func (r *GormRepository) CreateRefreshToken(ctx context.Context, refreshToken *AdminRefreshToken) error {
	if err := r.DB.WithContext(ctx).Create(refreshToken).Error; err != nil {
		return fmt.Errorf("创建 refresh token 记录失败: %w", err)
	}

	return nil
}

// GetRefreshTokenByHash 负责按 Refresh Token 哈希查询记录。
func (r *GormRepository) GetRefreshTokenByHash(ctx context.Context, refreshTokenHash string) (*AdminRefreshToken, error) {
	// refreshToken 保存查询到的 Refresh Token 记录。
	var refreshToken AdminRefreshToken
	if err := r.DB.WithContext(ctx).
		Where("refresh_token_hash = ?", refreshTokenHash).
		First(&refreshToken).Error; err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// RevokeRefreshTokenByJTI 负责按 JTI 撤销指定后台账号名下的 Refresh Token。
func (r *GormRepository) RevokeRefreshTokenByJTI(ctx context.Context, adminUserID string, tokenJTI string, revokedAt time.Time) (bool, error) {
	// updates 保存撤销 Refresh Token 时需要写回数据库的列值。
	updates := map[string]any{
		"revoked_at": revokedAt,
		"updated_at": revokedAt,
	}

	// result 保存本次撤销语句的执行结果。
	result := r.DB.WithContext(ctx).
		Model(&AdminRefreshToken{}).
		Where("admin_user_id = ? AND token_jti = ? AND revoked_at IS NULL", adminUserID, tokenJTI).
		Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("按 jti 撤销 refresh token 失败: %w", result.Error)
	}

	return result.RowsAffected > 0, nil
}

// RevokeRefreshTokenByHash 负责按哈希撤销指定后台账号名下的 Refresh Token。
func (r *GormRepository) RevokeRefreshTokenByHash(ctx context.Context, adminUserID string, refreshTokenHash string, revokedAt time.Time) (bool, error) {
	// updates 保存撤销 Refresh Token 时需要写回数据库的列值。
	updates := map[string]any{
		"revoked_at": revokedAt,
		"updated_at": revokedAt,
	}

	// result 保存本次撤销语句的执行结果。
	result := r.DB.WithContext(ctx).
		Model(&AdminRefreshToken{}).
		Where("admin_user_id = ? AND refresh_token_hash = ? AND revoked_at IS NULL", adminUserID, refreshTokenHash).
		Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("按哈希撤销 refresh token 失败: %w", result.Error)
	}

	return result.RowsAffected > 0, nil
}

// CreateLoginAudit 负责写入后台登录审计记录。
func (r *GormRepository) CreateLoginAudit(ctx context.Context, audit *AdminLoginAudit) error {
	if err := r.DB.WithContext(ctx).Create(audit).Error; err != nil {
		return fmt.Errorf("创建后台登录审计失败: %w", err)
	}

	return nil
}

// Transaction 负责在数据库事务中执行一组后台鉴权数据操作。
func (r *GormRepository) Transaction(ctx context.Context, fn func(repository Repository) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// txRepository 表示绑定到当前事务上下文的后台鉴权仓储实现。
		txRepository := &GormRepository{DB: tx}
		return fn(txRepository)
	})
}
