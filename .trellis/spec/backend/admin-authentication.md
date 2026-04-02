# 后台开发者鉴权规范

> 本项目后台开发者登录、令牌刷新、登出、当前用户和管理员建号的一期契约。

---

## 概览

当前后台开发者鉴权统一基于 `backend/internal/adminauth` 模块实现，接口统一挂在 `/admin/v1` 下，响应结构继续复用 `backend/internal/http/response` 中的 `code`、`message`、`data`、`timestamp` 契约。

当前真实实现位置：

- `backend/internal/adminauth`
- `backend/internal/http/middleware/admin_auth.go`
- `backend/internal/http/handler/auth.go`
- `backend/internal/http/handler/admin_user.go`
- `backend/internal/http/router/router.go`

---

## 数据模型

一期新增三张表，统一落在 `operator_portal` schema 下：

1. `operator_portal.admin_users`
2. `operator_portal.admin_refresh_tokens`
3. `operator_portal.admin_login_audits`

后台账号 ID 不使用自增整数直接对外暴露，而是通过 `operator_portal.admin_user_seq` 序列生成 `ADM-1001` 这类字符串 ID。

关键约束：

- `login_name` 全局唯一，入库前统一转小写
- `role_code` 当前固定为 `admin`
- `status` 仅允许 `active`、`disabled`
- `refresh_token_hash` 只存 SHA-256 哈希，不存原始 refresh token
- `password_hash` 只存 bcrypt 哈希，不存明文密码
- 登录审计表禁止记录密码、JWT 原文、refresh token 原文和密码哈希

---

## 接口约定

### 公开接口

- `POST /admin/v1/auth/login`
- `POST /admin/v1/auth/refresh`
- `GET /admin/v1/healthz`

### 受保护接口

- `POST /admin/v1/auth/logout`
- `GET /admin/v1/auth/me`
- `POST /admin/v1/admin-users`
- 其余 `/admin/v1` 业务接口，例如现有 `follower-users`

### 登录

请求体：

```json
{
  "loginName": "admin.root",
  "password": "12345678"
}
```

成功返回：

- `accessToken`
- `accessTokenExpiresIn`：单位为秒
- `refreshToken`
- `refreshTokenExpiresIn`：单位为秒
- `user`

失败规则：

- `loginName` 不合法或密码少于 8 位：`40001`
- 账号不存在、已删除或密码错误：`401`
- 账号状态为 `disabled`：`403`

### 刷新令牌

请求体：

```json
{
  "refreshToken": "..."
}
```

规则：

- 旧 Refresh Token 在刷新成功后立即撤销
- 过期、伪造、已撤销的 Refresh Token 统一返回 `401`

### 登出

Header：

```text
Authorization: Bearer <access token>
```

请求体：

```json
{
  "refreshToken": "..."
}
```

规则：

- 仅撤销当前后台账号名下与请求体匹配的 Refresh Token
- 找不到、已过期、已撤销都按成功处理
- 返回数据固定为：

```json
{
  "success": true
}
```

### 当前用户

Header：

```text
Authorization: Bearer <access token>
```

返回字段：

- `adminUserId`
- `loginName`
- `displayName`
- `roleCode`
- `status`
- `lastLoginAt`

### 管理员创建后台账号

Header：

```text
Authorization: Bearer <access token>
```

请求体：

```json
{
  "loginName": "dev.user",
  "password": "12345678",
  "displayName": "开发同学"
}
```

规则：

- 只有 `roleCode=admin` 可以调用
- 新账号固定写入 `role_code=admin`
- 新账号默认 `status=active`
- `loginName` 重复返回业务码 `40010`

---

## Token 与安全

- Access Token：JWT HS256
- Claims：`sub`、`loginName`、`roleCode`、`tokenType`、`jti`、`iat`、`exp`、`iss`
- `tokenType` 固定为 `access`
- Access Token 默认有效期：`2h`
- Refresh Token 默认有效期：`7d`
- Access Token 中间件除了校验 JWT 本身，还会回表确认后台账号仍是 `active` 且 `is_deleted=false`
- 后台账号被禁用后，旧 Access Token 会在下一次访问受保护接口时被拒绝

---

## 初始化方式

当前项目延续已有后端风格：

1. 服务启动时调用 `EnsureSchema + AutoMigrate` 自动创建 schema、序列和鉴权表
2. 首个管理员账号不由应用自动生成
3. 首个管理员通过 `backend/sql/seeds/0001_admin_user_seed.sql` 模板手工替换后执行

---

## 常见错误

- 在日志里输出原始密码、JWT、refresh token 或 bcrypt 哈希
- 新增后台接口时忘记放到受保护分组，导致匿名访问
- 把账号不存在和密码错误拆成两套错误文案，暴露账号枚举风险
- 刷新成功后不撤销旧 Refresh Token，导致同一令牌可重复使用
- Access Token 只验签不查库，导致禁用账号后旧 token 仍然可继续访问后台接口
