# 运营系统核心 CRUD 需求文档

## 1. 文档目标

本文档用于指导后端开发实现运营系统中的 3 个核心模块：

1. 小龙虾用户管理
2. 大龙虾用户管理
3. 大龙虾 Key 管理

文档重点覆盖以下内容：

- 业务对象定义
- 一期范围
- CRUD 接口建议
- 写入表与查询表
- 查询口径
- 校验规则
- 删除语义
- 验收标准

## 2. 业务对象定义

### 2.1 小龙虾用户

基于当前数据库结构，本文将“小龙虾用户”定义为跟单侧用户或普通 Lobster 使用者，满足以下任一条件即可视为小龙虾用户：

- 在 `public.user_lobster_profiles` 中存在配置档案
- 在 `public.account_links` 中存在 Lobster 账户关联
- 在 `public.follower_key_records` 中存在 follower key
- 在 `public.follower_node_bindings` 中存在节点绑定关系

### 2.2 大龙虾用户

基于当前数据库结构，本文将“大龙虾用户”定义为节点/槽位拥有者，满足以下任一条件即可视为大龙虾用户：

- 在 `public.lobster_slots` 中作为 `owner_user_id` 拥有槽位
- 在 `public.node_key_records` 中作为 `owner_user_id` 拥有节点 key

### 2.3 大龙虾 Key

本文将“大龙虾 Key”定义为 leader 节点侧 key，主记录来源为：

- `public.node_key_records`

以下表作为扩展查询信息来源：

- `public.unified_key_registry`
- `public.model_license_key_records`
- `public.asset_node`

### 2.4 重要说明

- 当前 SQL dump 中没有现成的业务注释和枚举说明。
- 本文档中的业务定义基于表结构推断，可直接作为开发输入，但上线前仍需业务方确认状态枚举和值域。
- 一期默认不执行物理删除，统一采用逻辑删除/停用。

## 3. 一期范围

### 3.1 本期必须支持

- 运营后台按条件分页查询小龙虾用户
- 运营后台按条件分页查询大龙虾用户
- 运营后台查看大龙虾用户详情
- 运营后台创建大龙虾用户基础档案与初始槽位
- 运营后台更新大龙虾用户基础档案与槽位配置
- 运营后台停用大龙虾用户
- 运营后台按条件分页查询大龙虾 Key
- 运营后台查看大龙虾 Key 详情
- 运营后台创建大龙虾 Key
- 运营后台更新大龙虾 Key
- 运营后台停用大龙虾 Key

### 3.2 本期不做

- 小龙虾用户详情
- 小龙虾用户新增
- 小龙虾用户更新
- 小龙虾用户删除/停用
- 自动审批流程
- 批量导入导出
- 复杂报表
- 跨系统消息通知
- 物理删除历史数据

## 4. 通用实现规则

### 4.1 接口风格

建议统一使用后台管理接口前缀：

```text
/admin/v1
```

### 4.2 通用能力

所有列表接口必须支持：

- 分页：`pageNum`、`pageSize`
- 关键字搜索：`keyword`
- 状态筛选：`status`
- 时间范围筛选：`createdFrom`、`createdTo`
- 排序：`sortBy`、`sortOrder`

### 4.3 审计要求

所有创建、更新、删除操作都需要记录操作审计，建议写入：

- `public.audit_logs`

建议审计字段至少包含：

- 操作人 ID
- 操作对象类型
- 操作对象 ID
- 操作类型
- 变更前快照
- 变更后快照
- trace_id

### 4.4 删除语义

本期“删除”统一定义为逻辑删除：

- 不执行 `DELETE FROM`
- 优先更新业务状态字段为停用/失效
- 若表中无状态字段，则仅停用主对象，并保留历史从表记录

### 4.5 事务要求

以下场景必须使用事务：

- 创建用户主档案时，同时创建扩展档案
- 更新用户主档案时，同时更新扩展档案
- 创建大龙虾用户时，同时创建初始槽位
- 创建大龙虾 Key 时，同时同步注册表

### 4.6 状态枚举建议

以下枚举是建议值，最终以现网约定为准：

- 用户状态：`ACTIVE`、`DISABLED`
- 账户关联状态：`LINKED`、`UNLINKED`
- 槽位状态：`AVAILABLE`、`RESERVED`、`ACTIVE`、`DISABLED`
- follower key 状态：`ISSUED`、`BOUND`、`EXPIRED`、`REVOKED`
- follower 绑定状态：`PENDING`、`BOUND`、`UNBOUND`
- 节点 key 状态：`ISSUED`、`ACTIVE`、`EXPIRED`、`REVOKED`

## 5. 数据表映射

| 模块 | 主写表 | 扩展读取表 | 说明 |
| --- | --- | --- | --- |
| 小龙虾用户 | `无（一期只读）` | `public.users`、`public.user_lobster_profiles`、`public.account_links`、`public.follower_node_bindings`、`public.mirror_configs` | 小龙虾列表聚合查询 |
| 大龙虾用户 | `public.users`、`public.lobster_slots` | `public.node_key_records`、`public.asset_node`、`public.model_license_key_records`、`public.activation_revenue_records` | 用户主档与槽位归属 |
| 大龙虾 Key | `public.node_key_records` | `public.unified_key_registry`、`public.model_license_key_records`、`public.asset_node`、`public.key_application_requests` | 节点 key 主记录 |

## 6. 模块一：小龙虾用户管理

> 详细后端可执行规格见：[小龙虾用户列表代码规格](/Users/leon/Desktop/my/operation_admin/.trellis/spec/backend/operator-follower-user-crud.md)

## 6.1 模块目标

一期只提供小龙虾用户列表查询，用于运营查看小龙虾用户经营状态。

页面只显示以下 6 个字段：

- 用户ID
- 账户状态
- 策略状态
- 绑定状态
- 责任域
- 更新时间

## 6.2 主接口列表

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/admin/v1/follower-users` | 小龙虾用户列表 |

## 6.3 列表接口需求

### 6.3.1 筛选条件

- `keyword`
- `accountStatus`
- `strategyStatus`
- `bindingStatus`
- `responsibilityDomain`
- `updatedFrom`
- `updatedTo`

### 6.3.2 返回字段

- `userId`
- `accountStatus`
- `strategyStatus`
- `bindingStatus`
- `responsibilityDomain`
- `updatedAt`

### 6.3.3 查询表与查询口径

候选用户集合来源：

- `public.user_lobster_profiles`
- `public.account_links`
- `public.follower_node_bindings`
- `public.mirror_configs`

展示字段口径：

- `userId`：`public.users.external_user_id`
- `accountStatus`：`public.account_links.link_status`，无记录返回 `UNLINKED`
- `strategyStatus`：由 `public.mirror_configs` 派生
- `bindingStatus`：`public.follower_node_bindings.status`，无记录返回 `UNBOUND`
- `responsibilityDomain`：优先从 `public.user_lobster_profiles.meta_json` 读取，其次从 `public.account_links.meta_json` 读取
- `updatedAt`：取 `users`、`user_lobster_profiles`、`account_links`、`follower_node_bindings`、`mirror_configs` 中最大 `updated_at`

查询约束：

- 必须先构造候选用户集合，不能直接扫全量 `public.users`
- `public.account_links`、`public.follower_node_bindings`、`public.mirror_configs` 都按用户唯一，可直接左连接
- 不读取 `public.follower_key_records` 和 `public.user_lobster_binding_requests`

## 6.4 一期不做

- 详情接口
- 新增接口
- 更新接口
- 删除/停用接口
- follower key 明细
- 绑定申请明细

## 6.5 验收标准

- 列表页只返回 6 个字段，不额外返回用户详情字段
- 只要用户存在于候选表任意一张，就必须能进入列表
- 只有 `mirror_configs` 的用户也必须能查出
- 同一用户同时存在于多张候选表时只能返回 1 行
- `responsibilityDomain` 不存在时返回 `null`

## 7. 模块二：大龙虾用户管理

## 7.1 模块目标

运营可以查询、创建、修改、停用大龙虾用户，并查看其名下槽位、节点、节点 key 和相关扩展信息。

## 7.2 主接口列表

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/admin/v1/leader-users` | 大龙虾用户列表 |
| `GET` | `/admin/v1/leader-users/{user_id}` | 大龙虾用户详情 |
| `POST` | `/admin/v1/leader-users` | 创建大龙虾用户 |
| `PATCH` | `/admin/v1/leader-users/{user_id}` | 更新大龙虾用户 |
| `DELETE` | `/admin/v1/leader-users/{user_id}` | 停用大龙虾用户 |

## 7.3 列表接口需求

### 7.3.1 筛选条件

- `user_id`
- `external_user_id`
- `display_name`
- `slot_status`
- `slot_no`
- `has_node_key`
- `node_key_status`

### 7.3.2 返回字段

- `user_id`
- `external_user_id`
- `display_name`
- `user_status`
- `slot_count`
- `occupied_slot_count`
- `node_key_count`
- `active_node_key_count`
- `latest_node_key_expire_at_ts`
- `last_node_seen_at`

### 7.3.3 查询表与查询口径

主查询表：

- `public.users u`

统计表：

- `public.lobster_slots ls`，按 `ls.owner_user_id = u.id`
- `public.node_key_records nkr`，按 `nkr.owner_user_id = u.id`
- `public.asset_node an`，按 `an.node_id = nkr.node_id`

建议查询策略：

1. 先统计每个用户拥有的槽位数量
2. 再统计每个用户拥有的节点 key 数量
3. 用 `asset_node.last_seen_at` 补充节点最近在线时间

建议伪 SQL：

```sql
SELECT
  u.id,
  u.external_user_id,
  u.display_name,
  u.status,
  slot_stat.slot_count,
  slot_stat.occupied_slot_count,
  key_stat.node_key_count,
  key_stat.active_node_key_count,
  key_stat.latest_node_key_expire_at_ts,
  key_stat.last_node_seen_at
FROM public.users u
LEFT JOIN (...) slot_stat ON slot_stat.user_id = u.id
LEFT JOIN (...) key_stat ON key_stat.user_id = u.id
WHERE slot_stat.user_id IS NOT NULL OR key_stat.user_id IS NOT NULL;
```

## 7.4 详情接口需求

详情页需要返回 5 个信息区块：

1. 用户基础信息：来自 `public.users`
2. 槽位列表：来自 `public.lobster_slots`
3. 节点 key 列表：来自 `public.node_key_records`
4. 节点信息：来自 `public.asset_node`
5. 模型许可证 key 列表：来自 `public.model_license_key_records`

## 7.5 创建接口需求

### 7.5.1 写入表

- `public.users`
- `public.lobster_slots`

### 7.5.2 入参

| 字段 | 是否必填 | 写入表 | 说明 |
| --- | --- | --- | --- |
| `external_user_id` | 是 | `users` | 外部用户 ID，必须唯一 |
| `display_name` | 是 | `users` | 用户昵称 |
| `status` | 否 | `users` | 默认 `ACTIVE` |
| `preferred_channel` | 否 | `users` | 偏好通知渠道 |
| `locale` | 否 | `users` | 地区 |
| `initial_slots` | 是 | `lobster_slots` | 初始槽位数组，至少 1 个 |

### 7.5.3 `initial_slots` 子字段

| 字段 | 是否必填 | 说明 |
| --- | --- | --- |
| `slot_no` | 是 | 槽位编号，同一 owner 下唯一 |
| `slot_code` | 是 | 槽位编码，全局唯一 |
| `title` | 是 | 槽位标题 |
| `status` | 否 | 默认 `AVAILABLE` |
| `meta_json` | 否 | 扩展信息 |

### 7.5.4 校验规则

- `external_user_id` 唯一
- `initial_slots` 至少 1 条
- `slot_code` 必须唯一
- 同一用户下 `slot_no` 不能重复

### 7.5.5 创建逻辑

1. 插入 `users`
2. 批量插入 `lobster_slots`
3. 返回完整详情

## 7.6 更新接口需求

### 7.6.1 允许修改字段

- `users.display_name`
- `users.status`
- `users.preferred_channel`
- `users.locale`

### 7.6.2 槽位更新规则

大龙虾用户更新接口允许同时传入 `slot_operations`：

- `create`
- `update`
- `disable`

建议后端在同一事务内处理。

### 7.6.3 槽位写入表

- `public.lobster_slots`

### 7.6.4 槽位修改约束

- 不允许修改 `owner_user_id`
- `slot_no` 修改后仍需满足 `(owner_user_id, slot_no)` 唯一
- `slot_code` 修改后仍需满足唯一

## 7.7 删除接口需求

### 7.7.1 删除语义

执行逻辑停用：

- 将 `users.status` 更新为停用状态

### 7.7.2 默认不自动级联修改

本期默认不自动停用以下数据：

- `lobster_slots`
- `node_key_records`

原因：

- 用户停用和资产停用可能不是同一步审批
- 节点和槽位存在单独运营流程

### 7.7.3 可选增强

后端可以预留参数：

```json
{
  "cascade_disable_slots": false
}
```

后续如业务确认需要，可支持联动停用槽位。

## 7.8 验收标准

- 可以按昵称、外部用户 ID、槽位编号查询
- 创建后能看到初始槽位
- 更新后槽位数量与用户信息保持一致
- 停用后用户不再出现在默认启用列表中

## 8. 模块三：大龙虾 Key 管理

## 8.1 模块目标

运营可以查询、创建、修改、停用 leader 节点 key，并查看 key 对应的节点、注册表和模型 key 扩展信息。

## 8.2 主接口列表

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| `GET` | `/admin/v1/leader-keys` | 大龙虾 Key 列表 |
| `GET` | `/admin/v1/leader-keys/{node_id}` | 大龙虾 Key 详情 |
| `POST` | `/admin/v1/leader-keys` | 创建大龙虾 Key |
| `PATCH` | `/admin/v1/leader-keys/{node_id}` | 更新大龙虾 Key |
| `DELETE` | `/admin/v1/leader-keys/{node_id}` | 停用大龙虾 Key |

## 8.3 列表接口需求

### 8.3.1 筛选条件

- `node_id`
- `owner_user_id`
- `display_name`
- `status`
- `key_source`
- `expire_from`
- `expire_to`
- `key_material_sha256`

### 8.3.2 返回字段

- `node_id`
- `owner_user_id`
- `owner_display_name`
- `status`
- `key_source`
- `activated_at_ts`
- `expire_at_ts`
- `key_material_sha256`
- `dashboard_url`
- `node_name`
- `node_status`
- `last_node_seen_at`
- `registry_key_ref`
- `registry_key_type`
- `registry_status`

### 8.3.3 查询表与查询口径

主查询表：

- `public.node_key_records nkr`

扩展关联表：

- `public.users u`，按 `u.id = nkr.owner_user_id`
- `public.asset_node an`，按 `an.node_id = nkr.node_id`
- `public.unified_key_registry ukr`，按 `ukr.owner_user_id = nkr.owner_user_id AND ukr.leader_node_id = nkr.node_id`

建议伪 SQL：

```sql
SELECT
  nkr.node_id,
  u.id AS owner_user_id,
  u.display_name AS owner_display_name,
  nkr.status,
  nkr.key_source,
  nkr.activated_at_ts,
  nkr.expire_at_ts,
  nkr.key_material_sha256,
  nkr.dashboard_url,
  an.name AS node_name,
  an.status AS node_status,
  an.last_seen_at,
  ukr.key_ref,
  ukr.key_type,
  ukr.status AS registry_status
FROM public.node_key_records nkr
JOIN public.users u ON u.id = nkr.owner_user_id
LEFT JOIN public.asset_node an ON an.node_id = nkr.node_id
LEFT JOIN public.unified_key_registry ukr
  ON ukr.owner_user_id = nkr.owner_user_id
 AND ukr.leader_node_id = nkr.node_id;
```

## 8.4 详情接口需求

详情页需要返回 4 个信息区块：

1. 节点 key 主信息：来自 `public.node_key_records`
2. 注册表信息：来自 `public.unified_key_registry`
3. 节点信息：来自 `public.asset_node`
4. 模型许可证 key 列表：来自 `public.model_license_key_records`

## 8.5 创建接口需求

### 8.5.1 主写表

- `public.node_key_records`

### 8.5.2 可选同步表

- `public.unified_key_registry`

### 8.5.3 入参

| 字段 | 是否必填 | 写入表 | 说明 |
| --- | --- | --- | --- |
| `node_id` | 是 | `node_key_records` | 节点 ID，必须唯一 |
| `owner_user_id` | 是 | `node_key_records` | 所属大龙虾用户 ID |
| `key_material_sha256` | 是 | `node_key_records` | key 指纹摘要 |
| `status` | 否 | `node_key_records` | 默认 `ISSUED` |
| `activated_at_ts` | 否 | `node_key_records` | 激活时间 |
| `expire_at_ts` | 否 | `node_key_records` | 过期时间 |
| `key_source` | 否 | `node_key_records` | key 来源 |
| `dashboard_url` | 否 | `node_key_records` | 管理地址 |
| `meta_json` | 否 | `node_key_records` | 扩展字段 |
| `registry_sync` | 否 | 控制字段 | 是否同步注册表 |
| `registry_payload` | 否 | `unified_key_registry` | 注册表写入内容 |

### 8.5.4 校验规则

- `node_id` 不能为空，且需满足 `node_key_records.node_id` 唯一约束
- `owner_user_id` 必须存在于 `users`
- `expire_at_ts` 不能早于 `activated_at_ts`

### 8.5.5 创建逻辑

1. 插入 `node_key_records`
2. 若 `registry_sync=true`，同步插入 `unified_key_registry`
3. 返回完整详情

## 8.6 更新接口需求

### 8.6.1 允许修改字段

- `status`
- `activated_at_ts`
- `expire_at_ts`
- `key_source`
- `dashboard_url`
- `meta_json`

### 8.6.2 注册表同步

若请求中显式传入 `registry_payload`，则同步更新 `unified_key_registry` 中当前节点对应记录。

建议匹配条件：

- `owner_user_id`
- `leader_node_id`
- `key_ref` 或当前生效记录

## 8.7 删除接口需求

### 8.7.1 删除语义

执行逻辑停用：

- 将 `node_key_records.status` 更新为停用状态
- 如启用注册表同步，则同步更新 `unified_key_registry.status`

### 8.7.2 不执行物理删除

原因：

- key 历史需要审计
- 节点 key 可能关联其他业务流程

## 8.8 验收标准

- 可按 `node_id`、owner 用户、状态、过期时间检索 key
- 创建后可在列表页立即查到
- 更新状态后详情页立即反映
- 删除后默认列表不再返回已停用 key

## 9. 错误处理约束

统一遵循 [错误处理](/Users/leon/Desktop/my/operation_admin/.trellis/spec/backend/error-handling.md)：

- 业务错误码使用数值型，不使用字符串型错误码
- 参数校验失败统一返回 `40001`
- 资源不存在统一返回 `404`
- 更细的业务错误码按模块细分，例如：
  - `40011`：外部用户 ID 重复
  - `40012`：配置档案不存在
  - `40013`：非法状态流转

## 10. 建议索引使用

开发时优先使用现有索引：

- `ix_users_status`
- `ix_account_links_portal_user_id`
- `ix_account_links_lobster_user_id`
- `ix_follower_key_records_user_id`
- `ix_follower_key_records_status`
- `ix_follower_node_bindings_user_id`
- `ix_follower_node_bindings_status`
- `ix_lobster_slots_owner_user_id`
- `uq_lobster_slots_owner_slot_no`
- `ix_node_key_records_owner_user_id`
- `ix_node_key_records_status`
- `ix_node_key_records_expire_at_ts`
- `ix_unified_key_registry_owner_user_id`
- `ix_unified_key_registry_leader_node_id`

## 11. 待业务确认项

以下项目不阻塞后端开始开发接口，但上线前必须确认：

1. “小龙虾 / 大龙虾”在业务上的正式定义是否与本文一致
2. 各状态字段的正式枚举值
3. 停用用户时是否需要自动停用槽位和 key
4. `unified_key_registry` 的同步策略是强同步还是可选同步
5. `asset_node.owner_user_id` 字段类型与 `users.id` 不一致，当前不建议作为强关联依据，需确认是否存在额外映射表

## 12. 后端拆分建议

建议后端按以下 3 个 service 实现：

1. `FollowerUserService`
2. `LeaderUserService`
3. `LeaderKeyService`

建议额外封装 2 个通用能力：

1. `AuditLogService`
2. `PaginationQueryBuilder`
