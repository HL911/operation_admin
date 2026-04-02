package adminauth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"operation_admin/backend/internal/http/response"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	// loginNamePattern 描述后台登录名允许的字符范围和长度约束。
	loginNamePattern = regexp.MustCompile(`^[a-z0-9._-]{4,32}$`)
)

// Service 描述后台鉴权模块对外暴露的业务能力。
type Service interface {
	// Login 负责校验后台账号凭证并返回新的登录会话。
	Login(ctx context.Context, command LoginCommand) (*LoginResult, error)
	// Refresh 负责基于有效 Refresh Token 轮换登录会话。
	Refresh(ctx context.Context, command RefreshCommand) (*RefreshResult, error)
	// Logout 负责撤销当前后台账号名下的指定 Refresh Token。
	Logout(ctx context.Context, command LogoutCommand) (*LogoutResult, error)
	// CreateAdminUser 负责由管理员创建新的后台开发者账号。
	CreateAdminUser(ctx context.Context, operator CurrentAdmin, command CreateAdminUserCommand) (*CreatedAdminUserView, error)
	// AuthenticateAccessToken 负责校验 Access Token 并返回当前后台账号上下文。
	AuthenticateAccessToken(ctx context.Context, accessToken string) (*CurrentAdmin, error)
}

// AdminAuthService 描述后台鉴权模块的领域服务实现。
type AdminAuthService struct {
	// Repository 提供后台鉴权流程所需的持久化能力。
	Repository Repository
	// TokenManager 提供 Access Token 与 Refresh Token 的签发和校验能力。
	TokenManager *TokenManager
}

// NewService 负责创建后台鉴权领域服务。
func NewService(repository Repository, tokenManager *TokenManager) *AdminAuthService {
	return &AdminAuthService{
		Repository:   repository,
		TokenManager: tokenManager,
	}
}

// Login 负责校验后台账号凭证并返回新的登录会话。
func (s *AdminAuthService) Login(ctx context.Context, command LoginCommand) (*LoginResult, error) {
	// normalizedLoginName 保存转小写并去除空白后的登录名。
	normalizedLoginName := normalizeLoginName(command.LoginName)
	if err := validateLoginName(normalizedLoginName); err != nil {
		return nil, err
	}

	if err := validatePassword(command.Password); err != nil {
		return nil, err
	}

	// adminUser 保存根据登录名查到的后台账号记录。
	adminUser, err := s.Repository.GetAdminUserByLoginName(ctx, normalizedLoginName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if auditErr := s.createLoginAudit(ctx, nil, normalizedLoginName, LoginAuditResultFailed, LoginAuditReasonInvalidCredentials, command.ClientIP, command.UserAgent); auditErr != nil {
				return nil, auditErr
			}

			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if adminUser.IsDeleted {
		if auditErr := s.createLoginAudit(ctx, &adminUser.AdminUserID, normalizedLoginName, LoginAuditResultFailed, LoginAuditReasonInvalidCredentials, command.ClientIP, command.UserAgent); auditErr != nil {
			return nil, auditErr
		}

		return nil, ErrInvalidCredentials
	}

	if adminUser.Status != AdminStatusActive {
		if auditErr := s.createLoginAudit(ctx, &adminUser.AdminUserID, normalizedLoginName, LoginAuditResultFailed, LoginAuditReasonAccountDisabled, command.ClientIP, command.UserAgent); auditErr != nil {
			return nil, auditErr
		}

		return nil, ErrAdminUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.PasswordHash), []byte(command.Password)); err != nil {
		if auditErr := s.createLoginAudit(ctx, &adminUser.AdminUserID, normalizedLoginName, LoginAuditResultFailed, LoginAuditReasonInvalidCredentials, command.ClientIP, command.UserAgent); auditErr != nil {
			return nil, auditErr
		}

		return nil, ErrInvalidCredentials
	}

	// accessToken 保存当前登录成功后签发的 Access Token 结果。
	accessToken, err := s.TokenManager.GenerateAccessToken(*adminUser)
	if err != nil {
		return nil, err
	}

	// refreshToken 保存当前登录成功后签发的 Refresh Token 结果。
	refreshToken, err := s.TokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// lastLoginAt 保存本次成功登录动作统一使用的登录时间。
	lastLoginAt := time.Now()

	if err := s.Repository.Transaction(ctx, func(repository Repository) error {
		// refreshTokenRecord 保存当前登录成功后需要落库的 Refresh Token 记录。
		refreshTokenRecord := &AdminRefreshToken{
			TokenJTI:         refreshToken.TokenJTI,
			AdminUserID:      adminUser.AdminUserID,
			RefreshTokenHash: refreshToken.TokenHash,
			ExpiresAt:        refreshToken.ExpiresAt,
			ClientIP:         command.ClientIP,
			UserAgent:        command.UserAgent,
		}

		// auditRecord 保存当前登录成功后需要写入的登录审计记录。
		auditRecord := &AdminLoginAudit{
			AdminUserID: &adminUser.AdminUserID,
			LoginName:   normalizedLoginName,
			Result:      LoginAuditResultSuccess,
			Reason:      "",
			ClientIP:    command.ClientIP,
			UserAgent:   command.UserAgent,
		}

		if err := repository.UpdateAdminUserLoginMetadata(ctx, adminUser.AdminUserID, lastLoginAt, command.ClientIP); err != nil {
			return err
		}

		if err := repository.CreateRefreshToken(ctx, refreshTokenRecord); err != nil {
			return err
		}

		if err := repository.CreateLoginAudit(ctx, auditRecord); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	adminUser.LastLoginAt = &lastLoginAt
	adminUser.LastLoginIP = command.ClientIP

	// loginResult 保存最终返回给接口层的登录结果。
	loginResult := &LoginResult{
		AccessToken:           accessToken.Token,
		AccessTokenExpiresIn:  accessToken.ExpiresInSeconds,
		RefreshToken:          refreshToken.Token,
		RefreshTokenExpiresIn: refreshToken.ExpiresInSeconds,
		User:                  buildAdminProfileView(*adminUser),
	}

	return loginResult, nil
}

// Refresh 负责基于有效 Refresh Token 轮换登录会话。
func (s *AdminAuthService) Refresh(ctx context.Context, command RefreshCommand) (*RefreshResult, error) {
	// normalizedRefreshToken 保存去除首尾空白后的原始 Refresh Token。
	normalizedRefreshToken := strings.TrimSpace(command.RefreshToken)
	if normalizedRefreshToken == "" {
		return nil, NewValidationError("refreshToken 不能为空")
	}

	// refreshTokenHash 保存原始 Refresh Token 对应的哈希值。
	refreshTokenHash := HashToken(normalizedRefreshToken)

	// currentRefreshToken 保存数据库中查到的 Refresh Token 记录。
	currentRefreshToken, err := s.Repository.GetRefreshTokenByHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenInvalid
		}

		return nil, err
	}

	// now 保存本次刷新动作统一使用的当前时间。
	now := time.Now()
	if currentRefreshToken.RevokedAt != nil || !currentRefreshToken.ExpiresAt.After(now) {
		return nil, ErrRefreshTokenInvalid
	}

	// adminUser 保存 Refresh Token 对应的后台账号记录。
	adminUser, err := s.Repository.GetAdminUserByAdminUserID(ctx, currentRefreshToken.AdminUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenInvalid
		}

		return nil, err
	}

	if adminUser.IsDeleted {
		return nil, ErrRefreshTokenInvalid
	}

	if adminUser.Status != AdminStatusActive {
		return nil, ErrAdminUserDisabled
	}

	// newAccessToken 保存本次刷新成功后新签发的 Access Token。
	newAccessToken, err := s.TokenManager.GenerateAccessToken(*adminUser)
	if err != nil {
		return nil, err
	}

	// newRefreshToken 保存本次刷新成功后新签发的 Refresh Token。
	newRefreshToken, err := s.TokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.Repository.Transaction(ctx, func(repository Repository) error {
		// revoked 表示旧 Refresh Token 是否在事务内成功撤销。
		revoked, err := repository.RevokeRefreshTokenByJTI(ctx, currentRefreshToken.AdminUserID, currentRefreshToken.TokenJTI, now)
		if err != nil {
			return err
		}

		if !revoked {
			return ErrRefreshTokenInvalid
		}

		// refreshTokenRecord 保存本次刷新成功后需要新建的 Refresh Token 记录。
		refreshTokenRecord := &AdminRefreshToken{
			TokenJTI:         newRefreshToken.TokenJTI,
			AdminUserID:      adminUser.AdminUserID,
			RefreshTokenHash: newRefreshToken.TokenHash,
			ExpiresAt:        newRefreshToken.ExpiresAt,
			ClientIP:         command.ClientIP,
			UserAgent:        command.UserAgent,
		}

		if err := repository.CreateRefreshToken(ctx, refreshTokenRecord); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// refreshResult 保存最终返回给接口层的刷新结果。
	refreshResult := &RefreshResult{
		AccessToken:           newAccessToken.Token,
		AccessTokenExpiresIn:  newAccessToken.ExpiresInSeconds,
		RefreshToken:          newRefreshToken.Token,
		RefreshTokenExpiresIn: newRefreshToken.ExpiresInSeconds,
	}

	return refreshResult, nil
}

// Logout 负责撤销当前后台账号名下的指定 Refresh Token。
func (s *AdminAuthService) Logout(ctx context.Context, command LogoutCommand) (*LogoutResult, error) {
	// normalizedAdminUserID 保存去除空白后的后台账号 ID。
	normalizedAdminUserID := strings.TrimSpace(command.AdminUserID)
	if normalizedAdminUserID == "" {
		return nil, NewValidationError("adminUserId 不能为空")
	}

	// normalizedRefreshToken 保存去除空白后的原始 Refresh Token。
	normalizedRefreshToken := strings.TrimSpace(command.RefreshToken)
	if normalizedRefreshToken == "" {
		return nil, NewValidationError("refreshToken 不能为空")
	}

	// refreshTokenHash 保存本次登出目标 Refresh Token 的哈希值。
	refreshTokenHash := HashToken(normalizedRefreshToken)
	// revokedAt 保存本次登出动作统一使用的撤销时间。
	revokedAt := time.Now()

	if _, err := s.Repository.RevokeRefreshTokenByHash(ctx, normalizedAdminUserID, refreshTokenHash, revokedAt); err != nil {
		return nil, err
	}

	return &LogoutResult{
		Success: true,
	}, nil
}

// CreateAdminUser 负责由管理员创建新的后台开发者账号。
func (s *AdminAuthService) CreateAdminUser(ctx context.Context, operator CurrentAdmin, command CreateAdminUserCommand) (*CreatedAdminUserView, error) {
	if operator.RoleCode != RoleCodeAdmin {
		return nil, ErrAdminPermissionDenied
	}

	// normalizedLoginName 保存转小写并去除空白后的新登录名。
	normalizedLoginName := normalizeLoginName(command.LoginName)
	if err := validateLoginName(normalizedLoginName); err != nil {
		return nil, err
	}

	if err := validatePassword(command.Password); err != nil {
		return nil, err
	}

	// normalizedDisplayName 保存去除首尾空白后的显示名。
	normalizedDisplayName := strings.TrimSpace(command.DisplayName)
	if err := validateDisplayName(normalizedDisplayName); err != nil {
		return nil, err
	}

	// adminUserID 保存本次创建动作生成的后台账号 ID。
	adminUserID, err := s.Repository.NextAdminUserID(ctx)
	if err != nil {
		return nil, err
	}

	// passwordHashBytes 保存 bcrypt 生成的密码哈希字节序列。
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("生成后台账号密码哈希失败: %w", err)
	}

	// adminUser 保存即将写入数据库的新后台账号记录。
	adminUser := &AdminUser{
		AdminUserID:  adminUserID,
		LoginName:    normalizedLoginName,
		PasswordHash: string(passwordHashBytes),
		DisplayName:  normalizedDisplayName,
		RoleCode:     RoleCodeAdmin,
		Status:       AdminStatusActive,
		IsDeleted:    false,
		RowVersion:   1,
	}

	if err := s.Repository.CreateAdminUser(ctx, adminUser); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrAdminUserDuplicated
		}

		return nil, err
	}

	// createdAdminUserView 保存创建成功后需要返回给接口层的结果对象。
	createdAdminUserView := buildCreatedAdminUserView(*adminUser)
	return &createdAdminUserView, nil
}

// AuthenticateAccessToken 负责校验 Access Token 并返回当前后台账号上下文。
func (s *AdminAuthService) AuthenticateAccessToken(ctx context.Context, accessToken string) (*CurrentAdmin, error) {
	// claims 保存 Access Token 解析成功后得到的业务 Claims。
	claims, err := s.TokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	// adminUser 保存根据 Access Token 主体查询到的后台账号记录。
	adminUser, err := s.Repository.GetAdminUserByAdminUserID(ctx, claims.Subject)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccessTokenInvalid
		}

		return nil, err
	}

	if adminUser.IsDeleted {
		return nil, ErrAccessTokenInvalid
	}

	if adminUser.Status != AdminStatusActive {
		return nil, ErrAdminUserDisabled
	}

	if adminUser.LoginName != claims.LoginName || adminUser.RoleCode != claims.RoleCode {
		return nil, ErrAccessTokenInvalid
	}

	// currentAdmin 保存最终需要注入请求上下文的后台账号信息。
	currentAdmin := buildCurrentAdmin(*adminUser)
	return &currentAdmin, nil
}

// createLoginAudit 负责创建一条后台登录审计记录。
func (s *AdminAuthService) createLoginAudit(ctx context.Context, adminUserID *string, loginName string, result string, reason string, clientIP string, userAgent string) error {
	// auditRecord 保存当前待写入数据库的登录审计记录。
	auditRecord := &AdminLoginAudit{
		AdminUserID: adminUserID,
		LoginName:   loginName,
		Result:      result,
		Reason:      reason,
		ClientIP:    clientIP,
		UserAgent:   userAgent,
	}

	return s.Repository.CreateLoginAudit(ctx, auditRecord)
}

// buildAdminProfileView 负责把后台账号模型转换为登录和 me 接口使用的返回对象。
func buildAdminProfileView(adminUser AdminUser) AdminProfileView {
	return AdminProfileView{
		AdminUserID: adminUser.AdminUserID,
		LoginName:   adminUser.LoginName,
		DisplayName: adminUser.DisplayName,
		RoleCode:    adminUser.RoleCode,
		Status:      adminUser.Status,
		LastLoginAt: formatOptionalTime(adminUser.LastLoginAt),
	}
}

// buildCreatedAdminUserView 负责把后台账号模型转换为创建账号接口使用的返回对象。
func buildCreatedAdminUserView(adminUser AdminUser) CreatedAdminUserView {
	return CreatedAdminUserView{
		AdminUserID: adminUser.AdminUserID,
		LoginName:   adminUser.LoginName,
		DisplayName: adminUser.DisplayName,
		RoleCode:    adminUser.RoleCode,
		Status:      adminUser.Status,
		CreatedAt:   adminUser.CreatedAt.Format(response.DateTimeLayout),
	}
}

// buildCurrentAdmin 负责把后台账号模型转换为中间件上下文使用的当前登录账号对象。
func buildCurrentAdmin(adminUser AdminUser) CurrentAdmin {
	return CurrentAdmin{
		AdminUserID: adminUser.AdminUserID,
		LoginName:   adminUser.LoginName,
		DisplayName: adminUser.DisplayName,
		RoleCode:    adminUser.RoleCode,
		Status:      adminUser.Status,
		LastLoginAt: adminUser.LastLoginAt,
	}
}

// formatOptionalTime 负责把可选时间字段转换为接口返回所需的统一字符串格式。
func formatOptionalTime(value *time.Time) *string {
	if value == nil {
		return nil
	}

	// formattedValue 保存格式化后的时间字符串。
	formattedValue := value.Format(response.DateTimeLayout)
	return &formattedValue
}

// normalizeLoginName 负责把登录名转换为去空白且统一小写的标准形态。
func normalizeLoginName(loginName string) string {
	return strings.ToLower(strings.TrimSpace(loginName))
}

// validateLoginName 负责校验后台登录名是否符合允许的字符范围和长度约束。
func validateLoginName(loginName string) error {
	if loginName == "" {
		return NewValidationError("loginName 不能为空")
	}

	if !loginNamePattern.MatchString(loginName) {
		return NewValidationError("loginName 格式不正确，仅支持 4 到 32 位小写字母、数字、点、中划线、下划线")
	}

	return nil
}

// validatePassword 负责校验后台账号密码是否满足最小长度要求。
func validatePassword(password string) error {
	if len(password) < 8 {
		return NewValidationError("password 长度不能少于 8 位")
	}

	return nil
}

// validateDisplayName 负责校验后台账号显示名是否为空或超长。
func validateDisplayName(displayName string) error {
	if displayName == "" {
		return NewValidationError("displayName 不能为空")
	}

	if utf8.RuneCountInString(displayName) > 64 {
		return NewValidationError("displayName 长度不能超过 64")
	}

	return nil
}

// isUniqueViolation 负责识别数据库唯一键冲突错误。
func isUniqueViolation(err error) bool {
	// pgError 保存从错误链中提取到的 PostgreSQL 错误对象。
	var pgError *pgconn.PgError
	if !errors.As(err, &pgError) {
		return false
	}

	return pgError.Code == "23505"
}
