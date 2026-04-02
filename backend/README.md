# Go 后端基础框架

## 技术栈

- `gin`：HTTP 路由与中间件
- `gorm`：数据库访问层
- `sonic`：JSON 编解码
- `viper`：配置加载与环境变量覆盖
- `zap`：结构化日志
- `lumberjack.v2`：日志文件滚动切割

## 目录结构

```text
backend/
├── cmd/server                 # 应用启动入口
├── configs/config.yaml        # 默认配置文件
├── sql/seeds                  # 首个管理员账号初始化 SQL 模板
├── internal/adminauth         # 后台开发者鉴权模型、仓储、服务与 Token 组件
├── internal/app               # 应用装配与生命周期管理
├── internal/config            # 配置加载与配置契约
├── internal/database          # Gorm / PostgreSQL 初始化
├── internal/followeruser      # 小龙虾用户 Gorm 模型、仓储与服务
├── internal/http
│   ├── handler                # HTTP 处理器
│   ├── middleware             # Gin 鉴权中间件
│   ├── response               # Sonic JSON 响应封装
│   ├── router                 # 路由注册
├── internal/logger            # Zap + Lumberjack 日志初始化
├── logs                       # 本地日志目录
├── go.mod
└── Makefile
```

## 默认能力

- 提供 `GET /healthz` 健康检查接口
- 提供 `POST /admin/v1/auth/login`、`POST /admin/v1/auth/refresh`、`POST /admin/v1/auth/logout`、`GET /admin/v1/auth/me`
- 提供 `POST /admin/v1/admin-users` 作为管理员创建后台开发者账号接口
- 预留 `GET /admin/v1/healthz` 作为后台管理接口分组示例
- 提供 `/admin/v1/follower-users` 的 Gorm CRUD 示例实现
- `/admin/v1` 下除健康检查、登录、刷新外的后台业务接口默认走 Access Token 鉴权
- 启动时自动读取配置并初始化日志、数据库与 HTTP 服务
- 支持通过环境变量覆盖配置文件中的字段
- 日志默认同时输出到控制台与 `backend/logs/app.log`

## 本地启动

### 1. 启动 PostgreSQL

项目根目录已经提供本地 PostgreSQL：

```bash
docker compose up -d
```

默认连接信息：

- Host: `127.0.0.1`
- Port: `15432`
- Database: `pincermarket`
- User: `postgres`
- Password: `postgres`

### 2. 安装依赖

```bash
cd backend
make tidy
```

### 3. 启动服务

```bash
cd backend
make run
```

服务默认监听：

- `http://127.0.0.1:8081/healthz`
- `http://127.0.0.1:8081/admin/v1/healthz`

### 4. 初始化首个管理员

服务首次启动后会自动建表，但不会自动创建首个后台管理员。请按以下步骤处理：

1. 启动服务一次，确保 `operator_portal.admin_users` 和 `operator_portal.admin_user_seq` 已创建
2. 打开 `backend/sql/seeds/0001_admin_user_seed.sql`
3. 替换其中的 `{{REPLACE_WITH_LOGIN_NAME}}` 和 `{{REPLACE_WITH_BCRYPT_HASH}}`
4. 手工执行该 SQL

注意：

- `login_name` 必须是小写，且满足 `^[a-z0-9._-]{4,32}$`
- `password_hash` 必须是 bcrypt 哈希，禁止写入明文密码

## 配置说明

默认配置文件路径：

- 在 `backend/` 目录执行时：`configs/config.yaml`
- 在仓库根目录执行时：`backend/configs/config.yaml`

支持使用环境变量覆盖，前缀统一为 `OPERATION_ADMIN`。

例如：

```bash
export OPERATION_ADMIN_SERVER_PORT=9090
export OPERATION_ADMIN_DATABASE_HOST=127.0.0.1
export OPERATION_ADMIN_AUTH_JWT_SECRET=replace-with-a-secure-secret
```

当前鉴权相关默认配置：

- `auth.issuer=operation-admin-backend`
- `auth.access_token_ttl=2h`
- `auth.refresh_token_ttl=168h`
- `auth.jwt_secret` 仅用于本地开发占位，生产环境必须替换

## 构建与测试

```bash
cd backend
make test
make build
```
