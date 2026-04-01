# 小龙虾用户 CRUD 简版

## 1. 页面字段来源

| 页面字段 | 示例值 | 来源表 | 来源字段 | 说明 |
| --- | --- | --- | --- | --- |
| `userId` | `U-1201` | `public.users` | `external_user_id` | 直接返回 |
| `accountStatus` | `active` | `public.users` | `status` | 建议统一转成小写返回 |
| `strategyStatus` | `enabled` | `public.mirror_configs` | `enabled` | `true -> enabled`，`false -> disabled`，无记录可返回 `unconfigured` |
| `bindingStatus` | `pending` | `public.follower_node_bindings` | `status` | 建议统一转成小写返回 |
| `responsibilityDomain` | `risk` | `public.user_lobster_profiles` | `meta_json ->> 'responsibilityDomain'` | 从 `meta_json` 读取 |
| `updatedAt` | `2026-03-31 15:30:25` | `users.updated_at` | `max(updated_at)` | 取 `users.updated_at`|

## 2. 统一返回对象

```json
{
  "userId": "U-1201",
  "accountStatus": "active",
  "strategyStatus": "enabled",
  "bindingStatus": "pending",
  "responsibilityDomain": "risk",
  "updatedAt": "2026-03-31 15:30:25"
}
```

### 2.1 返回字段来源

| 返回字段 | 来源表 | 来源字段 |
| --- | --- | --- |
| `userId` | `public.users` | `external_user_id` |
| `accountStatus` | `public.users` | `status` |
| `strategyStatus` | `public.mirror_configs` | `enabled` |
| `bindingStatus` | `public.follower_node_bindings` | `status` |
| `responsibilityDomain` | `public.user_lobster_profiles` | `meta_json ->> 'responsibilityDomain'` |
| `updatedAt` | `public.users`、`public.user_lobster_profiles`、`public.mirror_configs`、`public.follower_node_bindings` | 最大 `updated_at` |

## 3. Create

### 3.1 请求

- 方法：`POST`
- 路径：`/admin/v1/follower-users`

请求体：

```json
{
  "userId": "U-1201",
  "accountStatus": "active",
  "strategyStatus": "enabled",
  "bindingStatus": "pending",
  "responsibilityDomain": "risk"
}
```

### 3.2 请求参数来源

| 请求字段 | 写入表 | 写入字段 |
| --- | --- | --- |
| `userId` | `public.users` | `external_user_id` |
| `accountStatus` | `public.users` | `status` |
| `strategyStatus` | `public.mirror_configs` | `enabled` |
| `bindingStatus` | `public.follower_node_bindings` | `status` |
| `responsibilityDomain` | `public.user_lobster_profiles` | `meta_json.responsibilityDomain` |

### 3.3 返回

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "userId": "U-1201",
    "accountStatus": "active",
    "strategyStatus": "enabled",
    "bindingStatus": "pending",
    "responsibilityDomain": "risk",
    "updatedAt": "2026-03-31 15:30:25"
  }
}
```

返回字段来源：见“2.1 返回字段来源”。

## 4. Read

### 4.1 列表查询请求

- 方法：`GET`
- 路径：`/admin/v1/follower-users`

Query 参数：

| 参数 | 说明 |
| --- | --- |
| `pageNum` | 页码 |
| `pageSize` | 每页数量 |
| `userId` | 按用户 ID 查询 |
| `accountStatus` | 按账户状态查询 |
| `strategyStatus` | 按策略状态查询 |
| `bindingStatus` | 按绑定状态查询 |
| `responsibilityDomain` | 按责任域查询 |
| `updatedFrom` | 更新时间开始 |
| `updatedTo` | 更新时间结束 |

### 4.2 列表查询返回

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "userId": "U-1201",
        "accountStatus": "active",
        "strategyStatus": "enabled",
        "bindingStatus": "pending",
        "responsibilityDomain": "risk",
        "updatedAt": "2026-03-31 15:30:25"
      }
    ],
    "total": 1,
    "pageNum": 1,
    "pageSize": 20
  }
}
```

列表中每个字段来源：见“2.1 返回字段来源”。

### 4.3 详情查询请求

- 方法：`GET`
- 路径：`/admin/v1/follower-users/{userId}`

Path 参数：

| 参数 | 说明 |
| --- | --- |
| `userId` | 对应 `public.users.external_user_id` |

### 4.4 详情查询返回

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "userId": "U-1201",
    "accountStatus": "active",
    "strategyStatus": "enabled",
    "bindingStatus": "pending",
    "responsibilityDomain": "risk",
    "updatedAt": "2026-03-31 15:30:25"
  }
}
```

返回字段来源：见“2.1 返回字段来源”。

## 5. Update

### 5.1 请求

- 方法：`PATCH`
- 路径：`/admin/v1/follower-users/{userId}`

Path 参数：

| 参数 | 说明 |
| --- | --- |
| `userId` | 对应 `public.users.external_user_id` |

请求体：

```json
{
  "accountStatus": "active",
  "strategyStatus": "enabled",
  "bindingStatus": "pending",
  "responsibilityDomain": "risk"
}
```

### 5.2 请求参数来源

| 请求字段 | 写入表 | 写入字段 |
| --- | --- | --- |
| `accountStatus` | `public.users` | `status` |
| `strategyStatus` | `public.mirror_configs` | `enabled` |
| `bindingStatus` | `public.follower_node_bindings` | `status` |
| `responsibilityDomain` | `public.user_lobster_profiles` | `meta_json.responsibilityDomain` |

### 5.3 返回

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "userId": "U-1201",
    "accountStatus": "active",
    "strategyStatus": "enabled",
    "bindingStatus": "pending",
    "responsibilityDomain": "risk",
    "updatedAt": "2026-03-31 15:30:25"
  }
}
```

返回字段来源：见“2.1 返回字段来源”。

## 6. Delete

### 6.1 请求

- 方法：`DELETE`
- 路径：`/admin/v1/follower-users/{userId}`

Path 参数：

| 参数 | 说明 |
| --- | --- |
| `userId` | 对应 `public.users.external_user_id` |

### 6.2 删除动作

| 动作字段 | 写入表 | 写入字段 |
| --- | --- | --- |
| 账户状态置为停用 | `public.users` | `status` |
| 策略状态置为关闭 | `public.mirror_configs` | `enabled` |
| 绑定状态置为解绑 | `public.follower_node_bindings` | `status` |

### 6.3 返回

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "userId": "U-1201",
    "success": true
  }
}
```
