package followeruser

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"operation_admin/backend/internal/http/response"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Service 描述小龙虾用户领域层对外暴露的业务能力。
type Service interface {
	// Create 负责创建小龙虾用户。
	Create(ctx context.Context, command CreateCommand) (*FollowerUserView, error)
	// List 负责分页查询小龙虾用户。
	List(ctx context.Context, query ListQuery) ([]FollowerUserView, int64, int, int, error)
	// Get 负责查询单个小龙虾用户详情。
	Get(ctx context.Context, userID string) (*FollowerUserView, error)
	// Update 负责更新小龙虾用户信息。
	Update(ctx context.Context, command UpdateCommand) (*FollowerUserView, error)
	// Delete 负责逻辑删除小龙虾用户。
	Delete(ctx context.Context, command DeleteCommand) (*DeleteResult, error)
}

// FollowerUserService 描述小龙虾用户领域服务实现。
type FollowerUserService struct {
	// Repository 提供当前服务需要的数据库访问能力。
	Repository Repository
}

// NewService 负责创建小龙虾用户领域服务。
func NewService(repository Repository) *FollowerUserService {
	return &FollowerUserService{
		Repository: repository,
	}
}

// Create 负责创建小龙虾用户。
func (s *FollowerUserService) Create(ctx context.Context, command CreateCommand) (*FollowerUserView, error) {
	// normalizedUserID 保存去除空白后的 userId。
	normalizedUserID := strings.TrimSpace(command.UserID)
	if err := validateUserID(normalizedUserID); err != nil {
		return nil, err
	}

	// normalizedAccountStatus 保存去除空白后的账户状态。
	normalizedAccountStatus := strings.TrimSpace(command.AccountStatus)
	if err := validateAccountStatus(normalizedAccountStatus); err != nil {
		return nil, err
	}

	// normalizedStrategyStatus 保存去除空白后的策略状态。
	normalizedStrategyStatus := strings.TrimSpace(command.StrategyStatus)
	if err := validateStrategyStatus(normalizedStrategyStatus); err != nil {
		return nil, err
	}

	// normalizedBindingStatus 保存去除空白后的绑定状态。
	normalizedBindingStatus := strings.TrimSpace(command.BindingStatus)
	if err := validateBindingStatus(normalizedBindingStatus); err != nil {
		return nil, err
	}

	// normalizedResponsibilityDomain 保存去除空白后的责任域。
	normalizedResponsibilityDomain := strings.TrimSpace(command.ResponsibilityDomain)
	if err := validateResponsibilityDomain(normalizedResponsibilityDomain); err != nil {
		return nil, err
	}

	// followerUser 保存即将写入数据库的新记录。
	followerUser := &FollowerUser{
		UserID:               normalizedUserID,
		AccountStatus:        normalizedAccountStatus,
		StrategyStatus:       normalizedStrategyStatus,
		BindingStatus:        normalizedBindingStatus,
		ResponsibilityDomain: normalizedResponsibilityDomain,
		CreatedBy:            DefaultOperator,
		UpdatedBy:            DefaultOperator,
		RowVersion:           1,
	}

	if err := s.Repository.Create(ctx, followerUser); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrFollowerUserDuplicated
		}

		return nil, err
	}

	// followerUserView 保存创建成功后需要返回给前端的对象。
	followerUserView := buildFollowerUserView(*followerUser)
	return &followerUserView, nil
}

// List 负责分页查询小龙虾用户。
func (s *FollowerUserService) List(ctx context.Context, query ListQuery) ([]FollowerUserView, int64, int, int, error) {
	// normalizedPageNum 保存完成默认值处理后的页码。
	normalizedPageNum := normalizePageNum(query.PageNum)
	// normalizedPageSize 保存完成默认值和上限处理后的分页大小。
	normalizedPageSize := normalizePageSize(query.PageSize)

	// normalizedUserID 保存去除空白后的 userId 条件。
	normalizedUserID := strings.TrimSpace(query.UserID)
	if normalizedUserID != "" {
		if err := validateUserID(normalizedUserID); err != nil {
			return nil, 0, 0, 0, err
		}
	}

	// normalizedAccountStatus 保存去除空白后的账户状态条件。
	normalizedAccountStatus := strings.TrimSpace(query.AccountStatus)
	if normalizedAccountStatus != "" {
		if err := validateAccountStatus(normalizedAccountStatus); err != nil {
			return nil, 0, 0, 0, err
		}
	}

	// normalizedStrategyStatus 保存去除空白后的策略状态条件。
	normalizedStrategyStatus := strings.TrimSpace(query.StrategyStatus)
	if normalizedStrategyStatus != "" {
		if err := validateStrategyStatus(normalizedStrategyStatus); err != nil {
			return nil, 0, 0, 0, err
		}
	}

	// normalizedBindingStatus 保存去除空白后的绑定状态条件。
	normalizedBindingStatus := strings.TrimSpace(query.BindingStatus)
	if normalizedBindingStatus != "" {
		if err := validateBindingStatus(normalizedBindingStatus); err != nil {
			return nil, 0, 0, 0, err
		}
	}

	// normalizedResponsibilityDomain 保存去除空白后的责任域条件。
	normalizedResponsibilityDomain := strings.TrimSpace(query.ResponsibilityDomain)
	if normalizedResponsibilityDomain != "" {
		if err := validateResponsibilityDomain(normalizedResponsibilityDomain); err != nil {
			return nil, 0, 0, 0, err
		}
	}

	// updatedFrom 保存解析后的更新时间开始条件。
	updatedFrom, err := parseOptionalTime(query.UpdatedFrom, "updatedFrom")
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// updatedTo 保存解析后的更新时间结束条件。
	updatedTo, err := parseOptionalTime(query.UpdatedTo, "updatedTo")
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if updatedFrom != nil && updatedTo != nil && updatedFrom.After(*updatedTo) {
		return nil, 0, 0, 0, NewValidationError("updatedFrom 不能晚于 updatedTo")
	}

	// filter 保存仓储层实际使用的分页与筛选条件。
	filter := ListFilter{
		Offset:               (normalizedPageNum - 1) * normalizedPageSize,
		Limit:                normalizedPageSize,
		UserID:               normalizedUserID,
		AccountStatus:        normalizedAccountStatus,
		StrategyStatus:       normalizedStrategyStatus,
		BindingStatus:        normalizedBindingStatus,
		ResponsibilityDomain: normalizedResponsibilityDomain,
		UpdatedFrom:          updatedFrom,
		UpdatedTo:            updatedTo,
	}

	// followerUsers 保存仓储层查询到的模型列表。
	followerUsers, total, err := s.Repository.List(ctx, filter)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// followerUserViews 保存转换后的接口返回对象列表。
	followerUserViews := make([]FollowerUserView, 0, len(followerUsers))
	for _, followerUser := range followerUsers {
		// followerUserView 保存当前单条记录对应的返回对象。
		followerUserView := buildFollowerUserView(followerUser)
		followerUserViews = append(followerUserViews, followerUserView)
	}

	return followerUserViews, total, normalizedPageNum, normalizedPageSize, nil
}

// Get 负责查询单个小龙虾用户详情。
func (s *FollowerUserService) Get(ctx context.Context, userID string) (*FollowerUserView, error) {
	// normalizedUserID 保存去除空白后的目标 userId。
	normalizedUserID := strings.TrimSpace(userID)
	if err := validateUserID(normalizedUserID); err != nil {
		return nil, err
	}

	// followerUser 保存数据库中查到的当前用户记录。
	followerUser, err := s.Repository.GetByUserID(ctx, normalizedUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowerUserNotFound
		}

		return nil, err
	}

	// followerUserView 保存转换后的详情返回对象。
	followerUserView := buildFollowerUserView(*followerUser)
	return &followerUserView, nil
}

// Update 负责更新小龙虾用户信息。
func (s *FollowerUserService) Update(ctx context.Context, command UpdateCommand) (*FollowerUserView, error) {
	// normalizedUserID 保存去除空白后的目标 userId。
	normalizedUserID := strings.TrimSpace(command.UserID)
	if err := validateUserID(normalizedUserID); err != nil {
		return nil, err
	}

	// normalizedAccountStatus 保存去除空白后的账户状态。
	normalizedAccountStatus := strings.TrimSpace(command.AccountStatus)
	if err := validateAccountStatus(normalizedAccountStatus); err != nil {
		return nil, err
	}

	// normalizedStrategyStatus 保存去除空白后的策略状态。
	normalizedStrategyStatus := strings.TrimSpace(command.StrategyStatus)
	if err := validateStrategyStatus(normalizedStrategyStatus); err != nil {
		return nil, err
	}

	// normalizedBindingStatus 保存去除空白后的绑定状态。
	normalizedBindingStatus := strings.TrimSpace(command.BindingStatus)
	if err := validateBindingStatus(normalizedBindingStatus); err != nil {
		return nil, err
	}

	// normalizedResponsibilityDomain 保存去除空白后的责任域。
	normalizedResponsibilityDomain := strings.TrimSpace(command.ResponsibilityDomain)
	if err := validateResponsibilityDomain(normalizedResponsibilityDomain); err != nil {
		return nil, err
	}

	if command.RowVersion <= 0 {
		return nil, NewValidationError("rowVersion 必须大于 0")
	}

	// currentFollowerUser 保存数据库中当前版本的用户记录。
	currentFollowerUser, err := s.Repository.GetByUserID(ctx, normalizedUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowerUserNotFound
		}

		return nil, err
	}

	if currentFollowerUser.RowVersion != command.RowVersion {
		return nil, ErrFollowerUserVersionConflict
	}

	// now 保存本次更新动作统一使用的当前时间。
	now := time.Now()
	// updates 保存本次更新需要写回数据库的列值。
	updates := map[string]any{
		"account_status":        normalizedAccountStatus,
		"strategy_status":       normalizedStrategyStatus,
		"binding_status":        normalizedBindingStatus,
		"responsibility_domain": normalizedResponsibilityDomain,
		"updated_at":            now,
		"updated_by":            DefaultOperator,
		"row_version":           gorm.Expr("row_version + 1"),
	}

	// updated 保存本次更新是否命中了目标记录。
	updated, err := s.Repository.Update(ctx, normalizedUserID, command.RowVersion, updates)
	if err != nil {
		return nil, err
	}

	if !updated {
		return nil, ErrFollowerUserVersionConflict
	}

	// updatedFollowerUser 保存更新完成后重新查询到的最新记录。
	updatedFollowerUser, err := s.Repository.GetByUserID(ctx, normalizedUserID)
	if err != nil {
		return nil, err
	}

	// followerUserView 保存转换后的更新结果对象。
	followerUserView := buildFollowerUserView(*updatedFollowerUser)
	return &followerUserView, nil
}

// Delete 负责逻辑删除小龙虾用户。
func (s *FollowerUserService) Delete(ctx context.Context, command DeleteCommand) (*DeleteResult, error) {
	// normalizedUserID 保存去除空白后的目标 userId。
	normalizedUserID := strings.TrimSpace(command.UserID)
	if err := validateUserID(normalizedUserID); err != nil {
		return nil, err
	}

	if command.RowVersion <= 0 {
		return nil, NewValidationError("rowVersion 必须大于 0")
	}

	// currentFollowerUser 保存数据库中当前版本的用户记录。
	currentFollowerUser, err := s.Repository.GetByUserID(ctx, normalizedUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFollowerUserNotFound
		}

		return nil, err
	}

	if currentFollowerUser.RowVersion != command.RowVersion {
		return nil, ErrFollowerUserVersionConflict
	}

	// deletedAt 保存本次逻辑删除动作使用的统一时间。
	deletedAt := time.Now()
	// deleted 保存本次逻辑删除是否命中了目标记录。
	deleted, err := s.Repository.SoftDelete(ctx, normalizedUserID, command.RowVersion, deletedAt, DefaultOperator)
	if err != nil {
		return nil, err
	}

	if !deleted {
		return nil, ErrFollowerUserVersionConflict
	}

	// deleteResult 保存逻辑删除成功后的返回结果。
	deleteResult := &DeleteResult{
		UserID:  normalizedUserID,
		Success: true,
	}

	return deleteResult, nil
}

// buildFollowerUserView 负责把数据库模型转换为接口返回对象。
func buildFollowerUserView(followerUser FollowerUser) FollowerUserView {
	return FollowerUserView{
		UserID:               followerUser.UserID,
		AccountStatus:        followerUser.AccountStatus,
		StrategyStatus:       followerUser.StrategyStatus,
		BindingStatus:        followerUser.BindingStatus,
		ResponsibilityDomain: followerUser.ResponsibilityDomain,
		UpdatedAt:            followerUser.UpdatedAt.Format(response.DateTimeLayout),
	}
}

// validateUserID 负责校验 userId 是否符合当前接口要求。
func validateUserID(userID string) error {
	if userID == "" {
		return NewValidationError("userId 不能为空")
	}

	if len(userID) > 64 {
		return NewValidationError("userId 长度不能超过 64")
	}

	return nil
}

// validateAccountStatus 负责校验账户状态是否在允许范围内。
func validateAccountStatus(accountStatus string) error {
	switch accountStatus {
	case AccountStatusActive, AccountStatusDisabled:
		return nil
	default:
		return NewValidationError("accountStatus 只能是 active 或 disabled")
	}
}

// validateStrategyStatus 负责校验策略状态是否在允许范围内。
func validateStrategyStatus(strategyStatus string) error {
	switch strategyStatus {
	case StrategyStatusEnabled, StrategyStatusDisabled:
		return nil
	default:
		return NewValidationError("strategyStatus 只能是 enabled 或 disabled")
	}
}

// validateBindingStatus 负责校验绑定状态是否在允许范围内。
func validateBindingStatus(bindingStatus string) error {
	switch bindingStatus {
	case BindingStatusPending, BindingStatusBound, BindingStatusUnbound:
		return nil
	default:
		return NewValidationError("bindingStatus 只能是 pending、bound 或 unbound")
	}
}

// validateResponsibilityDomain 负责校验责任域字段是否为空或超长。
func validateResponsibilityDomain(responsibilityDomain string) error {
	if responsibilityDomain == "" {
		return NewValidationError("responsibilityDomain 不能为空")
	}

	if len(responsibilityDomain) > 32 {
		return NewValidationError("responsibilityDomain 长度不能超过 32")
	}

	return nil
}

// normalizePageNum 负责为页码补齐默认值。
func normalizePageNum(pageNum int) int {
	if pageNum <= 0 {
		return DefaultPageNum
	}

	return pageNum
}

// normalizePageSize 负责为分页大小补齐默认值并施加上限。
func normalizePageSize(pageSize int) int {
	if pageSize <= 0 {
		return DefaultPageSize
	}

	if pageSize > MaxPageSize {
		return MaxPageSize
	}

	return pageSize
}

// parseOptionalTime 负责解析可选时间字符串筛选条件。
func parseOptionalTime(rawValue string, fieldName string) (*time.Time, error) {
	// normalizedValue 保存去除首尾空白后的时间字符串。
	normalizedValue := strings.TrimSpace(rawValue)
	if normalizedValue == "" {
		return nil, nil
	}

	// parsedTime 保存按统一格式解析得到的时间值。
	parsedTime, err := time.ParseInLocation(response.DateTimeLayout, normalizedValue, time.Local)
	if err != nil {
		return nil, NewValidationError(fmt.Sprintf("%s 时间格式必须为 yyyy-MM-dd HH:mm:ss", fieldName))
	}

	return &parsedTime, nil
}

// isUniqueViolation 负责识别数据库唯一键冲突错误。
func isUniqueViolation(err error) bool {
	// pgError 保存从 error 链中提取到的 PostgreSQL 错误对象。
	var pgError *pgconn.PgError
	if !errors.As(err, &pgError) {
		return false
	}

	return pgError.Code == "23505"
}
