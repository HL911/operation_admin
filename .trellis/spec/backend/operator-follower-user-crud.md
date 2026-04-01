# 小龙虾用户 CRUD 代码规格

> 场景：运营后台管理小龙虾用户。
>
> 本文档是可执行代码规格，面向后端开发实现。统一响应格式、分页结构、错误处理方式必须同时遵循 [错误处理](./error-handling.md)。

---

## Scenario: 小龙虾用户后台 CRUD

### 1. Scope / Trigger

- Trigger：运营后台需要对“小龙虾用户”执行列表查询、详情查询、创建、更新、停用。
- Scope：只覆盖小龙虾用户主档案与基础 Lobster 关联信息。
- Mutable Tables：
  - `public.users`
  - `public.user_lobster_profiles`
  - `public.account_links`
- Read Only Tables：
  - `public.follower_key_records`
  - `public.follower_node_bindings`
  - `public.user_lobster_binding_requests`
  - `public.user_auth_identities`
  - `public.user_sessions`
- 不在本规格内的内容：
  - follower key 的签发、撤销、补发
  - follower 节点绑定关系的直接修改
  - 绑定申请审批流

### 1.1 业务定义

当前库表下，“小龙虾用户”按以下规则识别：

- 满足任一条件即可进入模块检索范围：
  - 在 `public.user_lobster_profiles` 中存在记录
  - 在 `public.account_links` 中存在记录
  - 在 `public.follower_key_records` 中存在记录
  - 在 `public.follower_node_bindings` 中存在记录
  - 在 `public.user_lobster_binding_requests` 中存在记录

当前模块的“主档案”定义为：

- `public.users` 中的用户基础信息
- `public.user_lobster_profiles` 中的 Lobster 配置档案
- `public.account_links` 中的当前有效账户关联

### 1.2 一期写入边界

一期仅允许通过本模块写入或更新以下字段：

- `users.external_user_id`
- `users.display_name`
- `users.status`
- `users.preferred_channel`
- `users.locale`
- `user_lobster_profiles.slot_cap`
- `user_lobster_profiles.tier`
- `user_lobster_profiles.meta_json`
- `account_links.lobster_user_id`
- `account_links.link_status`
- `account_links.meta_json`
- `account_links.last_sync_at`

以下数据只读：

- follower key 历史
- follower 节点绑定历史
- Lobster 绑定申请历史

---

### 2. Signatures

### 2.1 HTTP Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/admin/v1/follower-users` | 查询小龙虾用户分页列表 |
| `GET` | `/admin/v1/follower-users/{userId}` | 查询小龙虾用户详情 |
| `POST` | `/admin/v1/follower-users` | 创建小龙虾用户 |
| `PATCH` | `/admin/v1/follower-users/{userId}` | 更新小龙虾用户 |
| `DELETE` | `/admin/v1/follower-users/{userId}` | 停用小龙虾用户 |

### 2.2 Handler / Service / Repository 建议签名

```go
// handler
func (h *FollowerUserHandler) List(c *gin.Context)
func (h *FollowerUserHandler) Get(c *gin.Context)
func (h *FollowerUserHandler) Create(c *gin.Context)
func (h *FollowerUserHandler) Update(c *gin.Context)
func (h *FollowerUserHandler) Delete(c *gin.Context)

// service
type FollowerUserService interface {
    List(ctx context.Context, req ListFollowerUsersRequest) (PageResult[FollowerUserListItem], error)
    Get(ctx context.Context, userID int64) (*FollowerUserDetail, error)
    Create(ctx context.Context, req CreateFollowerUserRequest) (*FollowerUserDetail, error)
    Update(ctx context.Context, userID int64, req UpdateFollowerUserRequest) (*FollowerUserDetail, error)
    Delete(ctx context.Context, userID int64, req DeleteFollowerUserRequest) error
}
```

### 2.3 Query Params Signature

列表接口统一采用小驼峰命名：

```text
GET /admin/v1/follower-users?pageNum=1&pageSize=20&keyword=张三&status=ACTIVE&tier=BASIC&hasActiveKey=true
```

字段定义：

- `pageNum`：默认 `1`
- `pageSize`：默认 `20`，最大 `100`
- `keyword`：支持匹配 `externalUserId`、`displayName`、`lobsterUserId`
- `status`：用户状态
- `linkStatus`：账户关联状态
- `tier`：Lobster 等级
- `hasActiveKey`：是否有有效 follower key
- `hasActiveBinding`：是否有有效节点绑定
- `createdFrom`
- `createdTo`
- `sortBy`：默认 `updatedAt`
- `sortOrder`：默认 `desc`

---

### 3. Contracts

### 3.1 当前有效账户关联选择规则

`public.account_links` 当前没有 `(portal_user_id)` 唯一约束，因此查询时禁止直接简单 `LEFT JOIN`。

当前有效账户关联的选择规则固定为：

1. 优先取 `link_status = 'LINKED'` 的记录
2. 若存在多条 `LINKED`，取 `updated_at` 最新的一条
3. 若不存在 `LINKED`，取 `updated_at` 最新的一条历史记录
4. 若 `updated_at` 相同，则取 `id` 最大的一条

推荐 SQL：

```sql
WITH ranked_account_links AS (
  SELECT
    al.*,
    ROW_NUMBER() OVER (
      PARTITION BY al.portal_user_id
      ORDER BY
        CASE WHEN al.link_status = 'LINKED' THEN 0 ELSE 1 END,
        al.updated_at DESC,
        al.id DESC
    ) AS rn
  FROM public.account_links al
)
SELECT *
FROM ranked_account_links
WHERE rn = 1;
```

### 3.2 列表查询 Contract

#### 3.2.1 候选用户集合

列表必须以候选用户集合作为入口，不能只靠 `user_lobster_profiles` 或 `account_links`。

推荐候选集：

```sql
WITH candidate_users AS (
  SELECT user_id FROM public.user_lobster_profiles
  UNION
  SELECT portal_user_id AS user_id FROM public.account_links
  UNION
  SELECT user_id FROM public.follower_key_records
  UNION
  SELECT user_id FROM public.follower_node_bindings
  UNION
  SELECT user_id FROM public.user_lobster_binding_requests
)
```

#### 3.2.2 列表返回对象

```json
{
  "userId": 10001,
  "externalUserId": "lobster_u_10001",
  "displayName": "张三",
  "status": "ACTIVE",
  "preferredChannel": "whatsapp",
  "locale": "zh-CN",
  "tier": "BASIC",
  "slotCap": 2,
  "lobsterUserId": "lb_90001",
  "linkStatus": "LINKED",
  "activeKeyCount": 1,
  "activeBindingCount": 1,
  "lastBoundAtTs": 1775104225000,
  "createdAt": "2026-04-01 14:30:00",
  "updatedAt": "2026-04-01 15:00:00"
}
```

#### 3.2.3 列表字段来源

| API 字段 | 来源 |
| --- | --- |
| `userId` | `users.id` |
| `externalUserId` | `users.external_user_id` |
| `displayName` | `users.display_name` |
| `status` | `users.status` |
| `preferredChannel` | `users.preferred_channel` |
| `locale` | `users.locale` |
| `tier` | `user_lobster_profiles.tier` |
| `slotCap` | `user_lobster_profiles.slot_cap` |
| `lobsterUserId` | 当前有效 `account_links.lobster_user_id` |
| `linkStatus` | 当前有效 `account_links.link_status` |
| `activeKeyCount` | `follower_key_records` 聚合 |
| `activeBindingCount` | `follower_node_bindings` 聚合 |
| `lastBoundAtTs` | `follower_node_bindings.bound_at_ts` 最大值 |
| `createdAt` | `users.created_at` |
| `updatedAt` | `users.updated_at` |

#### 3.2.4 有效 key 统计规则

有效 key 默认统计以下状态：

- `ISSUED`
- `BOUND`
- `ACTIVE`

实现时若现网状态集不同，以配置常量维护，不允许散落在 SQL 字符串中。

#### 3.2.5 有效绑定统计规则

有效绑定默认统计以下状态：

- `BOUND`
- `ACTIVE`

### 3.3 详情查询 Contract

详情接口 `data` 必须返回以下结构：

```json
{
  "basic": {
    "userId": 10001,
    "externalUserId": "lobster_u_10001",
    "displayName": "张三",
    "status": "ACTIVE",
    "preferredChannel": "whatsapp",
    "locale": "zh-CN",
    "createdAt": "2026-04-01 14:30:00",
    "updatedAt": "2026-04-01 15:00:00"
  },
  "profile": {
    "profileId": 21,
    "slotCap": 2,
    "tier": "BASIC",
    "meta": {}
  },
  "accountLink": {
    "id": 31,
    "lobsterUserId": "lb_90001",
    "linkStatus": "LINKED",
    "linkedAt": "2026-04-01 14:32:00",
    "lastSyncAt": "2026-04-01 15:10:00",
    "meta": {}
  },
  "keySummary": {
    "total": 2,
    "active": 1,
    "latestExpireAtTs": 1776104225000
  },
  "bindingSummary": {
    "total": 2,
    "active": 1,
    "lastBoundAtTs": 1775104225000
  },
  "keys": [],
  "bindings": [],
  "bindingRequests": [],
  "authIdentities": []
}
```

详情接口默认返回：

- `keys`：最近 `20` 条 follower key，按 `updated_at DESC`
- `bindings`：最近 `20` 条绑定记录，按 `updated_at DESC`
- `bindingRequests`：最近 `20` 条绑定申请，按 `updated_at DESC`
- `authIdentities`：全部身份记录，按 `id DESC`

### 3.4 创建 Contract

#### 3.4.1 Request Body

```json
{
  "externalUserId": "lobster_u_10001",
  "displayName": "张三",
  "status": "ACTIVE",
  "preferredChannel": "whatsapp",
  "locale": "zh-CN",
  "profile": {
    "slotCap": 2,
    "tier": "BASIC",
    "meta": {}
  },
  "accountLink": {
    "lobsterUserId": "lb_90001",
    "linkStatus": "LINKED",
    "lastSyncAt": "2026-04-01 15:10:00",
    "meta": {}
  }
}
```

#### 3.4.2 Request Rules

- `externalUserId`：必填，唯一
- `displayName`：必填，长度 `1..120`
- `status`：可选，默认 `ACTIVE`
- `preferredChannel`：可选
- `locale`：可选
- `profile`：必填
- `profile.slotCap`：必填，`>= 0`
- `profile.tier`：必填
- `accountLink`：可选
- `accountLink.lobsterUserId`：若传入 `accountLink`，则必填
- `accountLink.linkStatus`：可选，默认 `LINKED`

#### 3.4.3 DB Write Sequence

必须在单事务内执行：

1. 插入 `public.users`
2. 插入 `public.user_lobster_profiles`
3. 若有 `accountLink`，插入 `public.account_links`
4. 记录审计日志

#### 3.4.4 Response Contract

- 成功后返回详情接口同结构对象
- `message` 建议固定为 `创建成功`

### 3.5 更新 Contract

#### 3.5.1 Request Body

PATCH 采用部分更新，未传字段表示不变：

```json
{
  "displayName": "张三（新）",
  "status": "ACTIVE",
  "preferredChannel": "whatsapp",
  "profile": {
    "slotCap": 3,
    "tier": "ADVANCED",
    "meta": {
      "operatorRemark": "人工升级"
    }
  },
  "accountLink": {
    "lobsterUserId": "lb_90002",
    "linkStatus": "LINKED",
    "lastSyncAt": "2026-04-01 16:00:00",
    "meta": {}
  }
}
```

#### 3.5.2 PATCH 语义

- 未传字段：不修改
- 传入对象但对象内缺失字段：对象内缺失字段不修改
- 本期不支持通过 `null` 表示清空字段
- 需要解除 Lobster 关联时，使用：
  - `accountLink.linkStatus = "UNLINKED"`

#### 3.5.3 DB Update Sequence

必须在单事务内执行：

1. `SELECT public.users ... FOR UPDATE`
2. 更新 `public.users`
3. `SELECT public.user_lobster_profiles ... FOR UPDATE`
4. 更新 `public.user_lobster_profiles`
5. 若请求中带 `accountLink`：
   - 查找当前有效 `account_links` 记录并 `FOR UPDATE`
   - 找到则更新
   - 未找到则插入一条新记录
6. 记录审计日志

### 3.6 删除 Contract

#### 3.6.1 Request Body

```json
{
  "reason": "人工停用"
}
```

#### 3.6.2 Delete 语义

`DELETE /admin/v1/follower-users/{userId}` 在本模块中表示“逻辑停用”：

1. `users.status` 更新为 `DISABLED`
2. 当前有效 `account_links` 若存在，则更新 `link_status = 'UNLINKED'`
3. 不修改：
   - `follower_key_records`
   - `follower_node_bindings`
   - `user_lobster_binding_requests`

#### 3.6.3 Delete Response

```json
{
  "code": 200,
  "message": "停用成功",
  "data": null,
  "timestamp": 1775104225000
}
```

### 3.7 DB Contract Summary

| Table | Create | Update | Delete | Notes |
| --- | --- | --- | --- | --- |
| `public.users` | insert | update | `status -> DISABLED` | 主表 |
| `public.user_lobster_profiles` | insert | update | no-op | 每个用户只应存在一条配置档案 |
| `public.account_links` | optional insert | update or insert | `link_status -> UNLINKED` | 当前库无 user 维度唯一约束，需按“当前有效账户关联规则”选中目标行 |
| `public.follower_key_records` | no-op | no-op | no-op | 只读 |
| `public.follower_node_bindings` | no-op | no-op | no-op | 只读 |

---

### 4. Validation & Error Matrix

| 场景 | 条件 | HTTP | code | message |
| --- | --- | --- | --- | --- |
| 列表查询 | `pageNum < 1` 或 `pageSize < 1` | `400` | `40001` | 参数校验失败：分页参数非法 |
| 列表查询 | `pageSize > 100` | `400` | `40001` | 参数校验失败：pageSize 超出上限 |
| 创建 | `externalUserId` 为空 | `400` | `40001` | 参数校验失败：externalUserId 不能为空 |
| 创建 | `displayName` 为空 | `400` | `40001` | 参数校验失败：displayName 不能为空 |
| 创建 | `slotCap < 0` | `400` | `40001` | 参数校验失败：slotCap 不能小于 0 |
| 创建 | `externalUserId` 重复 | `400` | `40011` | 外部用户 ID 已存在 |
| 创建 | `accountLink` 存在但 `lobsterUserId` 为空 | `400` | `40001` | 参数校验失败：lobsterUserId 不能为空 |
| 详情/更新/删除 | `userId` 不存在 | `404` | `404` | 用户不存在 |
| 更新 | 用户不存在 Lobster 配置档案 | `400` | `40012` | 小龙虾用户配置档案不存在 |
| 更新 | 状态非法流转 | `400` | `40013` | 状态流转不合法 |
| 更新 | `slotCap < 0` | `400` | `40001` | 参数校验失败：slotCap 不能小于 0 |
| 删除 | 用户已是 `DISABLED` | `200` | `200` | 停用成功 |
| 任意写接口 | 事务内 SQL 失败 | `500` | `500` | 服务器内部错误，请稍后重试 |

### 4.1 状态流转约束

`users.status` 一期允许的流转：

| from | to | allowed |
| --- | --- | --- |
| `ACTIVE` | `DISABLED` | yes |
| `DISABLED` | `ACTIVE` | yes |
| `ACTIVE` | `ACTIVE` | yes |
| `DISABLED` | `DISABLED` | yes |

`account_links.link_status` 一期允许的流转：

| from | to | allowed |
| --- | --- | --- |
| `LINKED` | `UNLINKED` | yes |
| `UNLINKED` | `LINKED` | yes |
| `LINKED` | `LINKED` | yes |
| `UNLINKED` | `UNLINKED` | yes |

若现网采用其他枚举，必须在 service 层维护允许流转表，禁止在 handler 中手写分支。

---

### 5. Good / Base / Bad Cases

### 5.1 Good Case

场景：创建一个完整的小龙虾用户，包含档案和 Lobster 账户关联。

```json
{
  "externalUserId": "lobster_u_10001",
  "displayName": "张三",
  "status": "ACTIVE",
  "preferredChannel": "whatsapp",
  "locale": "zh-CN",
  "profile": {
    "slotCap": 2,
    "tier": "BASIC",
    "meta": {}
  },
  "accountLink": {
    "lobsterUserId": "lb_90001",
    "linkStatus": "LINKED",
    "meta": {}
  }
}
```

预期：

- `users` 插入 1 条
- `user_lobster_profiles` 插入 1 条
- `account_links` 插入 1 条
- 返回完整详情对象

### 5.2 Base Case

场景：创建一个最小化小龙虾用户，只建用户和 Lobster 档案，不建账户关联。

```json
{
  "externalUserId": "lobster_u_10002",
  "displayName": "李四",
  "profile": {
    "slotCap": 0,
    "tier": "BASIC",
    "meta": {}
  }
}
```

预期：

- 可以创建成功
- 详情中的 `accountLink` 返回 `null`
- 列表中 `lobsterUserId` 返回 `null`

### 5.3 Bad Case

场景：重复创建同一个 `externalUserId`。

```json
{
  "externalUserId": "lobster_u_10001",
  "displayName": "重复用户",
  "profile": {
    "slotCap": 1,
    "tier": "BASIC"
  }
}
```

预期：

- 返回 `400`
- `code = 40011`
- 不产生任何新数据

### 5.4 Bad Case

场景：更新时传入负数 `slotCap`。

```json
{
  "profile": {
    "slotCap": -1
  }
}
```

预期：

- 返回 `400`
- `code = 40001`
- 不更新任何表

---

### 6. Tests Required

### 6.1 Handler Tests

- 列表接口参数绑定失败时返回统一错误响应
- 创建接口缺少必填字段时返回 `40001`
- PATCH 接口未传字段时不报错

### 6.2 Service Tests

- 创建成功时必须同时写入 `users` 和 `user_lobster_profiles`
- 创建时若 `account_links` 写入失败，事务必须整体回滚
- 更新时若 `accountLink` 不存在且请求带该对象，必须自动补建
- 删除时必须同时更新 `users.status` 与当前有效 `account_links.link_status`
- 删除已停用用户时应幂等成功

### 6.3 Repository / Query Tests

- 列表查询必须去重，不能因为 `account_links` 多条导致同一用户出现多行
- 列表查询必须能查出“只有 follower key、没有 profile/accountLink”的用户
- 详情查询必须只返回当前有效账户关联，而不是所有 `account_links` 混在主对象里
- `activeKeyCount` 和 `activeBindingCount` 必须按状态集统计

### 6.4 Assertion Points

- 成功响应顶层字段必须包含 `code`、`message`、`data`、`timestamp`
- 列表响应 `data` 必须包含 `list`、`total`、`pageNum`、`pageSize`、`pages`
- 所有时间字符串字段必须统一格式
- JSON 字段命名必须使用小驼峰

---

### 7. Wrong vs Correct

#### Wrong

```sql
SELECT *
FROM public.users u
LEFT JOIN public.account_links al ON al.portal_user_id = u.id
LEFT JOIN public.user_lobster_profiles up ON up.user_id = u.id
WHERE up.user_id IS NOT NULL OR al.portal_user_id IS NOT NULL;
```

问题：

- 若同一用户有多条 `account_links`，会出现重复行
- 只能查到有 profile 或 accountLink 的用户，漏掉只有 key / binding 的用户

#### Correct

```sql
WITH candidate_users AS (
  SELECT user_id FROM public.user_lobster_profiles
  UNION
  SELECT portal_user_id AS user_id FROM public.account_links
  UNION
  SELECT user_id FROM public.follower_key_records
  UNION
  SELECT user_id FROM public.follower_node_bindings
),
current_account_link AS (
  SELECT *
  FROM (
    SELECT
      al.*,
      ROW_NUMBER() OVER (
        PARTITION BY al.portal_user_id
        ORDER BY
          CASE WHEN al.link_status = 'LINKED' THEN 0 ELSE 1 END,
          al.updated_at DESC,
          al.id DESC
      ) AS rn
    FROM public.account_links al
  ) t
  WHERE t.rn = 1
)
SELECT ...
FROM candidate_users cu
JOIN public.users u ON u.id = cu.user_id
LEFT JOIN current_account_link cal ON cal.portal_user_id = u.id
LEFT JOIN public.user_lobster_profiles up ON up.user_id = u.id;
```

#### Wrong

```go
// 在 follower user PATCH 中直接改 follower_key_records 和 follower_node_bindings
```

问题：

- 写入边界失控
- 该模块与 key / binding 流程耦合
- 历史数据容易被错误覆盖

#### Correct

```go
// 该模块只写 users / user_lobster_profiles / account_links
// follower_key_records / follower_node_bindings 仅用于只读聚合展示
```

---

## Design Decision: account_links 不做 user 维度硬唯一

**Context**：当前 `public.account_links` 仅有主键和普通索引，没有 `portal_user_id` 唯一约束。

**Decision**：一期实现不新增数据库唯一约束，也不在本模块中清洗历史数据，而是通过“当前有效账户关联选择规则”稳定读取 1 条记录。

**Why**：

- 避免在现有脏数据场景下强推 schema 变更
- 先保证列表与详情结果稳定
- 后续如业务确认“一人一 Lobster 账户”后，再补唯一约束和数据清洗任务

## Common Mistake: 把所有 users 都当成小龙虾用户

**Symptom**：列表数量远大于运营预期，且大量用户没有任何 Lobster 信息。

**Cause**：查询时直接以 `users` 为主表全量扫描，没有使用候选用户集合。

**Fix**：必须先构造 `candidate_users`，再回查 `users` 主档。

**Prevention**：将候选用户 CTE 封装为统一查询片段，不允许在多个查询里重复手写不同版本。

