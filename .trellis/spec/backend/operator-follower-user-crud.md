# 小龙虾用户列表规格

> 场景：运营后台查询小龙虾用户列表。
>
> 当前版本为一期精简版，只做列表查询，只返回 6 个展示字段，不做详情、不做新增、不做编辑、不做删除。

---

## Scenario: 小龙虾用户列表查询

### 1. Scope / Trigger

- Trigger：运营后台需要查看小龙虾用户基础经营状态。
- Scope：只提供列表查询接口。
- Not In Scope：
  - 详情接口
  - 创建接口
  - 更新接口
  - 删除接口
  - follower key 明细
  - 绑定申请明细

### 1.1 一期页面只显示字段

小龙虾列表页只显示以下 6 个字段：

1. 用户ID
2. 账户状态
3. 策略状态
4. 绑定状态
5. 责任域
6. 更新时间

### 1.2 数据来源约束

本模块必须基于现有表做聚合，不新增库表，不修改现有 schema。

允许读取的表：

- `public.users`
- `public.user_lobster_profiles`
- `public.account_links`
- `public.follower_node_bindings`
- `public.mirror_configs`

本模块不读取以下表作为主口径：

- `public.follower_key_records`
- `public.user_lobster_binding_requests`

原因：

- 当前一期展示字段不需要 key 和申请明细
- 加入这些表会扩大候选用户范围，导致列表噪音增加

---

### 2. Signatures

### 2.1 HTTP Endpoint

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/admin/v1/follower-users` | 查询小龙虾用户分页列表 |

### 2.2 Query Params

```text
GET /admin/v1/follower-users?pageNum=1&pageSize=20&keyword=lobster_u_10001&accountStatus=LINKED&strategyStatus=RUNNING&bindingStatus=BOUND&responsibilityDomain=cn-east
```

字段定义：

- `pageNum`：默认 `1`
- `pageSize`：默认 `20`，最大 `100`
- `keyword`：支持匹配用户 ID
- `accountStatus`
- `strategyStatus`
- `bindingStatus`
- `responsibilityDomain`
- `updatedFrom`
- `updatedTo`
- `sortBy`：默认 `updatedAt`
- `sortOrder`：默认 `desc`

### 2.3 Handler / Service 建议签名

```go
type FollowerUserListService interface {
    List(ctx context.Context, req ListFollowerUsersRequest) (PageResult[FollowerUserListItem], error)
}
```

---

### 3. Contracts

### 3.1 用户候选集合

小龙虾列表不能直接扫全量 `users`，必须先构造候选用户集合。

候选规则：

- 在 `public.user_lobster_profiles` 中存在记录
- 或在 `public.account_links` 中存在记录
- 或在 `public.follower_node_bindings` 中存在记录
- 或在 `public.mirror_configs` 中存在记录

推荐 SQL：

```sql
WITH candidate_users AS (
  SELECT user_id FROM public.user_lobster_profiles
  UNION
  SELECT portal_user_id AS user_id FROM public.account_links
  UNION
  SELECT user_id FROM public.follower_node_bindings
  UNION
  SELECT user_id FROM public.mirror_configs
)
```

### 3.2 展示字段映射

| 展示字段 | API 字段 | 来源 | 规则 |
| --- | --- | --- | --- |
| 用户ID | `userId` | `users.external_user_id` | 页面展示固定取外部用户 ID，不展示内部自增主键 |
| 账户状态 | `accountStatus` | `account_links.link_status` | 取“当前有效账户关联”状态；无记录时返回 `UNLINKED` |
| 策略状态 | `strategyStatus` | `mirror_configs` 派生 | 按策略状态派生规则计算 |
| 绑定状态 | `bindingStatus` | `follower_node_bindings.status` | 取“当前有效绑定记录”状态；无记录时返回 `UNBOUND` |
| 责任域 | `responsibilityDomain` | `user_lobster_profiles.meta_json` 或 `account_links.meta_json` | 按责任域提取规则读取；无值时返回 `null` |
| 更新时间 | `updatedAt` | 多表派生 | 取用户相关关键表的最新更新时间 |

### 3.3 账户关联读取规则

`public.account_links` 在数据库中存在唯一约束 `uq_account_links_portal_user UNIQUE (portal_user_id)`。

固定规则：

1. 一个 `portal_user_id` 最多只有 1 条账户关联记录
2. 查询时允许直接按 `al.portal_user_id = u.id` 左连接
3. 若无记录，`accountStatus` 返回 `UNLINKED`
4. 一期不做历史账户关联择优逻辑

推荐 SQL：

```sql
LEFT JOIN public.account_links al
  ON al.portal_user_id = u.id
```

### 3.4 当前有效绑定规则

`public.follower_node_bindings` 在数据库中存在唯一约束 `uq_follower_node_binding_user UNIQUE (user_id)`。

固定规则：

1. 一个 `user_id` 最多只有 1 条绑定记录
2. 查询时允许直接按 `fnb.user_id = u.id` 左连接
3. 若有记录，`bindingStatus` 直接取 `fnb.status`
4. 若无记录，返回 `UNBOUND`

推荐 SQL：

```sql
LEFT JOIN public.follower_node_bindings fnb
  ON fnb.user_id = u.id
```

### 3.5 策略状态派生规则

策略状态来源于 `public.mirror_configs`，不是直接存储字段，而是计算字段。

固定派生规则：

```text
无 mirror_configs 记录         -> UNCONFIGURED
enabled = false               -> DISABLED
risk_circuit_open = true      -> RISK_HOLD
stop_pending_settlement = true -> STOPPING
enabled = true                -> RUNNING
其他未知组合                  -> UNKNOWN
```

推荐 SQL：

```sql
CASE
  WHEN mc.user_id IS NULL THEN 'UNCONFIGURED'
  WHEN mc.enabled = false THEN 'DISABLED'
  WHEN mc.risk_circuit_open = true THEN 'RISK_HOLD'
  WHEN mc.stop_pending_settlement = true THEN 'STOPPING'
  WHEN mc.enabled = true THEN 'RUNNING'
  ELSE 'UNKNOWN'
END AS strategy_status
```

### 3.6 责任域提取规则

当前数据库中没有独立的“责任域”结构化字段，一期按以下优先级提取：

1. `user_lobster_profiles.meta_json ->> 'responsibilityDomain'`
2. `user_lobster_profiles.meta_json ->> 'ownerDomain'`
3. `account_links.meta_json::jsonb ->> 'responsibilityDomain'`
4. `account_links.meta_json::jsonb ->> 'ownerDomain'`
5. 都不存在则返回 `null`

推荐 SQL：

```sql
COALESCE(
  up.meta_json ->> 'responsibilityDomain',
  up.meta_json ->> 'ownerDomain',
  NULLIF(BTRIM(al.meta_json), '')::jsonb ->> 'responsibilityDomain',
  NULLIF(BTRIM(al.meta_json), '')::jsonb ->> 'ownerDomain'
) AS responsibility_domain
```

实现要求：

- 若 `account_links.meta_json` 为空字符串或只有空白字符，不允许直接强转 `jsonb`
- 必须使用 `NULLIF(BTRIM(al.meta_json), '')::jsonb`

### 3.7 更新时间规则

更新时间不是单表字段，固定取以下时间的最大值：

- `users.updated_at`
- `user_lobster_profiles.updated_at`
- 当前有效 `account_links.updated_at`
- 当前有效 `follower_node_bindings.updated_at`
- `mirror_configs.updated_at`

推荐 SQL：

```sql
GREATEST(
  u.updated_at,
  COALESCE(up.updated_at, u.updated_at),
  COALESCE(al.updated_at, u.updated_at),
  COALESCE(fnb.updated_at, u.updated_at),
  COALESCE(mc.updated_at, u.updated_at)
) AS updated_at
```

### 3.8 推荐查询骨架

```sql
WITH candidate_users AS (
  SELECT user_id FROM public.user_lobster_profiles
  UNION
  SELECT portal_user_id AS user_id FROM public.account_links
  UNION
  SELECT user_id FROM public.follower_node_bindings
  UNION
  SELECT user_id FROM public.mirror_configs
)
SELECT
  u.external_user_id AS user_id,
  COALESCE(al.link_status, 'UNLINKED') AS account_status,
  CASE
    WHEN mc.user_id IS NULL THEN 'UNCONFIGURED'
    WHEN mc.enabled = false THEN 'DISABLED'
    WHEN mc.risk_circuit_open = true THEN 'RISK_HOLD'
    WHEN mc.stop_pending_settlement = true THEN 'STOPPING'
    WHEN mc.enabled = true THEN 'RUNNING'
    ELSE 'UNKNOWN'
  END AS strategy_status,
  COALESCE(fnb.status, 'UNBOUND') AS binding_status,
  COALESCE(
    up.meta_json ->> 'responsibilityDomain',
    up.meta_json ->> 'ownerDomain',
    NULLIF(BTRIM(al.meta_json), '')::jsonb ->> 'responsibilityDomain',
    NULLIF(BTRIM(al.meta_json), '')::jsonb ->> 'ownerDomain'
  ) AS responsibility_domain,
  GREATEST(
    u.updated_at,
    COALESCE(up.updated_at, u.updated_at),
    COALESCE(al.updated_at, u.updated_at),
    COALESCE(fnb.updated_at, u.updated_at),
    COALESCE(mc.updated_at, u.updated_at)
  ) AS updated_at
FROM candidate_users cu
JOIN public.users u ON u.id = cu.user_id
LEFT JOIN public.user_lobster_profiles up ON up.user_id = u.id
LEFT JOIN public.account_links al ON al.portal_user_id = u.id
LEFT JOIN public.follower_node_bindings fnb ON fnb.user_id = u.id
LEFT JOIN public.mirror_configs mc ON mc.user_id = u.id;
```

### 3.9 列表响应 Contract

```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "list": [
      {
        "userId": "lobster_u_10001",
        "accountStatus": "LINKED",
        "strategyStatus": "RUNNING",
        "bindingStatus": "BOUND",
        "responsibilityDomain": "cn-east",
        "updatedAt": "2026-04-01 16:30:00"
      }
    ],
    "total": 1,
    "pageNum": 1,
    "pageSize": 20,
    "pages": 1
  },
  "timestamp": 1775104225000
}
```

---

### 4. Validation & Error Matrix

| 场景 | 条件 | HTTP | code | message |
| --- | --- | --- | --- | --- |
| 列表查询 | `pageNum < 1` 或 `pageSize < 1` | `400` | `40001` | 参数校验失败：分页参数非法 |
| 列表查询 | `pageSize > 100` | `400` | `40001` | 参数校验失败：pageSize 超出上限 |
| 列表查询 | `sortBy` 非 `updatedAt` | `400` | `40001` | 参数校验失败：sortBy 非法 |
| 列表查询 | `sortOrder` 非 `asc/desc` | `400` | `40001` | 参数校验失败：sortOrder 非法 |
| 列表查询 | `updatedFrom > updatedTo` | `400` | `40001` | 参数校验失败：更新时间范围非法 |
| 列表查询 | SQL 执行失败 | `500` | `500` | 服务器内部错误，请稍后重试 |

允许的筛选枚举：

- `accountStatus`：按现网 `link_status` 值过滤
- `strategyStatus`：`UNCONFIGURED`、`DISABLED`、`RISK_HOLD`、`STOPPING`、`RUNNING`、`UNKNOWN`
- `bindingStatus`：按现网绑定状态值过滤，默认兼容 `BOUND`、`ACTIVE`、`UNBOUND`

---

### 5. Good / Base / Bad Cases

### 5.1 Good Case

场景：用户存在 profile、account link、mirror config、binding，全量信息齐全。

预期：

- `userId` 取 `users.external_user_id`
- `accountStatus` = 当前有效 `link_status`
- `strategyStatus` 按派生规则显示 `RUNNING`
- `bindingStatus` = 当前有效绑定状态
- `responsibilityDomain` = `meta_json` 中提取结果
- `updatedAt` = 多表最大更新时间

### 5.2 Base Case

场景：用户只有 `mirror_configs` 记录，没有 `account_links` 和 `follower_node_bindings`。

预期：

- 出现在列表中
- `accountStatus` = `UNLINKED`
- `strategyStatus` 根据 `mirror_configs` 派生
- `bindingStatus` = `UNBOUND`
- `responsibilityDomain` = `null`

### 5.3 Bad Case

场景：同一用户同时存在 `user_lobster_profiles`、`account_links`、`follower_node_bindings`、`mirror_configs` 4 张候选表。

预期：

- 列表中该用户只能出现 1 行
- 必须由 `candidate_users` 的 `UNION` 去重
- 4 张 1:1 扩展表均允许直接左连接

### 5.4 Bad Case

场景：`account_links.meta_json` 为空字符串。

预期：

- 查询不能因为 `::jsonb` 强转失败而报错
- `responsibilityDomain` 正常返回 `null`

---

### 6. Tests Required

### 6.1 Query Tests

- 列表必须对候选用户去重，不能因为用户同时存在于多张候选表而重复返回
- 列表必须能查出只有 `mirror_configs` 的用户
- 列表必须能查出只有 `follower_node_bindings` 的用户
- `updatedAt` 必须取多表最大值

### 6.2 Mapping Tests

- `userId` 必须返回 `users.external_user_id`，不是内部 `users.id`
- `accountStatus` 无记录时返回 `UNLINKED`
- `bindingStatus` 无记录时返回 `UNBOUND`
- `responsibilityDomain` 按优先级提取

### 6.3 Assertion Points

- 顶层响应必须包含 `code`、`message`、`data`、`timestamp`
- 分页响应必须包含 `list`、`total`、`pageNum`、`pageSize`、`pages`
- JSON 字段命名必须使用小驼峰

---

### 7. Wrong vs Correct

#### Wrong

```sql
SELECT
  u.external_user_id AS user_id,
  al.link_status,
  fnb.status
FROM public.users u
LEFT JOIN public.account_links al ON al.portal_user_id = u.id
LEFT JOIN public.follower_node_bindings fnb ON fnb.user_id = u.id;
```

问题：

- 会扫出不属于小龙虾候选集合的普通用户
- 没有策略状态
- 没有责任域提取
- 更新时间不准确

#### Correct

```sql
WITH candidate_users AS (...)
SELECT
  u.external_user_id AS user_id,
  COALESCE(al.link_status, 'UNLINKED') AS account_status,
  CASE ... END AS strategy_status,
  COALESCE(fnb.status, 'UNBOUND') AS binding_status,
  COALESCE(...) AS responsibility_domain,
  GREATEST(...) AS updated_at
FROM candidate_users cu
JOIN public.users u ON u.id = cu.user_id
LEFT JOIN public.user_lobster_profiles up ON up.user_id = u.id
LEFT JOIN public.account_links al ON al.portal_user_id = u.id
LEFT JOIN public.follower_node_bindings fnb ON fnb.user_id = u.id
LEFT JOIN public.mirror_configs mc ON mc.user_id = u.id;
```

---

## Design Decision: 用户ID固定显示 external_user_id

**Context**：`users` 同时有内部自增主键 `id` 和业务侧 `external_user_id`。

**Decision**：小龙虾列表页中的“用户ID”固定显示 `users.external_user_id`。

**Why**：

- 更符合运营识别习惯
- 与外部系统对账更方便
- `external_user_id` 当前有唯一索引且非空

## Design Decision: 账户状态固定映射 account_links.link_status

**Context**：`public.users.status` 和 `public.account_links.link_status` 都可能被误解成“账户状态”。

**Decision**：小龙虾列表页中的 `accountStatus` 固定映射 `public.account_links.link_status`，不映射 `public.users.status`。

**Why**：

- 当前列表核心目标是看小龙虾接入经营状态，而不是看后台主档是否启停
- 列表已经有 `strategyStatus` 与 `bindingStatus`，配套显示 Lobster 账户关联状态更完整
- `users.status` 若后续需要展示，应以新字段名 `userStatus` 单独增加，避免语义冲突

## Common Mistake: 把 users.status 当 accountStatus

**Symptom**：接口返回的 `accountStatus` 为 `ACTIVE` / `DISABLED`，与运营预期的 `LINKED` / `UNLINKED` 不一致。

**Cause**：把 `public.users.status` 误当作列表展示的账户状态。

**Fix**：`accountStatus` 固定读取 `public.account_links.link_status`；无记录时返回 `UNLINKED`。

**Prevention**：DTO 中显式区分 `accountStatus` 与未来可能新增的 `userStatus`，不要复用同一个字段名。

## Common Mistake: 把责任域当结构化列直接查询

**Symptom**：SQL 直接写 `up.responsibility_domain`，上线后报字段不存在。

**Cause**：当前数据库没有独立责任域列，责任域只能从 `meta_json` 提取。

**Fix**：按“责任域提取规则”从 `user_lobster_profiles.meta_json` 和 `account_links.meta_json` 中读取。

**Prevention**：后端 DTO 中把责任域定义成派生字段，不要映射成实体固定列。
