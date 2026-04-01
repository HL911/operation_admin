/**
 * FollowerUserAccountStatus 描述小龙虾用户账户状态的可选值。
 */
export type FollowerUserAccountStatus = "active" | "disabled";

/**
 * FollowerUserStrategyStatus 描述小龙虾用户策略状态的可选值。
 */
export type FollowerUserStrategyStatus = "enabled" | "disabled";

/**
 * FollowerUserBindingStatus 描述小龙虾用户绑定状态的可选值。
 */
export type FollowerUserBindingStatus = "pending" | "bound" | "unbound";

/**
 * FollowerUserRecord 描述接口文档中统一返回的小龙虾用户对象。
 */
export interface FollowerUserRecord {
  /** userId 表示运营侧使用的用户标识。 */
  userId: string;
  /** accountStatus 表示账户当前是否启用。 */
  accountStatus: FollowerUserAccountStatus;
  /** strategyStatus 表示策略能力当前是否启用。 */
  strategyStatus: FollowerUserStrategyStatus;
  /** bindingStatus 表示绑定流程所处的阶段。 */
  bindingStatus: FollowerUserBindingStatus;
  /** responsibilityDomain 表示该用户所属的责任域。 */
  responsibilityDomain: string;
  /** updatedAt 表示最近一次更新时间。 */
  updatedAt: string;
  /** rowVersion 表示用于乐观锁的版本号；接口文档未声明返回值时允许缺失。 */
  rowVersion?: number;
}

/**
 * FollowerUserListQuery 描述列表查询支持的筛选条件。
 */
export interface FollowerUserListQuery {
  /** pageNum 表示当前页码，从 1 开始。 */
  pageNum: number;
  /** pageSize 表示每页返回记录数。 */
  pageSize: number;
  /** userId 表示按用户 ID 的精确或前缀筛选值。 */
  userId: string;
  /** accountStatus 表示按账户状态筛选。 */
  accountStatus: "" | FollowerUserAccountStatus;
  /** strategyStatus 表示按策略状态筛选。 */
  strategyStatus: "" | FollowerUserStrategyStatus;
  /** bindingStatus 表示按绑定状态筛选。 */
  bindingStatus: "" | FollowerUserBindingStatus;
  /** responsibilityDomain 表示按责任域筛选。 */
  responsibilityDomain: string;
  /** updatedFrom 表示更新时间筛选起点。 */
  updatedFrom: string;
  /** updatedTo 表示更新时间筛选终点。 */
  updatedTo: string;
}

/**
 * FollowerUserListData 描述列表接口 `data` 字段中的分页结构。
 */
export interface FollowerUserListData {
  /** list 表示当前页的小龙虾用户列表。 */
  list: readonly FollowerUserRecord[];
  /** total 表示命中的总记录数。 */
  total: number;
  /** pageNum 表示当前页码。 */
  pageNum: number;
  /** pageSize 表示当前分页大小。 */
  pageSize: number;
  /** pages 表示总页数。 */
  pages: number;
}

/**
 * FollowerUserCreatePayload 描述创建请求体。
 */
export interface FollowerUserCreatePayload {
  /** userId 表示要创建的用户 ID。 */
  userId: string;
  /** accountStatus 表示账户状态。 */
  accountStatus: FollowerUserAccountStatus;
  /** strategyStatus 表示策略状态。 */
  strategyStatus: FollowerUserStrategyStatus;
  /** bindingStatus 表示绑定状态。 */
  bindingStatus: FollowerUserBindingStatus;
  /** responsibilityDomain 表示责任域。 */
  responsibilityDomain: string;
}

/**
 * FollowerUserUpdatePayload 描述更新请求体。
 */
export interface FollowerUserUpdatePayload {
  /** accountStatus 表示账户状态。 */
  accountStatus: FollowerUserAccountStatus;
  /** strategyStatus 表示策略状态。 */
  strategyStatus: FollowerUserStrategyStatus;
  /** bindingStatus 表示绑定状态。 */
  bindingStatus: FollowerUserBindingStatus;
  /** responsibilityDomain 表示责任域。 */
  responsibilityDomain: string;
  /** rowVersion 表示乐观锁版本号。 */
  rowVersion: number;
}

/**
 * FollowerUserDeletePayload 描述逻辑删除请求体。
 */
export interface FollowerUserDeletePayload {
  /** rowVersion 表示删除时要提交的乐观锁版本号。 */
  rowVersion: number;
}

/**
 * FollowerUserDeleteResult 描述删除成功后的返回数据。
 */
export interface FollowerUserDeleteResult {
  /** userId 表示被删除的小龙虾用户 ID。 */
  userId: string;
  /** success 表示删除动作是否成功。 */
  success: boolean;
}

/**
 * FollowerUserFormState 描述创建或编辑面板中的表单状态。
 */
export interface FollowerUserFormState {
  /** userId 表示表单当前操作的用户 ID。 */
  userId: string;
  /** accountStatus 表示当前账户状态。 */
  accountStatus: FollowerUserAccountStatus;
  /** strategyStatus 表示当前策略状态。 */
  strategyStatus: FollowerUserStrategyStatus;
  /** bindingStatus 表示当前绑定状态。 */
  bindingStatus: FollowerUserBindingStatus;
  /** responsibilityDomain 表示当前责任域输入。 */
  responsibilityDomain: string;
  /** rowVersion 表示当前编辑或删除使用的版本号。 */
  rowVersion: number;
}

/**
 * FollowerUserApiEnvelope 描述通用接口包装结构。
 */
export interface FollowerUserApiEnvelope<TData> {
  /** code 表示接口返回码。 */
  code: number | string;
  /** message 表示接口返回描述。 */
  message: string;
  /** data 表示当前请求真正关心的数据体。 */
  data: TData;
  /** timestamp 表示接口返回时间戳。 */
  timestamp: number;
}
